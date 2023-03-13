package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"onboarding-cli/types"

	"github.com/herumi/bls-go-binary/bls"
)

// CreateShares take the miner ids and the current miner setIndex and  id and returns the shares for all other miners .

func CreateShares(minerIds []string, setIndex int, currId string) []*types.ShareData {

	name := GetSummariesName(setIndex)

	shares := make([]*types.ShareData, 0)

	path := fmt.Sprintf("output/%s", name)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	dkg := make(map[string]any)

	err = json.Unmarshal(data, &dkg)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// partyId := dkg["ID"].(string)

	mskData := dkg["MSK"].([]interface{})

	var msk []bls.SecretKey

	for _, v := range mskData {
		var sk bls.SecretKey
		sk.SetHexString(v.(string))
		msk = append(msk, sk)
	}

	for _, id := range minerIds {
		otherPartyId := ComputeIDdkg(id)
		share, err := ComputeDKGKeyShare(msk, otherPartyId)
		if err != nil {
			log.Panic(err)
		}

		shareData := &types.ShareData{
			Share:     share.GetHexString(),
			FromMiner: currId,
			ToMiner:   id,
		}
		shares = append(shares, shareData)
	}

	return shares
}

// ComputeDKGKeyShare computes the share for a given miner PartyID using its own msk.
func ComputeDKGKeyShare(msk []bls.SecretKey, forID PartyID) (Key, error) {
	var secVec Key
	err := secVec.Set(msk, &forID)
	if err != nil {
		return Key{}, err
	}

	return secVec, nil
}
