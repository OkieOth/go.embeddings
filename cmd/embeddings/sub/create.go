package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates embeddings",
	Long:  `Create embeddings and stores them in a database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: implement embeddings create")
	},
}
