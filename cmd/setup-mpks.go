package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"onboarding-cli/types"
	"onboarding-cli/util"
)

var setupMPKS = &cobra.Command{
	Use:   "setup-mpks",
	Short: "Generates mpks and shares and sends them",
	Long:  "Generates the mpks, and responsible for sharding and sending them",
	Run: func(cmd *cobra.Command, args []string) {

		getReq, err := util.NewHTTPGetRequest("http://localhost:3000/nodes")
		if err != nil {
			panic(err)
		}
		getResponse, err := getReq.Get()
		respBody := getResponse.PostResponse.Body
		var nodes types.Nodes
		err = json.Unmarshal([]byte(respBody), &nodes)
		if err != nil {
			panic(err)
		}

		fmt.Println(nodes)

	},
}

func init() {
	rootCmd.AddCommand(setupMPKS)
}
