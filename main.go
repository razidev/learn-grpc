package main

import (
	"context"
	"errors"
	"grpc/pb/chat"
	"grpc/pb/user"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

type chatService struct {
	chat.UnimplementedChatServiceServer
}

func (cs *chatService) SendMessage(stream grpc.ClientStreamingServer[chat.ChatMessage, chat.ChatResponse]) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return status.Errorf(codes.Unknown, "Error receiving message: %v", err)
		}

		log.Printf("Received message from user %d: %s", req.UserId, req.Content)
	}

	return stream.SendAndClose(&chat.ChatResponse{
		Message: "Message received successfully",
	})
}

// func (UnimplementedChatServiceServer) ReceiveMessage(*ReceiveMessageRequest, grpc.ServerStreamingServer[ChatMessage]) error {
// 	return status.Errorf(codes.Unimplemented, "method ReceiveMessage not implemented")
// }
// func (UnimplementedChatServiceServer) Chat(grpc.BidiStreamingServer[ChatMessage, ChatMessage]) error {
// 	return status.Errorf(codes.Unimplemented, "method Chat not implemented")
// }

func main() {
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("There is an error in your net listen", err)
	}

	serv := grpc.NewServer()

	user.RegisterUserServiceServer(serv, &userService{})
	chat.RegisterChatServiceServer(serv, &chatService{})

	reflection.Register(serv)

	if err := serv.Serve(listen); err != nil {
		log.Fatal("error running server", err)
	}
}
