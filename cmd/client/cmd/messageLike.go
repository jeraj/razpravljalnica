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

var messageLikeCmd = &cobra.Command{
	Use:   "like <message_id>",
	Short: "Like a message",
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

		msg, err := grpcClient.LikeMessage(ctx, &pb.LikeMessageRequest{
			TopicId:   currentTopicID,
			MessageId: msgID,
			UserId:    currentUser.Id,
		})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.AlreadyExists {
				return fmt.Errorf("you already liked this message")
			}
			return err
		}
        magenta := color.New(color.FgMagenta, color.Bold)
		magenta.Printf("Message [%d] now has %d likes\n", msg.Id, msg.Likes)
		return nil
	},
}

