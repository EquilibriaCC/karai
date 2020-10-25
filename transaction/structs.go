package transaction

type RequestOracleData struct {
	Hash      string `json:"hash"`
	PubKey    string `json:"pubkey"`
	Signature string `json:"signature"`
	Data      string `json:"data"`
	Task      string `json:"task"`
	Height    string `json:"height"`
	Source    string `json:"source"`
	Epoc      string `json:"epoc"`
}

type RequestConsensus struct {
	Hash      string   `json:"hash"`
	PubKey    string   `json:"pubkey"`
	Signature string   `json:"signature"`
	Data      []string `json:"data"`
	Task      string   `json:"task"`
	Height    string   `json:"height"`
}

type Transaction struct {
	Time string `json:"time" db:"tx_time"`
	Type string `json:"type" db:"tx_type"`
	Hash string `json:"hash" db:"tx_hash"`
	Data string `json:"data" db:"tx_data"`
	Prev string `json:"prev" db:"tx_prev"`
	Epoc string `json:"epoc" db:"tx_epoc"`
	Subg string `json:"subg" db:"tx_subg"`
	Prnt string `json:"prnt" db:"tx_prnt"`
	Mile bool   `json:"mile" db:"tx_mile"`
	Lead bool   `json:"lead" db:"tx_lead"`
}

type RequestContract struct {
	Asset string `json:"asset"`
	Denom string `json:"denom"`
}
