package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/smockoro/grpc-microservice-sample/pkg/api"
)

type lastInsertIDError struct{}

func (lie *lastInsertIDError) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("error")
}
func (lie *lastInsertIDError) RowsAffected() (int64, error) {
	return 1, nil
}

type rowsAffectedError struct{}

func (rae *rowsAffectedError) LastInsertId() (int64, error) {
	return 1, nil
}
func (rae *rowsAffectedError) RowsAffected() (int64, error) {
	return 0, fmt.Errorf("error")
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ur := NewItemRepository(db)
	item := &api.Item{Name: "Apple", Description: "Red Apple", Price: 120}

	ctx := context.Background()
	if _, err = ur.Insert(ctx, item); err == nil {
		t.Errorf("error was expected while Insert stats: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	cancel()
	if _, err = ur.Insert(ctx, item); err == nil {
		t.Errorf("error was expected while Insert stats: %s", err)
	}

	mock.ExpectExec("INSERT INTO items").
		WithArgs(item.Name, item.Description, item.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))
	ctx = context.Background()
	if _, err = ur.Insert(ctx, item); err != nil {
		t.Errorf("error was not expected while Insert stats: %s", err)
	}

	mock.ExpectExec("INSERT INTO items").
		WithArgs(item.Name, item.Description, item.Price).
		WillReturnResult(&lastInsertIDError{})
	ctx = context.Background()
	if _, err = ur.Insert(ctx, item); err == nil {
		t.Errorf("error was expected while Insert stats: %s", err)
	}
}

func TestSelectByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	ur := NewItemRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	cancel()
	_, err = ur.SelectByID(ctx, 1)
	if err == nil {
		t.Errorf("error was expected while Select by ID stats: %s", err)
	}

	ctx = context.Background()
	_, err = ur.SelectByID(ctx, 1)
	if err == nil {
		t.Errorf("error was expected while Select by ID stats: %s", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
		AddRow(1, "Apple", "Red Apple", 120)
	mock.ExpectQuery("^SELECT (.+) FROM items WHERE").
		WillReturnRows(rows)
	ctx = context.Background()
	_, err = ur.SelectByID(ctx, 1)
	if err != nil {
		t.Errorf("error was not expected while Select by ID stats: %s", err)
	}

	rows = sqlmock.NewRows([]string{"id", "BAD"}).
		RowError(1, fmt.Errorf("error"))
	mock.ExpectQuery("^SELECT (.+) FROM items WHERE").
		WillReturnRows(rows)
	ctx = context.Background()
	_, err = ur.SelectByID(ctx, 1)
	if err == nil {
		t.Errorf("error was expected while Select by ID stats: %s", err)
	}
}

func TestSelectAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	ur := NewItemRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	cancel()
	if _, err = ur.SelectAll(ctx); err == nil {
		t.Errorf("error was expected while Select All stats: %s", err)
	}

	ctx = context.Background()
	if _, err = ur.SelectAll(ctx); err == nil {
		t.Errorf("error was expected while Select All stats: %s", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
		AddRow(1, "Apple", "Red Apple", 120).
		AddRow(2, "Pen", "HB pencil", 100)
	mock.ExpectQuery("^SELECT (.+) FROM items$").
		WillReturnRows(rows)
	ctx = context.Background()
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
	ur := NewItemRepository(db)
	item := &api.Item{Id: 1, Name: "Apple", Description: "Red Apple", Price: 500}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	cancel()
	if _, err = ur.Update(ctx, item); err == nil {
		t.Errorf("error was expected while Update stats: %s", err)
	}

	ctx = context.Background()
	if _, err = ur.Update(ctx, item); err == nil {
		t.Errorf("error was expected while Update stats: %s", err)
	}

	mock.ExpectExec("UPDATE items SET").WillReturnResult(sqlmock.NewResult(1, 1))
	ctx = context.Background()
	if _, err = ur.Update(ctx, item); err != nil {
		t.Errorf("error was not expected while Update stats: %s", err)
	}

	mock.ExpectExec("UPDATE items SET").WillReturnResult(&rowsAffectedError{})
	ctx = context.Background()
	if _, err = ur.Update(ctx, item); err == nil {
		t.Errorf("error was expected while Update stats: %s", err)
	}

	mock.ExpectExec("UPDATE items SET").WillReturnResult(sqlmock.NewResult(1, 0))
	ctx = context.Background()
	if _, err = ur.Update(ctx, item); err == nil {
		t.Errorf("error was expected while Update stats: %s", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	ur := NewItemRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	cancel()
	if _, err = ur.Delete(ctx, 1); err == nil {
		t.Errorf("error was expected while Delete stats: %s", err)
	}

	ctx = context.Background()
	if _, err = ur.Delete(ctx, 1); err == nil {
		t.Errorf("error was expected while Delete stats: %s", err)
	}

	mock.ExpectExec("DELETE FROM items WHERE").WillReturnResult(sqlmock.NewResult(1, 1))
	ctx = context.Background()
	if _, err = ur.Delete(ctx, 1); err != nil {
		t.Errorf("error was not expected while Delete stats: %s", err)
	}

	mock.ExpectExec("DELETE FROM items WHERE").WillReturnResult(&rowsAffectedError{})
	ctx = context.Background()
	if _, err = ur.Delete(ctx, 1); err == nil {
		t.Errorf("error was expected while Delete stats: %s", err)
	}

	mock.ExpectExec("DELETE FROM items WHERE").WillReturnResult(sqlmock.NewResult(1, 0))
	ctx = context.Background()
	if _, err = ur.Delete(ctx, 1); err == nil {
		t.Errorf("error was expected while Delete stats: %s", err)
	}
}
