package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"onboarding-cli/config"
	"onboarding-cli/types"
	"onboarding-cli/util"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var getInitialStatesCmd = &cobra.Command{
	Use:   "get-initialstates",
	Short: "Downloads the initial states",
	Long:  "Downloads the latest initial states created",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting initial states")
		server_url, err := config.Extract()
		if err != nil {
			log.Fatal(err)
		}

		getReq, err := util.NewHTTPGetRequest(server_url + "initialstates")
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

		data := make(map[string][]*types.InitialStateData)

		err = json.Unmarshal([]byte(respBody), &data)

		if err != nil {
			log.Fatal(err)
		}

		yamlData, err := yaml.Marshal(data)

		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(data["initial_states"])
		path := "initial_states.yaml"

		if err := ioutil.WriteFile(path, yamlData, 0644); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getInitialStatesCmd)
}
