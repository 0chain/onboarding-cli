This repository contains the code for Magic Block CLI

## How to setup

1. git clone `git@github.com:0chain/onboarding-cli.git`
2. cd `onboarding-cli`
3. Run `go mod download` to download the dependencies.
4. Create a `config.yaml` file in the directory as shown below

```yaml
miners:
  - n2n_ip: localhost
    public_ip: localhost
    port: 5000
    description: random description
  - n2n_ip: localhost
    public_ip: localhost
    port: 5000
    description: random description
  - n2n_ip: localhost
    public_ip: localhost
    port: 5000
    description: random description
sharders:
  - n2n_ip: localhost
    public_ip: localhost
    port: 6000
    description: random description
  - n2n_ip: localhost
    public_ip: localhost
    port: 6000
    description: random description
  - n2n_ip: localhost
    public_ip: localhost
    port: 6000
    description: random description
```

## Commands

- To create nodes.yaml,operational wallets, mpks and send the miner and sharder information to the server

```
go run main.go generate-keys --signature_scheme bls0chain --miners NO_OF_MINERS --sharders NO_OF_SHARDERS
```

**Note:** NO_OF_MINERS and NO_OF_SHARDERS passed must have same length as in the `config.yaml` file and if u want to overwrite existing nodes and wallet files use `--overwrite` flag.

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
