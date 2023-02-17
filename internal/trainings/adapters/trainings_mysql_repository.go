package adapters

import (
	"context"
	"database/sql"
	"os"
	"sort"
	"time"

	"github.com/dbaeka/workouts-go/internal/common/auth"
	"github.com/dbaeka/workouts-go/internal/trainings/app"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type mysqlTraining struct {
	UUID     string `db:"uuid"`
	UserUUID string `db:"user_uuid"`
	User     string `db:"user"`

	Time  time.Time `db:"time"`
	Notes string    `db:"notes"`

	ProposedTime   *time.Time `db:"proposed_time"`
	MoveProposedBy *string    `db:"move_proposed_by"`
}

type MySQLTrainingsRepository struct {
	db *sqlx.DB
}

func NewMySQLTrainingsRepository(db *sqlx.DB) MySQLTrainingsRepository {
	return MySQLTrainingsRepository{
		db: db,
	}
}

func dbTrainingsToApp(dbTrainings []mysqlTraining) ([]app.Training, error) {
	var trainings []app.Training

	for _, dbTraining := range dbTrainings {
		trainings = append(trainings, app.Training(dbTraining))
	}

	sort.Slice(trainings, func(i, j int) bool { return trainings[i].Time.Before(trainings[j].Time) })

	return trainings, nil
}

func (d MySQLTrainingsRepository) AllTrainings(ctx context.Context) ([]app.Training, error) {
	return d.getTrainings(ctx, nil)
}

func (d MySQLTrainingsRepository) FindTrainingsForUser(ctx context.Context, user auth.User) ([]app.Training, error) {
	return d.getTrainings(ctx, &user)
}

func (d MySQLTrainingsRepository) getTrainings(ctx context.Context, user *auth.User) ([]app.Training, error) {
	var dbTrainings []mysqlTraining

	query := "SELECT * FROM `trainings` WHERE `time` >= ?"

	var err error
	if user != nil {
		query = query + " AND `user_uuid` = ?"
		err = d.db.GetContext(ctx, &dbTrainings, query, time.Now().Add(-time.Hour*24), user.UUID)
	} else {
		err = d.db.GetContext(ctx, &dbTrainings, query, time.Now().Add(-time.Hour*24))
	}

	if errors.Is(err, sql.ErrNoRows) {
		// in reality this date exists, even if it's not persisted
		return []app.Training{}, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get hour from db")
	}

	return dbTrainingsToApp(dbTrainings)
}

func (d MySQLTrainingsRepository) CreateTraining(_ context.Context, training app.Training, createFn func() error) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	// Defer is executed on function just before exit.
	// With defer, we are always sure that we will close our transaction properly.
	defer func() {
		err = d.finishTransaction(err, tx)
	}()

	dbTraining := mysqlTraining(training)

	err = createFn()
	if err != nil {
		return err
	}

	if err := d.upsertTraining(tx, &dbTraining); err != nil {
		return err
	}

	return nil
}

func (d MySQLTrainingsRepository) getTraining(
	ctx context.Context,
	trainingUUID *string,
	hourTime *time.Time,
	forUpdate bool,
) (*app.Training, error) {
	dbTraining := mysqlTraining{}

	var query string

	if forUpdate {
		query += " FOR UPDATE"
	}

	var err error
	if trainingUUID != nil {
		query = "SELECT * FROM `trainings` WHERE `uuid` = ?" + query
		err = d.db.GetContext(ctx, &dbTraining, query, trainingUUID)

	} else if hourTime != nil {
		query = "SELECT * FROM `trainings` WHERE `time` = ?" + query
		err = d.db.GetContext(ctx, &dbTraining, query, hourTime)
	}

	if errors.Is(err, sql.ErrNoRows) {
		// in reality this date exists, even if it's not persisted
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get hour from db")
	}
	appTraining := app.Training(dbTraining)
	return &appTraining, nil
}

func (d MySQLTrainingsRepository) CancelTraining(ctx context.Context, trainingUUID string, deleteFn func(app.Training) error) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	// Defer is executed on function just before exit.
	// With defer, we are always sure that we will close our transaction properly.
	defer func() {
		err = d.finishTransaction(err, tx)
	}()

	appTraining, err := d.getTraining(ctx, &trainingUUID, nil, true)
	if err != nil || appTraining == nil {
		return err
	}

	err = deleteFn(*appTraining)
	if err != nil {
		return err
	}

	dbTraining := mysqlTraining(*appTraining)
	if err := d.deleteTraining(tx, &dbTraining); err != nil {
		return err
	}

	return nil
}

