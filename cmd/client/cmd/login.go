package cmd

import (
	"context"
	//"fmt"
	"time"
    
	pb "github.com/jeraj/razpravljalnica/gen"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"github.com/fatih/color"
)

var loginCmd = &cobra.Command{
    Use:   "login [username]",
    Short: "Register or login as a user",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        username := args[0]
        yellow := color.New(color.FgYellow, color.Bold)
        conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
        if err != nil {
            yellow.Println("Cannot connect:", err)
            return
        }

        grpcClient = pb.NewMessageBoardClient(conn)

        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        user, err := grpcClient.CreateUser(ctx, &pb.CreateUserRequest{Name: username})
        if err != nil {
            yellow.Println("Login failed:", err)
            return
        }

        currentUser = user
        yellow.Printf("Logged in as %s (id=%d)\n", user.Name, user.Id)
    },
}


