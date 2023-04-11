package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"onboarding-cli/config"
	"onboarding-cli/core"
	"onboarding-cli/types"
	"onboarding-cli/util"

	"github.com/0chain/gosdk/core/encryption"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type ShareResp struct {
	Shares []*types.ShareData `json:"shares"`
}

var validateShares = &cobra.Command{
	Use:   "validate-shares",
	Short: "Validates shares, creates sos, sends it, creates dkg local file",
	Long:  "Validates the shares, creates signatures or shares, sends them and then creates a dkg local file",
	Run: func(cmd *cobra.Command, args []string) {
		//Get all the miners from the server
		server_url, err := config.Extract()
		if err != nil {
			log.Fatal(err)
		}
		getReq, err := util.NewHTTPGetRequest(server_url + "nodes")
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
		var nodes types.Nodes
		err = json.Unmarshal([]byte(respBody), &nodes)
		if err != nil {
			panic(err)
		}

		minerMap := make(map[string][]string)

		for _, node := range nodes.Miners {
			minerMap[node.ID] = node.MPK
		}

		//Get all local miners
		yfile, err := ioutil.ReadFile("nodes.yaml")

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
			PrivKey  string
			PubKey   string
		}

		for _, miner := range data["miners"].([]interface{}) {

			minerData := miner.(map[string]any)

			id := minerData["id"]
			setIndex := minerData["set_index"]
			privKey := minerData["private_key"]
			pubKey := minerData["public_key"]

			miner := &struct {
				ID       string
				SetIndex int
				PrivKey  string
				PubKey   string
			}{id.(string), int(setIndex.(int)), privKey.(string), pubKey.(string)}

			miners = append(miners, miner)
		}

		//Get all the shares from the server

		signs := make([]*types.SignData, 0)

		for _, miner := range miners {

			signs = append(signs, SendSignedMessages(miner.ID, miner.PrivKey, miner.SetIndex, miner.PubKey, minerMap)...)

		}

		// Send the signs to the server

		postReq, err := util.NewHTTPPostRequest(server_url+"sos", signs)
		if err != nil {
			log.Fatal(err)
		}

		postResponse, err := postReq.Post()
		if err != nil {
			log.Fatal(err)
		}
		_, err = postResponse.CheckStatusCode()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Post Request Response", postResponse)

	},
}

func init() {
	rootCmd.AddCommand(validateShares)
}

func SendSignedMessages(currId string, privKey string, setIndex int, pubKey string, minerMap map[string][]string) []*types.SignData {

	server_url, err := config.Extract()
	if err != nil {
		log.Fatal(err)
	}

	var privateKey bls.SecretKey

	privateKeyBytes, err := hex.DecodeString(privKey)

	if err != nil {
		log.Fatal(err)
	}

	if err := privateKey.SetLittleEndian(privateKeyBytes); err != nil {
		log.Fatal(err)
	}

	hashData := fmt.Sprintf("%v%v", currId, pubKey)
	hash := encryption.Hash(hashData)
	signature := privateKey.Sign(hash).SerializeToHexStr()

	getReq, err := util.NewHTTPGetRequest(server_url + "shares/" + currId)
	if err != nil {
		panic(err)
	}

	headers := map[string]string{"X-App-Miner-ID": currId, "X-App-Public-Key": pubKey, "X-App-Client-Signature": signature}

	getReq.SetHeaders(headers)
	getResponse, err := getReq.Get()
	if err != nil {
		log.Fatal(err)
	}

	_, err = getResponse.CheckStatusCode()

	if err != nil {
		log.Fatal(err)
	}

	respBody := getResponse.PostResponse.Body
	var shares types.ShareServer
	err = json.Unmarshal([]byte(respBody), &shares)
	if err != nil {
		panic(err)
	}

	shareMap := make(map[string]string)
	shareDataMap := make(map[string]string)

	for _, share := range shares.Shares {

		partyId := core.ComputeIDdkg(share.FromMiner)

		shareDataMap[partyId.GetHexString()] = share.Share
		shareMap[share.FromMiner] = share.Share
	}

	mp := map[string]any{
		"id":             "1",
		"starting_round": 0,
		"secret_shares":  shareDataMap,
	}

	if err := saveDKGSummary(mp, setIndex); err != nil {
		panic(err)
	}

	//Sign the shares

	data := core.SignMessages(shareMap, minerMap, privKey, currId)

	//Send the signed shares

	return data

}

func saveDKGSummary(dkg map[string]any, index int) error {

	path := getPath(index)

	var err error
	var dkgData []byte
	if dkgData, err = json.MarshalIndent(dkg, "", " "); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, dkgData, 0644); err != nil {
		return err
	}

	return nil

}

func getPath(index int) string {
	return fmt.Sprintf("dkgSummary-%v_dkg.json", index+1)
}
