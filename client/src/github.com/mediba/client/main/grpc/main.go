package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mediba/client/entity"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/mediba/client/protobuf"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ListUsers(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("failed grpc: %v", err)
	}

	users := mapPbToUser(r.Users)

	for _, v := range users {
		log.Printf("ID: %d", v.ID)
		log.Printf("Password: %s", v.Password)
		log.Printf("Role: %d", v.Role)
		log.Printf("Email: %s", v.Email)
		log.Printf("Name: %s", v.Name)
		log.Printf("Created: %s", v.Created)
		log.Printf("Updated: %s", v.Updated)
		log.Printf("============")
	}
}

func mapPbToUser(pb []*pb.User) []entity.User {
	eu := make([]entity.User, len(pb))

	for i, v := range pb {
		eu[i] = entity.User{
			ID: v.Id,
			Password: v.Password,
			Role: v.Role,
			Email: v.Email,
			Name: v.Name,
			Created: v.Created,
			Updated: v.Updated,
		}
	}
	return eu
}
