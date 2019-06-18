package main

import (
	"context"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"

	"github.com/mediba/server/entity"
	pb "github.com/mediba/server/protobuf"
)

const (
	port = ":50051"
)

type server struct{}

var (
	id       int32
	password string
	role     int32
	email    string
	name     string
	created  string
	updated  string
)

func (s *server) ListUsers(context.Context, *empty.Empty) (*pb.ListUsersResponse, error) {
	log.Print("success")

	db, err := sql.Open("mysql", "awsuser:awspassword@/cms")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM user")

	var users []entity.User

	for rows.Next() {
		err := rows.Scan(&id, &password, &role, &email, &name, &created, &updated)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, entity.User{
			ID:       id,
			Password: password,
			Email:    email,
			Name:     name,
			Created:  created,
			Updated:  updated,
		})
	}

	return &pb.ListUsersResponse{
		Users: mapUserToPbUser(users),
	}, nil
}

func (s *server) GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return &pb.GetUserResponse{
		User: nil,
	}, nil
}

func mapUserToPbUser(e []entity.User) []*pb.User {
	pu := make([]*pb.User, len(e))

	for i, v := range e {
		pu[i] = &pb.User{
			Id:       v.ID,
			Password: v.Password,
			Role:     v.Role,
			Email:    v.Email,
			Name:     v.Name,
			Created:  v.Created,
			Updated:  v.Updated,
		}
	}
	return pu
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
