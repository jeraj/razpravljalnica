package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os"
    "time"

    pb "github.com/jeraj/razpravljalnica/gen"
    "google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewMessageBoardClient(conn)

    //preberi uporabnisko ime od terminala
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Select your username: ")
    username, _ := reader.ReadString('\n')
    username = username[:len(username)-1] // odstrani newline

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    user, err := client.CreateUser(ctx, &pb.CreateUserRequest{Name: username})
    if err != nil {
        log.Fatalf("CreateUser failed: %v", err)
    }

    fmt.Println("User successfully created:", user)
}