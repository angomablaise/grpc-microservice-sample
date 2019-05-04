package user

import (
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smockoro/grpc-microservice-sample/pkg/api"
)

type server struct {
	db *sql.DB
}

func NewUserServiceServer(db *sql.DB) api.UserServiceServer {
	return &server{db: db}
}

func (s *server) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "User Database Connect Error"+err.Error())
	}

	return c, nil
}

func (s *server) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx,
		"INSERT INTO users(`name`, `age`, `mail`, `address`) VALUES(?, ?, ?, ?)",
		req.User.Name, req.User.Age, req.User.Mail, req.User.Address)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert user"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve user id"+err.Error())
	}

	return &api.CreateResponse{Id: id}, nil
}

func (s *server) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.QueryContext(ctx,
		"SELECT `id`, `name`, `age`, `mail`, `address`) FROM users WHERE `id` = ?",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select operation"+err.Error())
	}
	defer res.Close()

	if !res.Next() {
		if err := res.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to get data "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID='%d' is not found",
			req.Id))
	}

	var user api.User
	if err := res.Scan(&user.Id, &user.Name, &user.Age, &user.Mail, &user.Address); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &api.GetResponse{User: &user}, nil
}

func (s *server) Update(ctx context.Context, req *api.UpdateRequest) (*api.UpdateResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx,
		"UPDATE users SET `name`=?, `age`=?, `mail`=?, `address`=? WHERE `id`=?",
		req.User.Name, req.User.Age, req.User.Mail, req.User.Address, req.User.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update user"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.Unknown,
			fmt.Sprintf("user id %d is not found", req.User.Id))
	}

	return &api.UpdateResponse{Updated: rows}, nil
}