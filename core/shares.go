package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"onboarding-cli/types"

	"github.com/herumi/bls-go-binary/bls"
)

func CreateShares(minerIds []string, setIndex int, currId string) map[string]*types.ShareData {

	name := GetSummariesName(setIndex)

	mp := make(map[string]*types.ShareData)

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

		mp[id] = &types.ShareData{
			Share:  share.GetHexString(),
			FromID: currId,
		}

	}

	return mp
}

func ComputeDKGKeyShare(msk []bls.SecretKey, forID PartyID) (Key, error) {
	var secVec Key
	err := secVec.Set(msk, &forID)
	if err != nil {
		return Key{}, err
	}

	return secVec, nil
}
