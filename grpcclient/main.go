package main

import (
	"context"
	"grpc/pb/user"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	clientConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("There is an error in your grpc dial", err)
	}

	userClient := user.NewUserServiceClient(clientConn)
	resp, err := userClient.CreateUser(context.Background(), &user.User{
		Id:      1,
		Age:     24,
		Balance: 130000,
		Address: &user.Address{
			Id:          1,
			FullAddress: "123 Main St, City, Country",
			Province:    "ProvinceName",
			City:        "CityName",
		},
	})
	if err != nil {
		log.Fatal("There is an error in your create user", err)
	}

	log.Println("Response from server:", resp.Message)

}
