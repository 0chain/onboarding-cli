package cmd

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"onboarding-cli/core"
	"onboarding-cli/types"
	"os"
	"strconv"

	"github.com/0chain/errors"
	"github.com/0chain/gosdk/core/encryption"
	"github.com/0chain/gosdk/core/zcncrypto"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/spf13/cobra"
)

var T = 3
var N = 66

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

		var clientSigScheme string

		clientSigScheme, err = flags.GetString("signature_scheme")
		if err != nil {
			log.Fatal(err)
		}
		if !(clientSigScheme == "bls0chain" || clientSigScheme == "ed25519") {
			log.Fatal("Signature Scheme can only take either bls0chain or ed25519")
		}

		var miners, sharders int

		miners, err = flags.GetInt("miners")
		if err != nil {
			log.Fatal(err)
		}

		sharders, err = flags.GetInt("sharders")
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.OpenFile("nodes.yml", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}

		minersData := "miners:\n"

		minerNodes := []types.Miner{}

		for i := 1; i <= miners; i++ {
			var wallet *zcncrypto.Wallet
			wallet, err = getWallet(clientSigScheme)
			if err != nil {
				panic(err)
			}
			minerNode, minerData, err := generateMinerNodeStructure(wallet, clientSigScheme, i)
			if err != nil {
				panic(err)
			}
			minersData += minerData
			minerNodes = append(minerNodes, minerNode)
		}

		shardersData := "sharders:\n"
		sharderNodes := []types.Sharder{}

		for i := 1; i <= sharders; i++ {
			var wallet *zcncrypto.Wallet
			wallet, err = getWallet(clientSigScheme)
			if err != nil {
				panic(err)
			}
			sharderNode, sharderData, err := generateSharderNodeStructure(wallet, clientSigScheme, i)
			if err != nil {
				panic(err)
			}
			shardersData += sharderData
			sharderNodes = append(sharderNodes, sharderNode)
		}

		endData := fmt.Sprintf("\nmessage: %s\nmagic_block_number: 1\nstarting_round: 0\nt_percent: 66\nk_percent: 75", "From CLI")

		completedData := minersData + shardersData + endData

		// nodes := types.Nodes{
		// 	Miners:   minerNodes,
		// 	Sharders: sharderNodes,
		// }

		// postReq, err := util.NewHTTPPostRequest("http://localhost:3000/nodes", nodes)

		// if err != nil {
		// 	fmt.Println(err)
		// }

		// if err != nil {
		// 	panic(err)
		// }
		// postResponse, err := postReq.Post()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println("Post Request Response", postResponse)

		var saveFlag bool
		saveFlag, err = flags.GetBool("save")
		if err != nil {
			log.Fatal(err)
		}
		if saveFlag {
			fmt.Println("Writing the files to nodes.yml")
			writeToFile(file, completedData)
		} else {
			fmt.Println(completedData)
		}
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
func generateMinerNodeStructure(wallet *zcncrypto.Wallet, scheme string, number int) (node types.Miner, details string, err error) {
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

	convertedIndex := strconv.Itoa(number)
	setIndex := convertedIndex

	if number < 10 {
		convertedIndex = "0" + convertedIndex
	}

	n2nIp := "localhost"
	publicIp := "localhost"
	port := "701" + setIndex
	path := "miner" + convertedIndex
	description := ""
	_ = core.CreateMpk(T, N, number, id)

	nodeStructure = fmt.Sprintf("- id: %s\n  public_key: %s\n  private_key: %s\n  n2n_ip: %s\n  public_ip: %s\n  port: %s\n  path: %s\n  description: %s\n  set_index: %s\n", id, pub, sec, n2nIp, publicIp, port, path, description, setIndex)

	node = types.Miner{
		ID:          id,
		N2NIp:       n2nIp,
		PublicKey:   pub,
		Port:        port,
		PublicIp:    publicIp,
		Path:        path,
		Description: description,
		SetIndex:    uint(number),
	}
	return node, nodeStructure, nil

}

func generateSharderNodeStructure(wallet *zcncrypto.Wallet, scheme string, number int) (node types.Sharder, details string, err error) {
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

	convertedIndex := strconv.Itoa(number)
	setIndex := convertedIndex

	if number < 10 {
		convertedIndex = "0" + convertedIndex
	}

	n2nIp := "localhost"
	publicIp := "localhost"
	port := "702" + setIndex
	path := "sharder" + convertedIndex
	description := ""

	nodeStructure = fmt.Sprintf("- id: %s\n  public_key: %s\n  private_key: %s\n  n2n_ip: %s\n  public_ip: %s\n  port: %s\n  path: %s\n  description: %s\n", id, pub, sec, n2nIp, publicIp, port, path, description)

	node = types.Sharder{
		ID:          id,
		N2NIp:       n2nIp,
		PublicKey:   pub,
		Port:        port,
		PublicIp:    publicIp,
		Path:        path,
		Description: description,
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
	generateKeys.PersistentFlags().Bool("save", false, "Save the generated key data in a file instead of printing")
}
