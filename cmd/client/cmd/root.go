/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
    "fmt"
    "strings"
	pb "github.com/jeraj/razpravljalnica/gen"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"github.com/fatih/color"
)

var (
	grpcConn   *grpc.ClientConn
	grpcClient pb.MessageBoardClient
)

var (
	currentUser    *pb.User
	currentTopicID int64
)

//var subscriptions = make(map[int64]context.CancelFunc)
var lastMessageID = make(map[int64]int64)


// rootCmd represents the base command when called without any subcommands
/*var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}*/

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Razpravljalnica CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			return err
		}

		grpcConn = conn
		grpcClient = pb.NewMessageBoardClient(conn)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
/*func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}*/

func Execute() {
	err := rootCmd.Execute()
	if grpcConn != nil {
		grpcConn.Close()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.client.yaml)")

    magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()

    // 2. Nastavi HelpFunc, ki bo obarval izpis
    rootCmd.SetHelpFunc(func(c *cobra.Command, s []string) {
        // Pridobimo originalno predlogo
        t := c.UsageTemplate()
        
        // Zamenjamo ključne besede z obarvanimi
        t = strings.ReplaceAll(t, "Usage:", magenta("Usage:"))
        t = strings.ReplaceAll(t, "Available Commands:", magenta("Available Commands:"))
        t = strings.ReplaceAll(t, "Flags:", magenta("Flags:"))
        
        // Začasno nastavimo to pobarvano predlogo in izpišemo pomoč
        c.SetUsageTemplate(t)
        c.Usage()
    })

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    
    rootCmd.AddCommand(shellCmd)
    rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(topicsCmd)
	rootCmd.AddCommand(topicCmd)
	topicCmd.AddCommand(topicCreateCmd)
	topicCmd.AddCommand(topicUseCmd)
	rootCmd.AddCommand(messageCmd)
	messageCmd.AddCommand(messageListCmd)
    messageCmd.AddCommand(messagePostCmd)
    messageCmd.AddCommand(messageLikeCmd)
    messageCmd.AddCommand(messageDeleteCmd)
    messageCmd.AddCommand(messageEditCmd)
    rootCmd.AddCommand(subscribeCmd)
    rootCmd.AddCommand(unsubscribeCmd)

}

func formatMessage(m *pb.Message, currentID int64) string {

    colorReset  := "\033[0m"
    colorRed    := "\033[31m"
    colorGreen  := "\033[32m"
    colorWhite  := "\033[37m"
    colorYellow := "\033[33m"

    userDisplay := fmt.Sprintf("%d", m.UserId)
    userColor := colorGreen
    

    if m.UserId == currentID {
        userColor = colorRed
        userDisplay = "Ti"
    }

    return fmt.Sprintf("[%d] %s%s%s: %s%s%s %s❤ %d%s",
        m.Id, 
        userColor, userDisplay, colorReset, 
        colorWhite, m.Text, colorReset, 
        colorYellow, m.Likes, colorReset)
}




