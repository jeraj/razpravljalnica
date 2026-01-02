package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

    pb "github.com/jeraj/razpravljalnica/gen"
    "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewMessageBoardClient(conn)

    //var currentTopicID int64 = 0 //izbrana še ni bila nobena tema

    //preberi uporabnisko ime od terminala
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Select your username: ")
    username, _ := reader.ReadString('\n')
    username = username[:len(username)-1]

    ctx, cancel := context.WithTimeout(context.Background(), time.Second) //tukaj omejim čas, da odejmalec ne bi čakal predolgo
    //context je mehanizem za časovne omejitve, pomembna za življensko dobo grpc klica
    defer cancel()

    user, err := client.CreateUser(ctx, &pb.CreateUserRequest{Name: username})
    if err != nil {
        log.Fatalf("CreateUser failed: %v", err)
    }

    fmt.Println("User successfully created:", user)
    fmt.Println("\nCOMMANDS\n")
    fmt.Println("  v: view other topics")
    fmt.Println("  c : create topic")


    for {
        fmt.Print("> ")
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }

        line = strings.TrimSpace(line)
        fmt.Println("You typed:", line)

        if line == "exit" {
            fmt.Println("Bye")
            break
        }

        if line == "v" { //zelimo videti teme
            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()

            resp, err := client.ListTopics(ctx, &emptypb.Empty{})
            if err != nil {
                log.Println("ListTopics failed:", err)
                continue
            }

            //implementirano za vsak slucaj, kljub temu, da je ze narejena neka default tema
            if len(resp.Topics) == 0 {
                fmt.Println("No topics yet.")
                continue
            }

            fmt.Println("Topics:")
            for _, t := range resp.Topics { //izpis tem
                fmt.Printf(" - [%d] %s\n", t.Id, t.Name)
            }
        }

        if line == "c" {
            fmt.Print("Enter topic name: ")
            topicName, _ := reader.ReadString('\n')
            topicName = strings.TrimSpace(topicName)

            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()

            topic, err := client.CreateTopic(
                ctx,
                &pb.CreateTopicRequest{Name: topicName},
            )
            if err != nil {
                log.Println("CreateTopic failed:", err)
                continue
            }

            fmt.Printf("Topic created: [%d] %s\n", topic.Id, topic.Name)
        }
    }
}