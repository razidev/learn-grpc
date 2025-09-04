package main

import (
	"context"
	"grpc/pb/user"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type userService struct {
	user.UnimplementedUserServiceServer
}

func (us *userService) CreateUser(ctx context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	log.Println("Create user request received")
	return &user.CreateResponse{
		Message: "User created successfully",
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("There is an error in your net listen", err)
	}

	serv := grpc.NewServer()

	user.RegisterUserServiceServer(serv, &userService{})

	reflection.Register(serv)

	if err := serv.Serve(listen); err != nil {
		log.Fatal("error running server", err)
	}
}
