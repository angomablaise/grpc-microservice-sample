package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/repository/mysql/user"
)

type lastInsertIdError struct{}

func (lie *lastInsertIdError) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("error")
}
func (lie *lastInsertIdError) RowsAffected() (int64, error) {
	return 1, nil
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	user := &api.User{Name: "Bob", Age: 11, Mail: "sample@sample.com", Address: "Tokyo"}

	ur := repo.NewUserRepository(sqlxDB)
	ctx := context.Background()
	if _, err = ur.Insert(ctx, user); err == nil {
		t.Errorf("error was expected while Insert stats: %s", err)
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Name, user.Age, user.Mail, user.Address).
		WillReturnResult(sqlmock.NewResult(1, 1))
	ctx = context.Background()
	if _, err = ur.Insert(ctx, user); err != nil {
		t.Errorf("error was not expected while Insert stats: %s", err)
	}

	mock.ExpectExec("INSERT INTO users").
		WillReturnResult(&lastInsertIdError{})
	if _, err = ur.Insert(ctx, user); err == nil {
		t.Errorf("error was expected while Insert stats: %s", err)
	}

}

func TestSelectByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	ur := repo.NewUserRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"id", "name", "age", "mail", "address"}).
		AddRow(1, "Bob", 11, "sample@sample.com", "Tokyo")
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").
		WillReturnRows(rows)
	ctx := context.Background()
	_, err = ur.SelectByID(ctx, 1)
	if err != nil {
		t.Errorf("error was not expected while Select by ID stats: %s", err)
	}

	ctx = context.Background()
	_, err = ur.SelectByID(ctx, 2)
	if err == nil {
		t.Errorf("error was not expected while Select by ID stats: %s", err)
	}

	rows = sqlmock.NewRows([]string{"id", "BAD"}).
		AddRow(1, "")
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").
		WillReturnRows(rows)
	ctx = context.Background()
	_, err = ur.SelectByID(ctx, 1)
	if err == nil {
		t.Errorf("error was not expected while Select by ID stats: %s", err)
	}

	rows = sqlmock.NewRows([]string{"id", "BAD"})
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").
		WillReturnRows(rows)
	ctx = context.Background()
	_, err = ur.SelectByID(ctx, 1)
	if err == nil {
		t.Errorf("error was not expected while Select by ID stats: %s", err)
	}

	rows = sqlmock.NewRows([]string{"id", "BAD"}).
		RowError(1, fmt.Errorf("error"))
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").
		WillReturnRows(rows)
	ctx = context.Background()
	_, err = ur.SelectByID(ctx, 1)
	if err == nil {
		t.Errorf("error was not expected while Select by ID stats: %s", err)
	}
}

func TestSelectAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	ur := repo.NewUserRepository(sqlxDB)

	ctx := context.Background()
	if _, err = ur.SelectAll(ctx); err == nil {
		t.Errorf("error was expected while Select All stats: %s", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "age", "mail", "address"}).
		AddRow(1, "Bob", 11, "sample@sample.com", "Tokyo").
		AddRow(2, "Alice", 13, "example@sample.com", "London")
	mock.ExpectQuery("^SELECT (.+) FROM users$").
		WillReturnRows(rows)
	ctx = context.Background()
	if _, err = ur.SelectAll(ctx); err != nil {
		t.Errorf("error was not expected while Select All stats: %s", err)
	}

	rows = sqlmock.NewRows([]string{"id", "BAD"}).
		AddRow(1, "Bob").
		AddRow(2, "Alice")
	mock.ExpectQuery("^SELECT (.+) FROM users$").
		WillReturnRows(rows)
	ctx = context.Background()
	if _, err = ur.SelectAll(ctx); err == nil {
		t.Errorf("error was expected while Select All stats: %s", err)
	}

	rows = sqlmock.NewRows([]string{"id", "BAD"}).
		AddRow(1, "Bob").
		RowError(1, fmt.Errorf("error"))
	mock.ExpectQuery("^SELECT (.+) FROM users$").
		WillReturnRows(rows)
	ctx = context.Background()
	if _, err = ur.SelectAll(ctx); err == nil {
		t.Errorf("error was expected while Select All stats: %s", err)
	}
}

type rowsAffectedError struct{}

func (rae *rowsAffectedError) LastInsertId() (int64, error) {
	return 1, nil
}
func (rae *rowsAffectedError) RowsAffected() (int64, error) {
	return 0, fmt.Errorf("error")
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	user := &api.User{Id: 1, Name: "Bob Olimar", Age: 88, Mail: "aa@sample.com", Address: "Tokyo"}
	ur := repo.NewUserRepository(sqlxDB)

	ctx := context.Background()
	if _, err = ur.Update(ctx, user); err == nil {
		t.Errorf("error was not expected while Update stats: %s", err)
	}

	mock.ExpectExec("UPDATE users SET").WillReturnResult(sqlmock.NewResult(1, 1))
	ctx = context.Background()
	if _, err = ur.Update(ctx, user); err != nil {
		t.Errorf("error was not expected while Update stats: %s", err)
	}

	mock.ExpectExec("UPDATE users SET").WillReturnResult(&rowsAffectedError{})
	ctx = context.Background()
	if _, err = ur.Update(ctx, user); err == nil {
		t.Errorf("error was not expected while Update stats: %s", err)
	}

	mock.ExpectExec("UPDATE users SET").WillReturnResult(sqlmock.NewResult(1, 0))
	ctx = context.Background()
	if _, err = ur.Update(ctx, user); err == nil {
		t.Errorf("error was not expected while Update stats: %s", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	ur := repo.NewUserRepository(sqlxDB)

	ctx := context.Background()
	if _, err = ur.Delete(ctx, 1); err == nil {
		t.Errorf("error was not expected while Delete stats: %s", err)
	}

	mock.ExpectExec("DELETE FROM users WHERE").WillReturnResult(sqlmock.NewResult(1, 1))
	ctx = context.Background()
	if _, err = ur.Delete(ctx, 1); err != nil {
		t.Errorf("error was not expected while Delete stats: %s", err)
	}

	mock.ExpectExec("DELETE FROM users WHERE").WillReturnResult(&rowsAffectedError{})
	ctx = context.Background()
	if _, err = ur.Delete(ctx, 1); err == nil {
		t.Errorf("error was not expected while Delete stats: %s", err)
	}

	mock.ExpectExec("DELETE FROM users WHERE").WillReturnResult(sqlmock.NewResult(1, 0))
	ctx = context.Background()
	if _, err = ur.Delete(ctx, 1); err == nil {
		t.Errorf("error was not expected while Delete stats: %s", err)
	}
}
