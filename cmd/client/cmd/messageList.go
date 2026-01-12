package cmd

import (
	"context"
	"fmt"
	"time"

	pb "github.com/jeraj/razpravljalnica/gen"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

var messageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List messages in current topic",
	RunE: func(cmd *cobra.Command, args []string) error {

		if currentUser == nil {
			return fmt.Errorf("you must login first")
		}
		if currentTopicID == 0 {
			return fmt.Errorf("no topic selected (use: topic use <id>)")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := grpcClient.GetMessages(ctx, &pb.GetMessagesRequest{
			TopicId:       currentTopicID,
			FromMessageId: 0,
			Limit:         50,
		})
		if err != nil {
			return err
		}
		
        magenta := color.New(color.FgMagenta, color.Bold)
		if len(resp.Messages) == 0 {
			magenta.Println("No messages.")
			return nil
		}

        for _, m := range resp.Messages {
            printFormattedMessage(m, currentUser.Id)
        }
		return nil
	},
}


func printFormattedMessage(m *pb.Message, currentID int64) {

    white  := color.New(color.FgWhite)
    green  := color.New(color.FgGreen).Add(color.Bold)
    red    := color.New(color.FgRed).Add(color.Bold)
    yellow := color.New(color.FgYellow)


    fmt.Printf("[%d] ", m.Id)


    if m.UserId == currentID {
        red.Print("Ti")
    } else {
        green.Printf("User %d", m.UserId)
    }


    fmt.Print(": ")
    white.Print(m.Text)


    fmt.Print(" ")
    yellow.Printf("❤ %d\n", m.Likes)
}

