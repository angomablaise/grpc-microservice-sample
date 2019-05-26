package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/service/item/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type itemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) repo.ItemRepository {
	return &itemRepository{db: db}
}

func (u *itemRepository) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := u.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Item Database Connect Error"+err.Error())
	}
	return c, nil
}

func (u *itemRepository) Insert(ctx context.Context, item *api.Item) (int64, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return -1, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx,
		"INSERT INTO items(`name`, `description`, `price`) VALUES(?, ?, ?)",
		item.Name, item.Description, item.Price)
	if err != nil {
		return -1, status.Error(codes.Unknown, "failed to insert item"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, status.Error(codes.Unknown, "failed to retrieve item id"+err.Error())
	}

	return id, nil
}

func (u *itemRepository) SelectByID(ctx context.Context, id int64) (*api.Item, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.QueryContext(ctx,
		"SELECT `id`, `name`, `description`, `price` FROM items WHERE `id` = ?",
		id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select operation"+err.Error())
	}
	defer res.Close()

	if !res.Next() {
		if err := res.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to get data "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID='%d' is not found",
			id))
	}

	var item api.Item
	if err := res.Scan(&item.Id, &item.Name, &item.Description, &item.Price); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &api.Item{
		Id:          item.Id,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
	}, nil
}

func (u *itemRepository) SelectAll(ctx context.Context) ([]*api.Item, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `id`, `name`, `description`, `price` FROM items")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select "+err.Error())
	}
	defer rows.Close()

	list := []*api.Item{}
	for rows.Next() {
		item := new(api.Item)
		if err := rows.Scan(&item.Id, &item.Name, &item.Description, &item.Price); err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return list, nil
}

func (u *itemRepository) Update(ctx context.Context, item *api.Item) (int64, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return -1, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx,
		"UPDATE items SET `name`=?, `description`=?, `price`=? WHERE `id`=?",
		item.Name, item.Description, item.Price, item.Id)
	if err != nil {
		return -1, status.Error(codes.Unknown, "failed to update item"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return -1, status.Error(codes.Unknown, err.Error())
	}

	if rows == 0 {
		return -1, status.Error(codes.Unknown,
			fmt.Sprintf("item id %d is not found", item.Id))
	}

	return rows, nil
}

func (u *itemRepository) Delete(ctx context.Context, id int64) (int64, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return -1, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "DELETE FROM items WHERE `id`=?", id)
	if err != nil {
		return -1, status.Error(codes.Unknown, "failed to delete "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return -1, status.Error(codes.Unknown, err.Error())
	}

	if rows == 0 {
		return -1, status.Error(codes.NotFound, fmt.Sprintf("ID='%d' is not found",
			id))
	}

	return rows, nil
}
