package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

var unsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe <topic_id>",
	Short: "Unsubscribe from a topic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if currentUser == nil {
			return fmt.Errorf("you must login first")
		}

		topicID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid topic id")
		}

		cancel, ok := subscriptions[topicID]
		if !ok {
			return fmt.Errorf("not subscribed to topic %d", topicID)
		}

		cancel()                 // prekliči goroutino
		delete(subscriptions, topicID) // odstrani iz mape
		magenta := color.New(color.FgMagenta, color.Bold)
		magenta.Printf("Unsubscribed from topic %d\n", topicID)
		return nil
	},
}

