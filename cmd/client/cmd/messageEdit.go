package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	pb "github.com/jeraj/razpravljalnica/gen"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/fatih/color"
)

var messageEditCmd = &cobra.Command{
	Use:   "edit <message_id> <new text>",
	Short: "Edit your message",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		if currentUser == nil {
			return fmt.Errorf("you must login first")
		}
		if currentTopicID == 0 {
			return fmt.Errorf("no topic selected")
		}

		msgID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid message id")
		}

		newText := strings.Join(args[1:], " ")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		msg, err := grpcClient.UpdateMessage(ctx, &pb.UpdateMessageRequest{
			TopicId:   currentTopicID,
			UserId:    currentUser.Id,
			MessageId: msgID,
			Text:      newText,
		})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.PermissionDenied {
				return fmt.Errorf("you can only edit your own messages")
			}
			return err
		}
        magenta := color.New(color.FgMagenta, color.Bold)
		magenta.Printf("Message [%d] updated\n", msg.Id)
		return nil
	},
}

