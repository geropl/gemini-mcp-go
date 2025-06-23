package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Helps with the initial setup of the Gemini MCP server",
	Long:  `Checks if the GEMINI_API_KEY environment variable is set and provides instructions if it is not.`,
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("GEMINI_API_KEY") == "" {
			fmt.Println("The GEMINI_API_KEY environment variable is not set.")
			fmt.Println("Please set it to your Gemini API key.")
			fmt.Println("You can get your API key from https://makersuite.google.com/app/apikey")
		} else {
			fmt.Println("The GEMINI_API_KEY environment variable is set.")
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
