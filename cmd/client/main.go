package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
    "strconv"

    pb "github.com/jeraj/razpravljalnica/gen"
    "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
)



func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewMessageBoardClient(conn)

    var currentTopicID int64 = 0 //izbrana še ni bila nobena tema
    topics := make(map[int64]string)
    lastMessageID := make(map[int64]int64)
    //preberi uporabnisko ime od terminala
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Connected to server!\nSelect your username to register: ")
    username, _ := reader.ReadString('\n')
    username = username[:len(username)-1]

    ctx, cancel := context.WithTimeout(context.Background(), time.Second) //tukaj omejim čas, da odejmalec ne bi čakal predolgo
    //context je mehanizem za časovne omejitve, pomembna za življensko dobo grpc klica
    defer cancel()

    user, err := client.CreateUser(ctx, &pb.CreateUserRequest{Name: username})
    if err != nil {
        log.Fatalf("CreateUser failed: %v", err)
    }

    fmt.Printf("User registered successfully! Welcome, %s (id=%d)\n\n", user.Name, user.Id)
    fmt.Println("\n<---------------------------------------------COMMANDS--------------------------------------------->\n")
    fmt.Println("    v ... view other topics")
    fmt.Println("    c ... create topic")
    fmt.Println("    o <msgId> ... open topic with id <msgId>")
    fmt.Println("    p <text> ... post a comment in current topic")
    fmt.Println("    d <msgId> ... delete a message with id <msgIdid>")
    fmt.Println("    l <msgId> ... like a message with id <msgId>")
    fmt.Println("    e <msgId> ... update a message with id <msgId>\n")

    fmt.Println("\n<-------------------------------------------------------------------------------------------------->\n")


    for {
        fmt.Print(">> ")
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }

        line = strings.TrimSpace(line)
        //fmt.Println("You typed:", line)

        if line == "exit" {
            fmt.Println("Goodbye!")
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
                topics[t.Id] = t.Name
            }
        }



        if line == "c" { //naredimo temo
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
        
        

        izbira_teme := strings.Split(line, " ")
        if izbira_teme[0] == "o"{ //ogled sporočil v temi in izbira teme hkrati
            topicID, err := strconv.ParseInt(izbira_teme[1], 10, 64)
            if err != nil {
                fmt.Println("Invalid topic id")
                continue
            }

            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()

            resp, err := client.GetMessages(ctx, &pb.GetMessagesRequest{
                TopicId:        topicID,
                FromMessageId:  0,
                Limit:          50, //limit sporocil znotraj teme
            })
            if err != nil {
                log.Println("GetMessages failed:", err)
                continue
            }

            currentTopicID = topicID //nastavimo na trenutno temo
            tema := topics[topicID] //pogledamo v slovar njeno ime, katera temaje
            fmt.Printf("\n--- Topic %d --- (%s)\n", currentTopicID, tema)

            if len(resp.Messages) == 0 {
                fmt.Println("No messages in this topic.")
                continue
            }

            for _, m := range resp.Messages {
                fmt.Printf("[%d] user=%d: %s (%d likes)\n", m.Id, m.UserId, m.Text, m.Likes)
            }
        }
        
        
        if strings.HasPrefix(line, "p ") && currentTopicID != 0 {
            text := strings.TrimPrefix(line, "p ")
            text = strings.TrimSpace(text)

            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()

            msg, err := client.PostMessage(ctx, &pb.PostMessageRequest{
                TopicId: currentTopicID,
                UserId:  user.Id,
                Text:    text,
            })
            
            if err != nil {
                log.Println("PostMessage failed:", err)
                continue
            }

            fmt.Printf("Message posted: [%d] user=%d: %s\n", msg.Id, msg.UserId, msg.Text)
            continue
        }
        
        
        
        parts := strings.SplitN(line, " ", 2)

        if parts[0] == "l" {
	        if currentTopicID == 0 {
		        fmt.Println("No topic selected")
		        continue
	        }

	        msgID, err := strconv.ParseInt(parts[1], 10, 64)
	        if err != nil {
		        fmt.Println("Invalid message id")
		        
		        continue
	        }

	        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	        defer cancel()

	        msg, err := client.LikeMessage(ctx, &pb.LikeMessageRequest{
		        TopicId:   currentTopicID,
		        MessageId: msgID,
		        UserId:    user.Id,
	        })
	        if err != nil {
		        //log.Println("LikeMessage failed:", err)
		        fmt.Println("You already liked this message!");
		        continue
	        }

	        fmt.Printf("Message [%d] now has %d likes\n", msg.Id, msg.Likes)
        }
        
        if parts[0] == "d" {
	        if currentTopicID == 0 {
		        fmt.Println("No topic selected")
		        continue
	        }

	        msgID, err := strconv.ParseInt(parts[1], 10, 64)
	        if err != nil {
		        fmt.Println("Invalid message id")
		        continue
	        }

	        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	        defer cancel()

	        _, err = client.DeleteMessage(ctx, &pb.DeleteMessageRequest{
		        TopicId:   currentTopicID,
		        UserId:    user.Id,
		        MessageId: msgID,
	        })
	        if err != nil {
		        st, ok := status.FromError(err)
		        if ok && st.Code() == codes.PermissionDenied {
			        fmt.Println("You can only delete your own messages")
		        } else {
			        log.Println("DeleteMessage failed:", err)
		        }
		        continue
	        }

	        fmt.Println("Message deleted")
        }
        
        if parts[0] == "e" {
	        if currentTopicID == 0 {
		        fmt.Println("No topic selected")
		        continue
	        }

	        args := strings.SplitN(line, " ", 3)
	        if len(args) < 3 {
		        fmt.Println("Usage: -e <message_id> <new text>")
		        continue
	        }

	        msgID, err := strconv.ParseInt(args[1], 10, 64)
	        if err != nil {
		        fmt.Println("Invalid message id")
		        continue
	        }

	        newText := args[2]

	        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	        defer cancel()

	        msg, err := client.UpdateMessage(ctx, &pb.UpdateMessageRequest{
		        TopicId:   currentTopicID,
		        UserId:    user.Id,
		        MessageId: msgID,
		        Text:      newText,
	        })
	        if err != nil {
		        st, ok := status.FromError(err)
		        if ok && st.Code() == codes.PermissionDenied {
			        fmt.Println("You can only update your own messages")
		        } else {
			        log.Println("UpdateMessage failed:", err)
		        }
		        continue
	        }

	        fmt.Printf("Message updated: [%d] %s\n", msg.Id, msg.Text)
        }
        if strings.HasPrefix(line, "s ") {
	        args := strings.Split(line, " ")
	        if len(args) != 2 {
		        fmt.Println("Usage: s <topicID>")
		        continue
	        }

	        topicID, err := strconv.ParseInt(args[1], 10, 64)
	        if err != nil {
		        fmt.Println("Invalid topic id")
		        continue
	        }

	        if _, ok := topics[topicID]; !ok {
		        fmt.Println("Topic does not exist")
		        continue
	        }

	        fmt.Printf("Subscribed to topic %d (%s)\n", topicID, topics[topicID])
	        go subscribeToTopic(client, topicID, user.Id, lastMessageID)

        }



    }
}

