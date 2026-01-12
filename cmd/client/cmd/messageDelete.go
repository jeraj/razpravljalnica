package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/jeraj/razpravljalnica/gen"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/fatih/color"
)

var messageDeleteCmd = &cobra.Command{
	Use:   "delete <message_id>",
	Short: "Delete your message",
	Args:  cobra.ExactArgs(1),
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err = grpcClient.DeleteMessage(ctx, &pb.DeleteMessageRequest{
			TopicId:   currentTopicID,
			UserId:    currentUser.Id,
			MessageId: msgID,
		})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.PermissionDenied {
				return fmt.Errorf("you can only delete your own messages")
			}
			return err
		}
        magenta := color.New(color.FgMagenta, color.Bold)
		magenta.Println("Message deleted")
		return nil
	},
}

