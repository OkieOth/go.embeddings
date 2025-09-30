package main

import (
	"fmt"
	"okieoth/schemaguesser/cmd/embeddings/sub"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "embeddings",
	Short: "Tool to generate and query embeddings",
	Long:  `A tool to support creating and investigating embeddings`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("embeddings - call the tool with one of the provided sub commands")
	},
}

func init() {
	rootCmd.AddCommand(sub.CreateCmd)
}

func main() {
	rootCmd.Execute()
}
