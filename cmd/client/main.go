package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/yishak-cs/CleanGrpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// establish a connection to the gRPC server
	connection, err := grpc.Dial("localhost:50000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to server: ", err)
	}
	defer connection.Close()

	// et client of UserService or the stub
	client := pb.NewUserServiceClient(connection)

	//context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 4 {
			fmt.Println("Usage: client create <name> <email>")
			return
		}
		createUser(ctx, client, os.Args[2], os.Args[3])

	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Usage: client get <user_id>")
			return
		}
		getUser(ctx, client, os.Args[2])

	case "list":
		listUsers(ctx, client)

	case "update":
		if len(os.Args) < 5 {
			fmt.Println("Usage: client update <user_id> <name> <email>")
			return
		}
		id, err := strconv.ParseUint(os.Args[2], 10, 32)
		if err != nil {
			fmt.Println("Invalid user ID:", err)
			return
		}
		updateUser(ctx, client, uint32(id), os.Args[3], os.Args[4])

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: client delete <user_id>")
			return
		}
		deleteUser(ctx, client, os.Args[2])

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  client create <name> <email>")
	fmt.Println("  client get <user_id>")
	fmt.Println("  client list")
	fmt.Println("  client update <user_id> <name> <email>")
	fmt.Println("  client delete <user_id>")
}

func createUser(ctx context.Context, client pb.UserServiceClient, name, email string) {
	req := &pb.CreateUserRequest{
		Name:  name,
		Email: email,
	}

	resp, err := client.CreateUser(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("Response: %s\n", resp.Status)
}

func getUser(ctx context.Context, client pb.UserServiceClient, id string) {
	req := &pb.SingleUserRequest{
		Id: id,
	}

	user, err := client.GetUser(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	fmt.Printf("User ID: %s\n", user.Id)
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("Email: %s\n", user.Email)
}

func listUsers(ctx context.Context, client pb.UserServiceClient) {
	resp, err := client.GetUsersList(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}

	fmt.Printf("Total users: %d\n", len(resp.Users))
	for i, user := range resp.Users {
		fmt.Printf("\nUser #%d:\n", i+1)
		fmt.Printf("  ID: %s\n", user.Id)
		fmt.Printf("  Name: %s\n", user.Name)
		fmt.Printf("  Email: %s\n", user.Email)
	}
}

func updateUser(ctx context.Context, client pb.UserServiceClient, id uint32, name, email string) {
	req := &pb.UpdateUserRequest{
		Id:    int64(id),
		Name:  name,
		Email: email,
	}

	resp, err := client.UpdateUser(ctx, req)
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

	fmt.Printf("Response: %s\n", resp.Status)
}

func deleteUser(ctx context.Context, client pb.UserServiceClient, id string) {
	req := &pb.SingleUserRequest{
		Id: id,
	}

	resp, err := client.DeleteUser(ctx, req)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}

	fmt.Printf("Response: %s\n", resp.Status)
}
