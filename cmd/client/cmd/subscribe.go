package cmd

import (
	"context"
	"fmt"
    "strconv"
	pb "github.com/jeraj/razpravljalnica/gen"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
)
var subscriptions = make(map[int64]context.CancelFunc)


var subscribeCmd = &cobra.Command{
	Use:   "subscribe <topic_id>",
	Short: "Subscribe to a topic to receive new messages",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if currentUser == nil {
			return fmt.Errorf("you must login first")
		}

		topicID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid topic id")
		}

		if _, ok := subscriptions[topicID]; ok {
			return fmt.Errorf("already subscribed to topic %d", topicID)
		}

		_, cancel := context.WithCancel(context.Background())
		subscriptions[topicID] = cancel

        go func() {
            fromMessageID := lastMessageID[topicID]
            req := &pb.SubscribeTopicRequest{
                TopicId:       []int64{topicID},
                UserId:        currentUser.Id,
                FromMessageId: fromMessageID,
                SubscribeToken: "dummy-token",
            }

            stream, err := grpcClient.SubscribeTopic(context.Background(), req)
            if err != nil {
                fmt.Println("Subscribe failed:", err)
                return
            }

            for {
                msgEvent, err := stream.Recv()
                if err != nil {
                    fmt.Println("Subscription ended:", err)
                    return
                }

                
                color.Magenta("\n[NOVO SPOROČILO - Tema %d]", msgEvent.Message.TopicId)
                printFormattedMessage(msgEvent.Message, currentUser.Id)
                
                lastMessageID[topicID] = msgEvent.Message.Id
            }
        }()

        magenta := color.New(color.FgMagenta, color.Bold)
		magenta.Printf("Subscribed to topic %d\n", topicID)
		return nil
	},
}

