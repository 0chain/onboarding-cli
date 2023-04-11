package cmd

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"onboarding-cli/config"
	"onboarding-cli/core"
	"onboarding-cli/types"
	"onboarding-cli/util"
	"os"
	"strconv"

	"github.com/0chain/errors"
	"github.com/0chain/gosdk/core/encryption"
	"github.com/0chain/gosdk/core/zcncrypto"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var T = 2
var N = 3

var generateKeys = &cobra.Command{
	Use:   "generate-keys",
	Short: "Generates the keys and nodes structure",
	Long:  "Responsible or generating the keys and the node structures for miners and sharders",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating Keys...")
		var (
			flags = cmd.Flags()
			err   error
		)

		server_url, err := config.Extract()
		if err != nil {
			log.Fatal(err)
		}

		var clientSigScheme string

		clientSigScheme, err = flags.GetString("signature_scheme")
		if err != nil {
			log.Fatal(err)
		}
		if !(clientSigScheme == "bls0chain" || clientSigScheme == "ed25519") {
			log.Fatal("Signature Scheme can only take either bls0chain or ed25519")
		}

		overwrite, _ := flags.GetBool("overwrite")

		if !overwrite {

			if _, err := os.Stat("nodes.yaml"); err == nil || !os.IsNotExist(err) {
				// some nodes exists
				log.Fatal("nodes.yaml already exists. Aborting please delete the file or provide --overwrite flag and try again")
			}
		}

		var miners, sharders int

		yfile, err := ioutil.ReadFile("config.yaml")
		if err != nil {
			log.Fatal(err)
		}

		configData := make(map[string][]types.ConfigNodeData)
		err = yaml.Unmarshal(yfile, &configData)
		if err != nil {
			log.Fatal(err)
		}
		minersConfigData := configData["miners"]
		shardersConfigData := configData["sharders"]

		miners, err = flags.GetInt("miners")
		if err != nil {
			log.Fatal(err)
		}

		sharders, err = flags.GetInt("sharders")
		if err != nil {
			log.Fatal(err)
		}

		if len(minersConfigData) != miners {
			log.Fatal("Number of miners entered and miner data in config.yaml have mismatched length. Aborting")
			return
		}

		if len(shardersConfigData) != sharders {
			log.Fatal("Number of sharder entered and sharder data in config.yaml have mismatched length. Aborting")
			return
		}

		file, err := os.OpenFile("nodes.yaml", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		err = os.RemoveAll("keys")

		if err != nil {
			log.Fatal(err)
		}

		err = os.MkdirAll("keys", os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}
		minersData := ""
		if miners > 0 {
			minersData = "miners:\n"
			err = os.MkdirAll("output", os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}

		minerNodes := []types.Miner{}

		for i := 1; i <= miners; i++ {
			var wallet *zcncrypto.Wallet
			wallet, err = getWallet(clientSigScheme)
			if err != nil {
				panic(err)
			}
			minerNode, minerData, err := generateMinerNodeStructure(wallet, clientSigScheme, i, minersConfigData[i-1])
			if err != nil {
				panic(err)
			}
			minersData += minerData
			minerNodes = append(minerNodes, minerNode)
			path := fmt.Sprintf("keys/b0mnode%d_keys.json", i)
			err = saveWallet(path, wallet)
			if err != nil {
				log.Fatal(err)
			}
		}

		shardersData := "sharders:\n"
		if sharders == 0 {
			shardersData = ""
		}
		sharderNodes := []types.Sharder{}

		for i := 1; i <= sharders; i++ {
			var wallet *zcncrypto.Wallet
			wallet, err = getWallet(clientSigScheme)
			if err != nil {
				panic(err)
			}
			sharderNode, sharderData, err := generateSharderNodeStructure(wallet, clientSigScheme, i, shardersConfigData[i-1])
			if err != nil {
				panic(err)
			}
			shardersData += sharderData
			sharderNodes = append(sharderNodes, sharderNode)
			path := fmt.Sprintf("keys/b0snode%d_keys.json", i)
			err = saveWallet(path, wallet)
			if err != nil {
				log.Fatal(err)
			}
		}

		completedData := minersData + shardersData

		nodes := types.Nodes{
			Miners:   minerNodes,
			Sharders: sharderNodes,
		}

		postReq, err := util.NewHTTPPostRequest(server_url+"nodes", nodes)

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

		fmt.Println("Writing the files to nodes.yaml")
		writeToFile(file, completedData)
		fmt.Println(completedData)
	},
}

func getWallet(scheme string) (wallet *zcncrypto.Wallet, err error) {

	sigScheme := zcncrypto.NewSignatureScheme(scheme)

	switch sigScheme.(type) {
	case *zcncrypto.ED255190chainScheme:
		wallet, err = sigScheme.GenerateKeys()
		if err != nil {
			return nil, err
		}
	case *zcncrypto.HerumiScheme:
		wallet, err = sigScheme.GenerateKeys()
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("key-gen", fmt.Sprintf("unrecognized scheme %s", scheme))
	}

	return wallet, err
}

// TODO: refactor miner and sharder structures to a single function later
// TODO: Need to map the return type which was causing some complications
func generateMinerNodeStructure(wallet *zcncrypto.Wallet, scheme string, number int, minerNodeData types.ConfigNodeData) (node types.Miner, details string, err error) {
	if len(wallet.Keys) == 0 {
		return types.Miner{}, "", errors.New("key-gen", "Writing keys failed. Empty wallet.")
	}

	if scheme != "bls0chain" {
		// TODO: Discuss what to write here
		//b := bufio.NewWriter(w)

		//if _, err = b.WriteString(fmt.Sprintf("%s\n", wallet.Keys[0].PublicKey)); err != nil {
		//	return sec, pub, id, err
		//}
		//
		//if _, err = b.WriteString(wallet.Keys[0].PrivateKey + "\n"); err != nil {
		//	return sec, pub, id, err
		//}
		//
		//if err = b.Flush(); err != nil {
		//	return sec, pub, id, err
		//}

		return
	}

	privKey, _ := wallet.Keys[0].PrivateKey, wallet.Keys[0].PublicKey

	var sk bls.SecretKey

	err = sk.DeserializeHexStr(privKey)
	if err != nil {
		return types.Miner{}, "", err
	}
	sec := hex.EncodeToString(sk.GetLittleEndian())
	pub := sk.GetPublicKey().SerializeToHexStr()

	decodeString, _ := hex.DecodeString(pub)
	id := encryption.Hash(decodeString)

	var nodeStructure string

	setIndex := strconv.Itoa(number - 1)

	n2nIp := minerNodeData.N2NIp
	publicIp := minerNodeData.PublicIp
	port := minerNodeData.Port
	path := "miner01"
	description := minerNodeData.Description
	mpk := core.CreateMpk(T, N, number-1, id)

	if mpk == nil {
		log.Fatal("mpk could not be saved")
	}

	nodeStructure = fmt.Sprintf("- id: %s\n  public_key: %s\n  private_key: %s\n  n2n_ip: %s\n  public_ip: %s\n  port: %s\n  path: %s\n  description: %s\n  set_index: %s\n", id, pub, sec, n2nIp, publicIp, port, path, description, setIndex)

	node = types.Miner{
		ID:          id,
		N2NIp:       n2nIp,
		PublicKey:   pub,
		Port:        port,
		PublicIp:    publicIp,
		Path:        path,
		Description: description,
		SetIndex:    uint(number - 1),
		MPK:         mpk,
	}

	filePath := fmt.Sprintf("keys/b0mnode%d_keys.txt", number)

	val := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", pub, privKey, n2nIp, publicIp, port)
	data := []byte(val)

	err = ioutil.WriteFile(filePath, data, 0644)

	if err != nil {
		log.Fatal(err)
	}

	return node, nodeStructure, nil

}

func generateSharderNodeStructure(wallet *zcncrypto.Wallet, scheme string, number int, sharderNodeData types.ConfigNodeData) (node types.Sharder, details string, err error) {
	if len(wallet.Keys) == 0 {
		return types.Sharder{}, "", errors.New("key-gen", "Writing keys failed. Empty wallet.")
	}

	if scheme != "bls0chain" {
		// TODO: Discuss what to write here
		//b := bufio.NewWriter(w)

		//if _, err = b.WriteString(fmt.Sprintf("%s\n", wallet.Keys[0].PublicKey)); err != nil {
		//	return sec, pub, id, err
		//}
		//
		//if _, err = b.WriteString(wallet.Keys[0].PrivateKey + "\n"); err != nil {
		//	return sec, pub, id, err
		//}
		//
		//if err = b.Flush(); err != nil {
		//	return sec, pub, id, err
		//}

		return
	}

	privKey, _ := wallet.Keys[0].PrivateKey, wallet.Keys[0].PublicKey

	var sk bls.SecretKey

	err = sk.DeserializeHexStr(privKey)
	if err != nil {
		return types.Sharder{}, "", err
	}
	sec := hex.EncodeToString(sk.GetLittleEndian())
	pub := sk.GetPublicKey().SerializeToHexStr()

	decodeString, _ := hex.DecodeString(pub)
	id := encryption.Hash(decodeString)

	var nodeStructure string

	setIndex := "0"

	n2nIp := sharderNodeData.N2NIp
	publicIp := sharderNodeData.PublicIp
	port := sharderNodeData.Port
	path := "sharder01"
	description := sharderNodeData.Description

	nodeStructure = fmt.Sprintf("- id: %s\n  public_key: %s\n  private_key: %s\n  n2n_ip: %s\n  public_ip: %s\n  port: %s\n  path: %s\n  description: %s\n  set_index: %s\n", id, pub, sec, n2nIp, publicIp, port, path, description, setIndex)

	node = types.Sharder{
		ID:          id,
		N2NIp:       n2nIp,
		PublicKey:   pub,
		Port:        port,
		PublicIp:    publicIp,
		Path:        path,
		Description: description,
	}

	filePath := fmt.Sprintf("keys/b0snode%d_keys.txt", number)

	val := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", pub, privKey, n2nIp, publicIp, port)
	data := []byte(val)

	err = ioutil.WriteFile(filePath, data, 0644)

	if err != nil {
		log.Fatal(err)
	}

	return node, nodeStructure, nil

}

func writeToFile(w io.Writer, data string) (err error) {
	b := bufio.NewWriter(w)
	if _, err = b.WriteString(data); err != nil {
		return err
	}

	if err = b.Flush(); err != nil {
		return err
	}

	return
}

func init() {
	rootCmd.AddCommand(generateKeys)

	generateKeys.PersistentFlags().String("signature_scheme", "", "Defines the signature scheme used for signing contracts. Either of: ed25519 or bls0chain")
	generateKeys.MarkPersistentFlagRequired("signature_scheme")
	generateKeys.PersistentFlags().Int("miners", 3, "Number of miners for which keys needs to be generated")
	generateKeys.PersistentFlags().Int("sharders", 3, "Number of sharders for which keys needs to be generated")
	generateKeys.PersistentFlags().Bool("overwrite", false, "Overwrite existing nodes.yaml file")
}

func saveWallet(path string, wallet *zcncrypto.Wallet) error {

	data, err := json.Marshal(wallet)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}
