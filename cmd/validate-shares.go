package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateShares = &cobra.Command{
	Use:   "validate-shares",
	Short: "Validates shares, creates sos, sends it, creates dkg local file",
	Long:  "Validates the shares, creates signatures or shares, sends them and then creates a dkg local file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Validate Shares called")
	},
}

func init() {
	rootCmd.AddCommand(setupMPKS)
}
