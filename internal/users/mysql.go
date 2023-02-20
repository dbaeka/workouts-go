package main

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"os"
)

type db struct {
	mysqlClient *sqlx.DB
}

type mysqlUser struct {
	UUID        string `db:"uuid"`
	Balance     int    `db:"balance"`
	DisplayName string `db:"display_name"`
	Role        string `db:"role"`
	LastIP      string `db:"last_ip"`
}

func (d db) GetUser(ctx context.Context, userID string) (*mysqlUser, error) {
	return d.getUser(ctx, d.mysqlClient, userID, false)
}

// sqlContextGetter is an interface provided both by transaction and standard db connection
type sqlContextGetter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func (d db) getUser(
	ctx context.Context,
	db sqlContextGetter,
	userID string,
	forUpdate bool,
) (*mysqlUser, error) {
	dbUser := mysqlUser{}
	query := "SELECT * FROM `hours` WHERE `hour` = ?"
	if forUpdate {
		query += " FOR UPDATE"
	}

	err := db.GetContext(ctx, &dbUser, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		// in reality this date exists, even if it's not persisted
		return &mysqlUser{
			Balance: 0,
		}, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get user from db")
	}

	return &dbUser, nil
}

func (d db) UpdateUser(
	ctx context.Context,
	userID string,
	updateFn func(u *mysqlUser) (*mysqlUser, error),
) error {
	tx, err := d.mysqlClient.Beginx()
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	// Defer is executed on function just before exit.
	// With defer, we are always sure that we will close our transaction properly.
	defer func() {
		// In `UpdateHour` we are using named return - `(err error)`.
		// Thanks to that, that can check if function exits with error.
		//
		// Even if function exits without error, commit still can return error.
		// In that case we can override nil to err `err = m.finish...`.
		err = d.finishTransaction(err, tx)
	}()

	existingUser, err := d.getUser(ctx, tx, userID, true)
	if err != nil {
		return err
	}

	updatedUser, err := updateFn(existingUser)
	if err != nil {
		return err
	}

	if err := d.upsertUser(tx, updatedUser); err != nil {
		return err
	}

	return nil
}

// upsertUser updates hour if hour already exists in the database.
// If your doesn't exist, it's inserted.
func (d db) upsertUser(tx *sqlx.Tx, userToUpdate *mysqlUser) error {
	_, err := tx.NamedExec(
		`INSERT INTO 
			users (uuid, display_name, role, balance, last_ip) 
		VALUES 
			(:uuid, :display_name, :role, :balance, :last_ip)
		ON DUPLICATE KEY UPDATE 
			balance = :balance, last_ip = :last_ip`,
		userToUpdate,
	)
	if err != nil {
		return errors.Wrap(err, "unable to upsert hour")
	}

	return nil
}

// finishTransaction rollbacks transaction if error is provided.
// If err is nil transaction is committed.
//
// If the rollback fails, we are using multierr library to add error about rollback failure.
// If the commit fails, commit error is returned.
func (d db) finishTransaction(err error, tx *sqlx.Tx) error {
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
