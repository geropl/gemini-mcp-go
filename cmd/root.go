package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gemini-mcp-go",
	Short: "A Go implementation of an MCP server for Google Gemini",
	Long: `gemini-mcp-go is a server that implements the Model Context Protocol (MCP)
to act as a bridge between an MCP-compatible client and Google's Gemini API.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
