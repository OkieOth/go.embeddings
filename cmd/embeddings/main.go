package main

import (
	"fmt"
	"okieoth/goembeddings/cmd/embeddings/sub"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "embeddings",
	Short: "Tool to generate and query embeddings",
	Long:  `A tool to support creating and investigating embeddings`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmdToRun, _, err := sub.RunInteractive(cmd); err != nil {
			fmt.Println("Error while running program interactive:", err)
		} else {
			cmdToRun.Run(cmdToRun, args)
		}
	},
}

func init() {
	rootCmd.AddCommand(sub.CreateCmd)
}

func main() {
	rootCmd.Execute()
}
