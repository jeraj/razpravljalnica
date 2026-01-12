package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/fatih/color"
	
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start interactive shell",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Connected to Razpravljalnica server!")
		
        
        blueBanner := color.New(color.FgCyan, color.Bold)
        fmt.Println() 
        blueBanner.Println("+--------------------------------------------+")
        blueBanner.Println("|         WELCOME TO RAZPRAVLJALNICA         |")
        blueBanner.Println("+--------------------------------------------+")
        fmt.Println()
        fmt.Println("Type '--help' to see the available commands!\nType 'exit' to quit\n\nPlease login to use the Razpravljalnica by using the command login <username>!\n")
        
		reader := bufio.NewReader(os.Stdin)

		for {
		    blue := color.New(color.FgBlue, color.Bold)
			blue.Print("\nrazpravljalnica>> ")
			line, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if line == "exit" {
                alert := color.New(color.FgRed, color.Bold)
                    
                fmt.Println()
				alert.Println("+--------------------------------------------+")
                alert.Println("|                  GOODBYE!!                 |")
                alert.Println("+--------------------------------------------+")
                fmt.Println()
				return nil
			}

			
			args := strings.Split(line, " ")
			rootCmd.SetArgs(args)

			if err := rootCmd.Execute(); err != nil {
				fmt.Println("Error:", err)
			}
		}
	},
}

