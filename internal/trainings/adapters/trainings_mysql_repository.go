package adapters

import (
	"context"
	"database/sql"
	"os"
	"sort"
	"time"

	"github.com/dbaeka/workouts-go/internal/trainings/app/query"
	"github.com/dbaeka/workouts-go/internal/trainings/domain/training"

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

	Canceled bool `db:"canceled"`
}

type MySQLTrainingsRepository struct {
	db *sqlx.DB
}

func NewMySQLTrainingsRepository(db *sqlx.DB) MySQLTrainingsRepository {
	return MySQLTrainingsRepository{
		db: db,
	}
}

func (r MySQLTrainingsRepository) marshalTraining(tr *training.Training) mysqlTraining {
	dbTraining := mysqlTraining{
		UUID:     tr.UUID(),
		UserUUID: tr.UserUUID(),
		User:     tr.UserName(),
		Time:     tr.Time(),
		Notes:    tr.Notes(),
		Canceled: tr.IsCanceled(),
	}

	if tr.IsRescheduleProposed() {
		proposedBy := tr.MovedProposedBy().String()
		proposedTime := tr.ProposedNewTime()

		dbTraining.MoveProposedBy = &proposedBy
		dbTraining.ProposedTime = &proposedTime
	}

	return dbTraining
}

func (r MySQLTrainingsRepository) unmarshalTraining(dbTraining *mysqlTraining) (*training.Training, error) {
	var moveProposedBy training.UserType
	var err error
	if dbTraining.MoveProposedBy != nil {
		moveProposedBy, err = training.NewUserTypeFromString(*dbTraining.MoveProposedBy)
		if err != nil {
			return nil, err
		}
	}

	var proposedTime time.Time
	if dbTraining.ProposedTime != nil {
		proposedTime = *dbTraining.ProposedTime
	}

	return training.UnmarshalTrainingFromDatabase(
		dbTraining.UUID,
		dbTraining.UserUUID,
		dbTraining.User,
		dbTraining.Time,
		dbTraining.Notes,
		dbTraining.Canceled,
		proposedTime,
		moveProposedBy,
	)
}

func (r MySQLTrainingsRepository) trainingModelsToQuery(dbTrainings *[]mysqlTraining) ([]query.Training, error) {
	var trainings []query.Training

	for _, dbTraining := range *dbTrainings {
		dbTraining := dbTraining
		tr, err := r.unmarshalTraining(&dbTraining)
		if err != nil {
			return nil, err
		}

		queryTraining := query.Training{
			UUID:           tr.UUID(),
			UserUUID:       tr.UserUUID(),
			User:           tr.UserName(),
			Time:           tr.Time(),
			Notes:          tr.Notes(),
			CanBeCancelled: tr.CanBeCanceledForFree(),
		}

		if tr.IsRescheduleProposed() {
			proposedTime := tr.ProposedNewTime()
			queryTraining.ProposedTime = &proposedTime

			proposedBy := tr.MovedProposedBy().String()
			queryTraining.MoveProposedBy = &proposedBy
		}

		trainings = append(trainings, queryTraining)
	}

	sort.Slice(trainings, func(i, j int) bool { return trainings[i].Time.Before(trainings[j].Time) })

	return trainings, nil
}

func (r MySQLTrainingsRepository) AllTrainings(ctx context.Context) ([]query.Training, error) {
	return r.getTrainings(ctx, nil)
}

func (r MySQLTrainingsRepository) FindTrainingsForUser(ctx context.Context, userUUID string) ([]query.Training, error) {
	return r.getTrainings(ctx, &userUUID)
}

func (r MySQLTrainingsRepository) getTrainings(ctx context.Context, userUUID *string) ([]query.Training, error) {
	var dbTrainings []mysqlTraining

	sqlQuery := "SELECT * FROM `trainings` WHERE `time` >= ?"

	var err error
	if userUUID != nil {
		sqlQuery = sqlQuery + " AND `user_uuid` = ?"
		err = r.db.GetContext(ctx, &dbTrainings, sqlQuery, time.Now().Add(-time.Hour*24), userUUID)
	} else {
		err = r.db.GetContext(ctx, &dbTrainings, sqlQuery, time.Now().Add(-time.Hour*24))
	}

	if errors.Is(err, sql.ErrNoRows) {
		// in reality this date exists, even if it's not persisted
		return []query.Training{}, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get training from db")
	}

	return r.trainingModelsToQuery(&dbTrainings)
}

func (r MySQLTrainingsRepository) AddTraining(_ context.Context, tr *training.Training) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	// Defer is executed on function just before exit.
	// With defer, we are always sure that we will close our transaction properly.
	defer func() {
		err = r.finishTransaction(err, tx)
	}()

	dbTraining := r.marshalTraining(tr)

	if err := r.upsertTraining(tx, &dbTraining); err != nil {
		return err
	}

	return nil
}

func (r MySQLTrainingsRepository) GetTraining(
	ctx context.Context,
	trainingUUID string,
	user training.User,
) (*training.Training, error) {
	return r.getTraining(ctx, trainingUUID, user, false)
}

func (r MySQLTrainingsRepository) getTraining(
	ctx context.Context,
	trainingUUID string,
	user training.User,
	forUpdate bool,
) (*training.Training, error) {
	dbTraining := mysqlTraining{}

	var sqlQuery string

	if forUpdate {
		sqlQuery += " FOR UPDATE"
	}

	var err error
	sqlQuery = "SELECT * FROM `trainings` WHERE `uuid` = ?" + sqlQuery
	err = r.db.GetContext(ctx, &dbTraining, sqlQuery, trainingUUID)

	if errors.Is(err, sql.ErrNoRows) {
		// in reality this date exists, even if it's not persisted
		return nil, training.NotFoundError{TrainingUUID: trainingUUID}
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get training from db")
	}

	tr, err := r.unmarshalTraining(&dbTraining)
	if err != nil {
		return nil, err
	}

	if err := training.CanUserSeeTraining(user, *tr); err != nil {
		return nil, err
	}
	return tr, nil
}

func (r MySQLTrainingsRepository) UpdateTraining(
	ctx context.Context,
	trainingUUID string,
	user training.User,
	updateFn func(ctx context.Context, tr *training.Training) (*training.Training, error)) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	// Defer is executed on function just before exit.
	// With defer, we are always sure that we will close our transaction properly.
	defer func() {
		err = r.finishTransaction(err, tx)
	}()

	appTraining, err := r.getTraining(ctx, trainingUUID, user, true)
	if err != nil || appTraining == nil {
		return err
	}

	updatedTraining, err := updateFn(ctx, appTraining)
	if err != nil {
		return err
	}

	dbTraining := r.marshalTraining(updatedTraining)
	if err := r.upsertTraining(tx, &dbTraining); err != nil {
		return err
	}

	return nil
}

// RemoveAllTrainings was designed for tests for doing data cleanups
func (r MySQLTrainingsRepository) RemoveAllTrainings(ctx context.Context) error {
	sqlQuery := `TRUNCATE trainings`
	_, err := r.db.ExecContext(ctx, sqlQuery)
	if err != nil {
		return errors.Wrap(err, "unable to delete all trainings")
	}

	return nil
}

// upsertTraining updates training if training already exists in the database.
// If it doesn't exist, it's inserted.
func (r MySQLTrainingsRepository) upsertTraining(tx *sqlx.Tx, trainingToUpdate *mysqlTraining) error {
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

func (r MySQLTrainingsRepository) finishTransaction(err error, tx *sqlx.Tx) error {
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
