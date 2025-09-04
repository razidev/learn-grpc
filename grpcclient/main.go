package main

import (
	"context"
	"errors"
	"grpc/pb/chat"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	clientConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("There is an error in your grpc dial", err)
	}

	/* For User Service
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
	*/

	/* For Chat Service */
	chatClient := chat.NewChatServiceClient(clientConn)

	/*
		stream, err := chatClient.SendMessage(context.Background())
		if err != nil {
			log.Fatal("There is an error in your send message", err)
		}

		err = stream.Send(&chat.ChatMessage{
			UserId:  1,
			Content: "Hello, this is a test message",
		})
		if err != nil {
			log.Fatal("There is an error in your sending message", err)
		}
		err = stream.Send(&chat.ChatMessage{
			UserId:  1,
			Content: "Hello, again",
		})
		if err != nil {
			log.Fatal("There is an error in your sending message", err)
		}

		time.Sleep(5 * time.Second)
		err = stream.Send(&chat.ChatMessage{
			UserId:  1,
			Content: "Hello, after 5 seconds",
		})
		if err != nil {
			log.Fatal("There is an error in your sending message", err)
		}

		res, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatal("There is an error in your receiving response", err)
		}

		log.Println("Response from server:", res.Message)
	*/

	stream, err := chatClient.ReceiveMessage(context.Background(), &chat.ReceiveMessageRequest{
		UserId: 22,
	})
	if err != nil {
		log.Fatal("There is an error in your receive message", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal("There is an error in your receiving message", err)
		}

		log.Printf("Received message from user %d: %s", msg.UserId, msg.Content)
	}
}
