package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/service/user/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repo.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := u.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "User Database Connect Error"+err.Error())
	}
	return c, nil
}

func (u *userRepository) Insert(ctx context.Context, user *api.User) (int64, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return -1, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx,
		"INSERT INTO users(`name`, `age`, `mail`, `address`) VALUES(?, ?, ?, ?)",
		user.Name, user.Age, user.Mail, user.Address)
	if err != nil {
		return -1, status.Error(codes.Unknown, "failed to insert user"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, status.Error(codes.Unknown, "failed to retrieve user id"+err.Error())
	}

	return id, nil
}

func (u *userRepository) SelectByID(ctx context.Context, id int64) (*api.User, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.QueryContext(ctx,
		"SELECT `id`, `name`, `age`, `mail`, `address` FROM users WHERE `id` = ?",
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

	var user api.User
	if err := res.Scan(&user.Id, &user.Name, &user.Age, &user.Mail, &user.Address); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &api.User{
		Id:      user.Id,
		Name:    user.Name,
		Age:     user.Age,
		Mail:    user.Mail,
		Address: user.Address,
	}, nil
}

func (u *userRepository) SelectAll(ctx context.Context) ([]*api.User, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `id`, `name`, `age`, `mail`, `address` FROM users")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select "+err.Error())
	}
	defer rows.Close()

	list := []*api.User{}
	for rows.Next() {
		user := new(api.User)
		if err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Mail, &user.Address); err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
		list = append(list, user)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return list, nil
}

func (u *userRepository) Update(ctx context.Context, user *api.User) (int64, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return -1, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx,
		"UPDATE users SET `name`=?, `age`=?, `mail`=?, `address`=? WHERE `id`=?",
		user.Name, user.Age, user.Mail, user.Address, user.Id)
	if err != nil {
		return -1, status.Error(codes.Unknown, "failed to update user"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return -1, status.Error(codes.Unknown, err.Error())
	}

	if rows == 0 {
		return -1, status.Error(codes.Unknown,
			fmt.Sprintf("user id %d is not found", user.Id))
	}

	return rows, nil
}

func (u *userRepository) Delete(ctx context.Context, id int64) (int64, error) {
	c, err := u.connect(ctx)
	if err != nil {
		return -1, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "DELETE FROM users WHERE `id`=?", id)
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
