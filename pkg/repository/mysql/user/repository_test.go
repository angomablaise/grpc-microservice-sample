package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/smockoro/grpc-microservice-sample/pkg/api"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := &api.User{Name: "Bob", Age: 11, Mail: "sample@sample.com", Address: "Tokyo"}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Name, user.Age, user.Mail, user.Address).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ur := NewUserRepository(db)
	ctx := context.Background()
	if _, err = ur.Insert(ctx, user); err != nil {
		t.Errorf("error was not expected while Insert stats: %s", err)
	}
}

func TestSelectByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "age", "mail", "address"}).
		AddRow(1, "Bob", 11, "sample@sample.com", "Tokyo")
		//AddRow(2, "Alice", 13, "example@sample.com", "London")

	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").
		WillReturnRows(rows)

	ur := NewUserRepository(db)
	ctx := context.Background()
	_, err = ur.SelectByID(ctx, 1)
	if err != nil {
		t.Errorf("error was not expected while Select by ID stats: %s", err)
	}
}

func TestSelectAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "age", "mail", "address"}).
		AddRow(1, "Bob", 11, "sample@sample.com", "Tokyo").
		AddRow(2, "Alice", 13, "example@sample.com", "London")

	mock.ExpectQuery("^SELECT (.+) FROM users$").
		WillReturnRows(rows)

	ur := NewUserRepository(db)
	ctx := context.Background()
	if _, err = ur.SelectAll(ctx); err != nil {
		t.Errorf("error was not expected while Select All stats: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := &api.User{Id: 1, Name: "Bob Olimar", Age: 88, Mail: "aa@sample.com", Address: "Tokyo"}
	mock.ExpectExec("UPDATE users SET").WillReturnResult(sqlmock.NewResult(1, 1))

	ur := NewUserRepository(db)
	ctx := context.Background()
	if _, err = ur.Update(ctx, user); err != nil {
		t.Errorf("error was not expected while Update stats: %s", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM users WHERE").WillReturnResult(sqlmock.NewResult(1, 1))

	ur := NewUserRepository(db)
	ctx := context.Background()
	if _, err = ur.Delete(ctx, 1); err != nil {
		t.Errorf("error was not expected while Delete stats: %s", err)
	}
}
