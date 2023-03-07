package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getMagicBlock = &cobra.Command{
	Use:   "get-magicblock",
	Short: "Downloads the magic block",
	Long:  "Downloads the latest magic block created",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Get magic block called")
	},
}

func init() {
	rootCmd.AddCommand(getMagicBlock)
}