func (d MySQLTrainingsRepository) RescheduleTraining(
	ctx context.Context,
	trainingUUID string,
	newTime time.Time,
	updateFn func(app.Training) (app.Training, error),
) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	// Defer is executed on function just before exit.
	// With defer, we are always sure that we will close our transaction properly.
	defer func() {
		err = d.finishTransaction(err, tx)
	}()

	appTraining, err := d.getTraining(ctx, &trainingUUID, nil, true)
	if err != nil || appTraining == nil {
		return err
	}

	// check if new time already taken
	existingTraining, err := d.getTraining(ctx, nil, &newTime, false)
	if err != nil {
		return err
	}

	if existingTraining != nil {
		return errors.Errorf("there is training already at %s", newTime)
	}

	updatedTraining, err := updateFn(*appTraining)
	if err != nil {
		return err
	}

	dbTraining := mysqlTraining(updatedTraining)
	if err := d.upsertTraining(tx, &dbTraining); err != nil {
		return err
	}

	return nil
}

func (d MySQLTrainingsRepository) ApproveTrainingReschedule(ctx context.Context, trainingUUID string, updateFn func(app.Training) (app.Training, error)) error {
	return d.updateTraining(ctx, trainingUUID, updateFn)
}

func (d MySQLTrainingsRepository) RejectTrainingReschedule(ctx context.Context, trainingUUID string, updateFn func(app.Training) (app.Training, error)) error {
	return d.updateTraining(ctx, trainingUUID, updateFn)
}

func (d MySQLTrainingsRepository) updateTraining(ctx context.Context, trainingUUID string, updateFn func(app.Training) (app.Training, error)) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	// Defer is executed on function just before exit.
	// With defer, we are always sure that we will close our transaction properly.
	defer func() {
		err = d.finishTransaction(err, tx)
	}()

	appTraining, err := d.getTraining(ctx, &trainingUUID, nil, true)
	if err != nil || appTraining == nil {
		return err
	}

	updatedTraining, err := updateFn(*appTraining)
	if err != nil {
		return err
	}

	dbTraining := mysqlTraining(updatedTraining)
	if err := d.upsertTraining(tx, &dbTraining); err != nil {
		return err
	}

	return nil
}

// deleteTraining updates training if training already exists in the database.
// If it doesn't exist, it's inserted.
func (d MySQLTrainingsRepository) deleteTraining(tx *sqlx.Tx, trainingToDelete *mysqlTraining) error {
	_, err := tx.NamedExec(
		`DELETE FROM trainings WHERE uuid = :uuid`,
		trainingToDelete,
	)
	if err != nil {
		return errors.Wrap(err, "unable to delete training")
	}

	return nil
}

// upsertTraining updates training if training already exists in the database.
// If it doesn't exist, it's inserted.
func (d MySQLTrainingsRepository) upsertTraining(tx *sqlx.Tx, trainingToUpdate *mysqlTraining) error {
	_, err := tx.NamedExec(
		`INSERT INTO 
			trainings (uuid, user_uuid, user, time, notes, proposed_time, move_proposed_by) 
		VALUES 
			(:uuid, :user_uuid, :user, :time, :notes, :proposed_time, :move_proposed_by)
		ON DUPLICATE KEY UPDATE 
			user = :user`,
		trainingToUpdate,
	)
	if err != nil {
		return errors.Wrap(err, "unable to upsert training")
	}

	return nil
}

func (d MySQLTrainingsRepository) finishTransaction(err error, tx *sqlx.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return multierr.Combine(err, rollbackErr)
		}

		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return errors.Wrap(err, "failed to commit tx")
		}

		return nil
	}
}

func NewMySQLConnection() (*sqlx.DB, error) {
	config := mysql.Config{
		Addr:      os.Getenv("MYSQL_ADDRESS"),
		User:      os.Getenv("MYSQL_USERNAME"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		DBName:    os.Getenv("MYSQL_DATABASE"),
		ParseTime: true, // with that parameter, we can use time.Time in mysqlHour.Hour
	}

	db, err := sqlx.Connect("mysql", config.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to MySQL")
	}

	return db, nil
}
