package v1

import (
	v1 "aniqma/aniqma/crudgrpc/pkg/api/v1"
	"database/sql"

	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/lib/pq"
)

const (
	apiVersion = "v1"
)

type uSersServiceServer struct {
	db *sql.DB
}

// NewUsesServiceServer
func NewUsersServiceServer(db *sql.DB) v1.UsersServiceServer {
	return &uSersServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *uSersServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// connect returns SQL database connection from the pool
func (s *uSersServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create new todo task
func (s *uSersServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.USers.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	// insert Users entity data
	res, err := c.ExecContext(ctx, "INSERT INTO users(Username, Password, Reminder) VALUES($1, $2, $3) RETURNING id",
		req.USers.Username, req.USers.Password, reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Users-> "+err.Error())
	}
	// get ID of creates Users
	id, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created Users-> "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

// Read users task
func (s *uSersServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// query Users by ID
	rows, err := c.QueryContext(ctx, "SELECT id, username, password, reminder FROM users WHERE id= $1",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Users-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from Users-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Users with ID='%d' is not found",
			req.Id))
	}

	// get Users data
	var td v1.Users
	var reminder time.Time
	if err := rows.Scan(&td.Id, &td.Username, &td.Password, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from Users row-> "+err.Error())
	}
	td.Reminder, err = ptypes.TimestampProto(reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple Users rows with ID='%d'",
			req.Id))
	}

	return &v1.ReadResponse{
		Api:   apiVersion,
		USers: &td,
	}, nil

}

// Update Users task
func (s *uSersServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.USers.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	// update Users
	res, err := c.ExecContext(ctx, "UPDATE users SET username=$1, password=$2, reminder=$3 WHERE id=$4",
		req.USers.Username, req.USers.Password, reminder, req.USers.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update Users-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Users with ID='%d' is not found",
			req.USers.Id))
	}

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil
}

// Delete Users task
func (s *uSersServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// delete Users
	res, err := c.ExecContext(ctx, "DELETE FROM users WHERE id=$1", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete Users-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Users with ID='%d' is not found",
			req.Id))
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: rows,
	}, nil
}

// Read all Users tasks
func (s *uSersServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// get Users list
	rows, err := c.QueryContext(ctx, "SELECT id, username, password, reminder FROM users")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Users-> "+err.Error())
	}
	defer rows.Close()

	var reminder time.Time
	list := []*v1.Users{}
	for rows.Next() {
		td := new(v1.Users)
		if err := rows.Scan(&td.Id, &td.Username, &td.Password, &reminder); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from Users row-> "+err.Error())
		}
		td.Reminder, err = ptypes.TimestampProto(reminder)
		if err != nil {
			return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
		}
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from Users-> "+err.Error())
	}

	return &v1.ReadAllResponse{
		Api:    apiVersion,
		USerse: list,
	}, nil
}
