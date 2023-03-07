package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupMPKS = &cobra.Command{
	Use:   "setup-mpks",
	Short: "Generates mpks and shares and sends them",
	Long:  "Generates the mpks, and responsible for sharding and sending them",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setup MPKS called")
	},
}

func init() {
	rootCmd.AddCommand(setupMPKS)
}
