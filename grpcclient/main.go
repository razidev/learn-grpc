package main

import (
	"context"
	"grpc/pb/user"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	clientConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("There is an error in your grpc dial", err)
	}

	/* For User Service */
	userClient := user.NewUserServiceClient(clientConn)
	resp, err := userClient.CreateUser(context.Background(), &user.User{
		Id:      1,
		Age:     -2,
		Balance: 130000,
		Address: &user.Address{
			Id:          1,
			FullAddress: "123 Main St, City, Country",
			Province:    "ProvinceName",
			City:        "CityName",
		},
		BirthDate: timestamppb.New(time.Now()),
	})
	if err != nil {
		// st, ok := status.FromError(err)
		// if ok {
		// 	if st.Code() == codes.InvalidArgument {
		// 		log.Println("Invalid argument error:", st.Message())
		// 	} else if st.Code() == codes.Unknown {
		// 		log.Println("Unknown error:", st.Message())
		// 	} else if st.Code() == codes.Internal {
		// 		log.Println("Internal server error:", st.Message())
		// 	}

		// 	return
		// }
		log.Println("There is an error in your create user", err)
		return
	}

	if !resp.Base.IsSuccess {
		switch resp.Base.StatusCode {
		case 400:
			log.Println("Client error:", resp.Base.Message)
		case 500:
			log.Println("Server error:", resp.Base.Message)
		}

		return
	}

	log.Println("Response from server:", resp.Base.Message)

	/* For Chat Service */
	// chatClient := chat.NewChatServiceClient(clientConn)

	/* client streaming
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

	/* server streaming
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
	*/

	/* bidirectional streaming
	stream, err := chatClient.Chat(context.Background())
	if err != nil {
		log.Fatal("There is an error in your chat", err)
	}

	err = stream.Send(&chat.ChatMessage{
		UserId:  1237,
		Content: "Hello, this is a bidirectional streaming message",
	})
	if err != nil {
		log.Fatal("There is an error in your sending message", err)
	}

	msg, err := stream.Recv()
	if err != nil {
		log.Fatal("There is an error in your receiving message", err)
	}
	log.Printf("Received message from user %d: %s", msg.UserId, msg.Content)
	msg, err = stream.Recv()
	if err != nil {
		log.Fatal("There is an error in your receiving message", err)
	}
	log.Printf("Received message from user %d: %s", msg.UserId, msg.Content)

	time.Sleep(5 * time.Second)

	err = stream.Send(&chat.ChatMessage{
		UserId:  1237,
		Content: "Hello, this is a bidirectional streaming message 1",
	})
	if err != nil {
		log.Fatal("There is an error in your sending message 1", err)
	}

	msg, err = stream.Recv()
	if err != nil {
		log.Fatal("There is an error in your receiving message", err)
	}
	log.Printf("Received message from user %d: %s", msg.UserId, msg.Content)
	msg, err = stream.Recv()
	if err != nil {
		log.Fatal("There is an error in your receiving message", err)
	}
	log.Printf("Received message from user %d: %s", msg.UserId, msg.Content)
	*/
}
