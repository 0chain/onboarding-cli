package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateKeys = &cobra.Command{
	Use:   "generate-keys",
	Short: "Generates the keys and nodes structure",
	Long:  "Responsible or generating the keys and the node structures for miners and sharders",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generate Keys called")
	},
}

func init() {
	rootCmd.AddCommand(generateKeys)
}
