package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/0chain/gosdk/core/common"
	"github.com/spf13/cobra"
	"math"
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

		magicBlockNumber := 1
		startingRound := 0
		t_percent := 66
		k_percent := 75
		N := len(nodes.Miners)

		T := int(math.Ceil(float64(N) * (float64(t_percent) / 100.0)))
		K := int(math.Ceil(float64(N) * (float64(k_percent) / 100.0)))

		fmt.Println(magicBlockNumber, startingRound, T, K)

		var minersmpks []types.MinerMpks

		for _, v := range nodes.Miners {
			minmpk := types.MinerMpks{
				Miner:        v,
				CreationDate: common.Now(),
				Type:         "Miner",
			}
			minersmpks = append(minersmpks, minmpk)
		}

	},
}

func init() {
	rootCmd.AddCommand(setupMPKS)
}
