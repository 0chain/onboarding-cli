This repository contains the code for Magic Block CLI

## Commands

- To create nodes.yaml,operational wallets,mpks and send the miner and sharder information to the server

```
go run main.go generate-keys --signature_scheme bls0chain --miners NO_OF_MINERS --sharders NO_OF_SHARDERS
```

- To create and send shares to the server

```
go run main.go send-shares
```

- To verify the shares and send the signatures to the server

```
go run main.go validate-shares
```

-- To get the magic block

```
go run main.go get-magicblock
```

## Generated Files

- nodes.yaml
- keys/b0mnode${i}\_keys.json (miner)
- keys/b0snode${i}\_keys.json (sharder)
- dkgSummary-${i}\_dkg.json
- b0magicBlock.json
