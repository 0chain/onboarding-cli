package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"onboarding-cli/config"
	"onboarding-cli/util"

	"github.com/spf13/cobra"
)

var getMagicBlock = &cobra.Command{
	Use:   "get-magicblock",
	Short: "Downloads the magic block",
	Long:  "Downloads the latest magic block created",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Get magic block called")

		server_url, err := config.Extract()
		if err != nil {
			log.Fatal(err)
		}

		getReq, err := util.NewHTTPGetRequest(server_url + "magicblock")
		if err != nil {
			panic(err)
		}
		getResponse, err := getReq.Get()
		if err != nil {
			log.Fatal(err)
		}

		_, err = getResponse.CheckStatusCode()
		if err != nil {
			log.Fatal(err)
		}

		respBody := getResponse.PostResponse.Body

		path := "b0magicBlock.json"

		if err := ioutil.WriteFile(path, []byte(respBody), 0644); err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(getMagicBlock)
}
