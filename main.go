package main

import (
	"context"
	"errors"
	"fmt"
	"grpc/pb/chat"
	"grpc/pb/common"
	"grpc/pb/user"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type userService struct {
	user.UnimplementedUserServiceServer
}

func (us *userService) CreateUser(ctx context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	if userRequest.Age < 1 {
		return &user.CreateResponse{
			Base: &common.BaseResponse{
				StatusCode: 400,
				IsSuccess:  false,
				Message:    "Age must be greater than 0",
			},
		}, nil
	}

	return &user.CreateResponse{
		Base: &common.BaseResponse{
			StatusCode: 200,
			IsSuccess:  true,
			Message:    "User created successfully",
		},
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

func (cs *chatService) ReceiveMessage(req *chat.ReceiveMessageRequest, stream grpc.ServerStreamingServer[chat.ChatMessage]) error {
	log.Printf("ReceiveMessage called for user %d", req.UserId)

	for i := 0; i < 10; i++ {
		err := stream.Send(&chat.ChatMessage{
			UserId:  req.UserId,
			Content: fmt.Sprintf("Hello from server %d!", i),
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "Error sending message: %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}

func (cs *chatService) Chat(stream grpc.BidiStreamingServer[chat.ChatMessage, chat.ChatMessage]) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return status.Errorf(codes.Unknown, "Error receiving message: %v", err)
		}

		log.Printf("Received message from user %d: %s", msg.UserId, msg.Content)

		time.Sleep(2 * time.Second)

		err = stream.Send(&chat.ChatMessage{
			UserId:  1,
			Content: "Echo from server",
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "Error sending message: %v", err)
		}

		err = stream.Send(&chat.ChatMessage{
			UserId:  1,
			Content: "Echo from server #2",
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "Error sending message: %v", err)
		}
	}
	return nil
}

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
