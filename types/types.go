package types

type Miner struct {
	ID          string `json:"id" gorm:"primary_key"`
	N2NIp       string `json:"n2n_ip" gorm:"column:n2n_ip"`
	PublicKey   string `json:"public_key"`
	Port        string `json:"port"`
	PublicIp    string `json:"public_ip"`
	Path        string `json:"path"`
	Description string `json:"description"`
	SetIndex    uint   `json:"set_index"`
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

type MPK struct {
	ID  string   `json:"id"`
	Mpk []string `json:"mpk"`
}

type Mpks struct {
	Mpks map[string]*MPK `json:"mpks"`
}