func subscribeToTopic(client pb.MessageBoardClient, topicID, userID int64, lastMessageID map[int64]int64) {
    fromMessageID := lastMessageID[topicID]

    req := &pb.SubscribeTopicRequest{
        TopicId:        []int64{topicID},
        UserId:         userID,
        FromMessageId:  fromMessageID,
        SubscribeToken: "dummy-token",
    }

    stream, err := client.SubscribeTopic(context.Background(), req)
    if err != nil {
        log.Println("SubscribeTopic failed:", err)
        return
    }

    for {
        msgEvent, err := stream.Recv()
        if err != nil {
            log.Println("Subscription ended:", err)
            return
        }
        fmt.Printf("\n[SUBSCRIBE] Topic %d: [%d] user=%d: %s (%d likes)\n",
            msgEvent.Message.TopicId,
            msgEvent.Message.Id,
            msgEvent.Message.UserId,
            msgEvent.Message.Text,
            msgEvent.Message.Likes,
        )

        // posodobi lastMessageID za ta topic
        lastMessageID[topicID] = msgEvent.Message.Id
    }
}


/*
package main
var currentTopicID int64 = 0
var lastMessageID int64 = 0


func startMessagePolling() {
    go func() {
        for {
            if currentTopicID == 0 {
                time.Sleep(1 * time.Second)
                continue
            }

            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            resp, err := grpcClient.GetMessages(ctx, &pb.GetMessagesRequest{
                TopicId:       currentTopicID,
                FromMessageId: lastMessageID, // pridobi samo nove
                Limit:         50,
            })
            cancel()

            if err != nil {
                fmt.Println("Polling error:", err)
                time.Sleep(1 * time.Second)
                continue
            }

            for _, m := range resp.Messages {
                fmt.Printf("\n>> [%d] user=%d: %s (%d likes)\n", m.Id, m.UserId, m.Text, m.Likes)
                lastMessageID = m.Id
            }

            time.Sleep(1 * time.Second)
        }
    }()
}




func main() {
	Execute()
}*/

