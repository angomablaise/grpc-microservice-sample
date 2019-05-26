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

	item := &api.Item{Name: "Apple", Description: "Red Apple", Price: 120}

	mock.ExpectExec("INSERT INTO items").
		WithArgs(item.Name, item.Description, item.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ur := NewItemRepository(db)
	ctx := context.Background()
	if _, err = ur.Insert(ctx, item); err != nil {
		t.Errorf("error was not expected while Insert stats: %s", err)
	}
}

func TestSelectByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
		AddRow(1, "Apple", "Red Apple", 120)

	mock.ExpectQuery("^SELECT (.+) FROM items WHERE").
		WillReturnRows(rows)

	ur := NewItemRepository(db)
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

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price"}).
		AddRow(1, "Apple", "Red Apple", 120).
		AddRow(2, "Pen", "HB pencil", 100)

	mock.ExpectQuery("^SELECT (.+) FROM items$").
		WillReturnRows(rows)

	ur := NewItemRepository(db)
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

	item := &api.Item{Id: 1, Name: "Apple", Description: "Red Apple", Price: 500}
	mock.ExpectExec("UPDATE items SET").WillReturnResult(sqlmock.NewResult(1, 1))

	ur := NewItemRepository(db)
	ctx := context.Background()
	if _, err = ur.Update(ctx, item); err != nil {
		t.Errorf("error was not expected while Update stats: %s", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM items WHERE").WillReturnResult(sqlmock.NewResult(1, 1))

	ur := NewItemRepository(db)
	ctx := context.Background()
	if _, err = ur.Delete(ctx, 1); err != nil {
		t.Errorf("error was not expected while Delete stats: %s", err)
	}
}
