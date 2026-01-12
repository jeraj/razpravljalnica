package cmd

import (
	"context"
	"fmt"
	"time"

	//pb "github.com/jeraj/razpravljalnica/gen"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
	"github.com/fatih/color"
)

var topicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "List all topics",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
        magenta := color.New(color.FgMagenta, color.Bold)
		resp, err := grpcClient.ListTopics(ctx, &emptypb.Empty{})
		if err != nil {
			return err
		}

		if len(resp.Topics) == 0 {
			fmt.Println("No topics yet.")
			return nil
		}

		magenta.Println("\nTopics:")
		for _, t := range resp.Topics {
			fmt.Printf("  [%d] %s\n", t.Id, t.Name)
		}
		fmt.Printf("\n");
		return nil
	},
}


