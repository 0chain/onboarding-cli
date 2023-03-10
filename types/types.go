package types

type Miner struct {
	ID          string
	PublicKey   string
	PrivateKey  string
	N2NIp       string
	Host        string
	Port        string
	PublicIp    string
	Path        string
	Description string
	SetIndex    uint
}

type Sharder struct {
	ID          string
	N2NIp       string
	PublicKey   string
	Host        string
	Port        string
	PublicIp    string
	Path        string
	Description string
}

type Nodes struct {
	Miners   []Miner
	Sharders []Sharder
}
