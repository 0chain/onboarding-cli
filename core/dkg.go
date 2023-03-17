package core

import (
	"fmt"

	"github.com/herumi/bls-go-binary/bls"
)

type PartyID = bls.ID

type Key = bls.SecretKey

type PublicKey = bls.PublicKey

type DKG struct {
	ID               PartyID
	Msk              []Key
	Mpks             []PublicKey
	MagicBlockNumber int64
	StartingRound    int64
}

func MakeDKG(t, n int, id string) *DKG {
	dkg := &DKG{}
	var secKey Key
	secKey.SetByCSPRNG()

	dkg.ID = ComputeIDdkg(id)
	dkg.Msk = secKey.GetMasterSecretKey(t)
	dkg.Mpks = bls.GetMasterPublicKey(dkg.Msk)
	return dkg
}

func ComputeIDdkg(minerID string) PartyID {
	var forID PartyID
	if err := forID.SetHexString("1" + minerID[:31]); err != nil {
		fmt.Printf("Error while computing ID %s\n", forID.GetHexString())
	}
	return forID
}
