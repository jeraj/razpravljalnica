package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

var topicUseCmd = &cobra.Command{
	Use:   "use <topicID>",
	Short: "Select active topic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid topic id")
		}

		currentTopicID = id
		magenta := color.New(color.FgMagenta, color.Bold)
		magenta.Printf("Using topic %d\n", currentTopicID)
		return nil
	},
}

