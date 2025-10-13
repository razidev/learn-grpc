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
	"strings"
	"time"

	protovalidate "buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userService struct {
	user.UnimplementedUserServiceServer
}

func loggingMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("Logging middleware")
	log.Println(info.FullMethod)
	res, err := handler(ctx, req)

	log.Println("Setelah request")
	return res, err
}

func authMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("Auth middleware")

	if info.FullMethod == "/user.UserService/Login" {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authToken, ok := md["authorization"]
	if !ok || len(authToken) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	log.Println("Token:", authToken[0])

	splitToken := strings.Split(authToken[0], " ")
	token := splitToken[1]

	if token != "AccessToken" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}

func (us *userService) Login(ctx context.Context, loginRequest *user.LoginRequest) (*user.LoginResponse, error) {
	return &user.LoginResponse{
		Base: &common.BaseResponse{
			StatusCode: 200,
			IsSuccess:  true,
			Message:    "Login successful",
		},
		AccessToken:  "AccessToken",
		RefreshToken: "RefreshToken",
	}, nil
}

func (us *userService) CreateUser(ctx context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	if err := protovalidate.Validate(userRequest); err != nil {
		if ve, ok := err.(*protovalidate.ValidationError); ok {
			var validations []*common.ValidationError = make([]*common.ValidationError, 0)
			for _, fieldErr := range ve.Violations {
				log.Printf("Field %s: %s", *fieldErr.Proto.Field.Elements[0].FieldName, *fieldErr.Proto.Message)
				validations = append(validations, &common.ValidationError{
					Field:   *fieldErr.Proto.Field.Elements[0].FieldName,
					Message: *fieldErr.Proto.Message,
				})
			}

			return &user.CreateResponse{
				Base: &common.BaseResponse{
					ValidationErrors: validations,
					StatusCode:       400,
					IsSuccess:        false,
					Message:          "Validation error",
				},
			}, nil
		}
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	return &user.CreateResponse{
		Base: &common.BaseResponse{
			StatusCode: 200,
			IsSuccess:  true,
			Message:    "User created successfully",
		},
		CreatedAt: timestamppb.Now(),
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

	serv := grpc.NewServer(grpc.ChainUnaryInterceptor(loggingMiddleware, authMiddleware))

	user.RegisterUserServiceServer(serv, &userService{})
	chat.RegisterChatServiceServer(serv, &chatService{})

	reflection.Register(serv)

	if err := serv.Serve(listen); err != nil {
		log.Fatal("error running server", err)
	}
}
