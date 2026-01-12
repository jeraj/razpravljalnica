package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/jeraj/razpravljalnica/gen"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

var messagePostCmd = &cobra.Command{
	Use:   "post <text>",
	Short: "Post a message to current topic",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if currentUser == nil {
			return fmt.Errorf("you must login first")
		}
		if currentTopicID == 0 {
			return fmt.Errorf("no topic selected")
		}

		text := strings.Join(args, " ")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		msg, err := grpcClient.PostMessage(ctx, &pb.PostMessageRequest{
			TopicId: currentTopicID,
			UserId:  currentUser.Id,
			Text:    text,
		})
		if err != nil {
			return err
		}
        magenta := color.New(color.FgMagenta, color.Bold)
        magenta.Printf("Sporočilo objavljeno z ID: %d\n", msg.Id)
		return nil
	},
}

