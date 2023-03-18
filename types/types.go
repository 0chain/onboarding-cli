package types

type Miner struct {
	ID          string   `json:"id" gorm:"primary_key" yaml:"id"`
	N2NIp       string   `json:"n2n_ip" gorm:"column:n2n_ip" yaml:"n2n_ip"`
	PublicKey   string   `json:"public_key" yaml:"public_key"`
	Port        string   `json:"port" yaml:"port"`
	PublicIp    string   `json:"public_ip" yaml:"public_ip"`
	Path        string   `json:"path" yaml:"path"`
	Description string   `json:"description" yaml:"description"`
	SetIndex    uint     `json:"set_index" yaml:"set_index"`
	MPK         []string `json:"mpk" yaml:"mpk"`
}

type Sharder struct {
	ID          string `json:"id" gorm:"primary_key"`
	N2NIp       string `json:"n2n_ip" gorm:"column:n2n_ip"`
	PublicKey   string `json:"public_key"`
	Port        string `json:"port"`
	PublicIp    string `json:"public_ip"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type Nodes struct {
	Miners   []Miner   `json:"miners"`
	Sharders []Sharder `json:"sharders"`
}

type ShareData struct {
	Share     string `json:"share"`
	FromMiner string `json:"from_miner"`
	ToMiner   string `json:"to_miner"`
}

type SignData struct {
	Sign      string `json:"sign"`
	FromMiner string `json:"from_miner"`
	ToMiner   string `json:"to_miner"`
}

type ShareServer struct {
	Shares []*ShareData `json:"shares"`
}
