package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func CreateMpk(t, n, setIndex int, minerID string) []string {

	var mpk []string

	dkg := MakeDKG(t, n, minerID)

	for _, v := range dkg.Mpks {
		mpk = append(mpk, v.GetHexString())
	}

	var msk []string

	for _, v := range dkg.Msk {
		msk = append(msk, v.GetHexString())
	}

	mp := map[string]any{
		"ID":  dkg.ID.GetHexString(),
		"MSK": msk,
		"MPK": mpk,
	}

	err := saveDKG(mp, setIndex)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("no err")

	return mpk
}

func saveDKG(dkg map[string]any, index int) error {
	var err error
	var dkgData []byte
	if dkgData, err = json.MarshalIndent(dkg, "", " "); err != nil {
		return err
	}

	name := GetSummariesName(index)

	path := fmt.Sprintf("output/%s", name)

	if err := ioutil.WriteFile(path, dkgData, 0644); err != nil {
		return err
	}

	return nil
}

func GetSummariesName(index int) string {
	return fmt.Sprintf("b0mnode%v_dkg.json", index+1)
}
