package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/service/user/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repo.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Insert(ctx context.Context, user *api.User) (int64, error) {
	res, err := u.db.NamedExecContext(ctx,
		"INSERT INTO users(`name`, `age`, `mail`, `address`) VALUES(:name, :age, :mail, :address)",
		user)
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
	res, err := u.db.QueryxContext(ctx,
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
	if err := res.StructScan(&user); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &user, nil
}

func (u *userRepository) SelectAll(ctx context.Context) ([]*api.User, error) {
	rows, err := u.db.QueryxContext(ctx, "SELECT `id`, `name`, `age`, `mail`, `address` FROM users")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select "+err.Error())
	}
	defer rows.Close()

	list := []*api.User{}
	for rows.Next() {
		var user api.User
		if err := rows.StructScan(&user); err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
		list = append(list, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return list, nil
}

func (u *userRepository) Update(ctx context.Context, user *api.User) (int64, error) {
	res, err := u.db.NamedExecContext(ctx,
		"UPDATE users SET `name`=:name, `age`=:age, `mail`=:mail, `address`=:address WHERE `id`=:id",
		user)
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
	res, err := u.db.ExecContext(ctx, "DELETE FROM users WHERE `id`=:id", id)
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
