package cmd

import (
	"context"
	//"fmt"
	"time"

	pb "github.com/jeraj/razpravljalnica/gen"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

var topicCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new topic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
        magenta := color.New(color.FgMagenta, color.Bold)
		topic, err := grpcClient.CreateTopic(ctx, &pb.CreateTopicRequest{
			Name: args[0],
		})
		if err != nil {
			return err
		}

		magenta.Printf("Topic created: [%d] %s\n", topic.Id, topic.Name)
		return nil
	},
}

