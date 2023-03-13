package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"onboarding-cli/core"
	"onboarding-cli/types"
	"onboarding-cli/util"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var setupMPKS = &cobra.Command{
	Use:   "send-shares",
	Short: "Generates shares and sends them",
	Long:  "Generates the mpks, and responsible for sharding and sending them",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Send Shares called")

		getReq, err := util.NewHTTPGetRequest("http://localhost:3000/nodes")
		if err != nil {
			panic(err)
		}
		getResponse, err := getReq.Get()
		if err != nil {
			log.Fatal(err)
		}
		respBody := getResponse.PostResponse.Body
		var nodes types.Nodes
		err = json.Unmarshal([]byte(respBody), &nodes)
		if err != nil {
			panic(err)
		}

		minerIds := make([]string, len(nodes.Miners))

		for i := 0; i < len(nodes.Miners); i++ {
			minerIds[i] = nodes.Miners[i].ID
		}

		yfile, err := ioutil.ReadFile("nodes.yml")

		if err != nil {
			log.Fatal(err)
		}

		data := make(map[string]any)

		err = yaml.Unmarshal(yfile, &data)

		if err != nil {
			log.Fatal(err)
		}

		var miners []*struct {
			ID       string
			SetIndex int
		}

		for _, miner := range data["miners"].([]interface{}) {

			minerData := miner.(map[string]any)

			id := minerData["id"]
			setIndex := minerData["set_index"]

			miner := &struct {
				ID       string
				SetIndex int
			}{id.(string), int(setIndex.(int))}

			miners = append(miners, miner)
		}

		shares := make([]*types.ShareData, 0)

		for _, miner := range miners {
			shares = append(shares, core.CreateShares(minerIds, miner.SetIndex,
				miner.ID)...)
		}

		postReq, err := util.NewHTTPPostRequest("http://localhost:3000/mpks", shares)

		if err != nil {
			log.Fatal(err)
		}

		postResponse, err := postReq.Post()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Post Request Response", postResponse)

	},
}

func init() {
	rootCmd.AddCommand(setupMPKS)
}
