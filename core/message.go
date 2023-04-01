package core

import (
	"encoding/hex"
	"fmt"
	"log"
	"onboarding-cli/types"

	"github.com/0chain/gosdk/core/encryption"
	"github.com/herumi/bls-go-binary/bls"
)

// SignMessages verifies the share and signs the share with the private key for every share sent to it by other miners
func SignMessages(shares map[string]string, mpks map[string][]string, privKey string, currId string) []*types.SignData {

	mp := make([]*types.SignData, 0)

	var privateKey bls.SecretKey

	privateKeyBytes, err := hex.DecodeString(privKey)

	if err != nil {
		log.Panic(err)
	}

	if err := privateKey.SetLittleEndian(privateKeyBytes); err != nil {
		log.Panic(err)
	}

	for id, share := range shares {

		mpk, ok := mpks[id]

		if !ok {

			fmt.Printf("PANIC:mpk not found for id %s\n", id)
			continue
		}
		var jpk []PublicKey

		for _, v := range mpk {
			var pk PublicKey
			pk.SetHexString(v)
			jpk = append(jpk, pk)
		}

		var shareKey bls.SecretKey

		shareKey.SetHexString(share)

		if !ValidateShare(jpk, shareKey, ComputeIDdkg(currId)) {
			fmt.Println("PANIC:share validation failed")
			continue
		}

		message := encryption.Hash(shareKey.GetHexString())

		sign := privateKey.Sign(message).SerializeToHexStr()

		signData := &types.SignData{
			Sign:      sign,
			FromMiner: currId,
			ToMiner:   id,
			Message:   message,
		}

		mp = append(mp, signData)
	}

	return mp
}

// ValidateShare verifies the share sent by other miners using their mpk and the PartyID of the local miner
func ValidateShare(jpk []PublicKey, sij bls.SecretKey, id PartyID) bool {
	var expectedSijPK PublicKey
	if err := expectedSijPK.Set(jpk, &id); err != nil {
		return false
	}
	sijPK := sij.GetPublicKey()
	return expectedSijPK.IsEqual(sijPK)
}
