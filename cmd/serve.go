package cmd

import (
	"fmt"
	"github.com/geropl/gemini-mcp-go/pkg/server"
	"github.com/spf13/cobra"
	"os"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the Gemini MCP server",
	Long:  `Starts the Gemini MCP server, listening for requests from an MCP client.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("GEMINI_API_KEY environment variable must be set")
		}
		return server.Run(cmd.Context(), apiKey)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
