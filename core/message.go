package core

import (
	"encoding/hex"
	"log"

	"github.com/0chain/gosdk/core/encryption"
	"github.com/herumi/bls-go-binary/bls"
)

func SignMessages(shares map[string]string, mpks map[string][]string, privKey string, currId string) map[string]string {

	mp := make(map[string]string)

	for id, share := range shares {

		mpk, ok := mpks[id]

		if !ok {

			log.Panicf("mpk not found for id %s", id)
			continue
		}

		var jpk []PublicKey

		for _, v := range mpk {
			var pk PublicKey
			pk.SetHexString(v)
			jpk = append(jpk, pk)
		}

		var privateKey bls.SecretKey

		privateKeyBytes, err := hex.DecodeString(privKey)

		if err != nil {
			log.Panic(err)
		}

		if err := privateKey.SetLittleEndian(privateKeyBytes); err != nil {
			log.Panic(err)
		}

		var shareKey bls.SecretKey

		shareKey.SetHexString(share)

		if !ValidateShare(jpk, shareKey, ComputeIDdkg(currId)) {
			log.Panicf("invalid share for id %s", id)
			continue
		}

		message := encryption.Hash(shareKey.GetHexString())

		sign := privateKey.Sign(message).GetHexString()

		mp[id] = sign
	}

	return mp
}

func ValidateShare(jpk []PublicKey, sij bls.SecretKey, id PartyID) bool {
	var expectedSijPK PublicKey
	if err := expectedSijPK.Set(jpk, &id); err != nil {
		return false
	}
	sijPK := sij.GetPublicKey()
	return expectedSijPK.IsEqual(sijPK)
}
