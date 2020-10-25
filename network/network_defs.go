package network

import (
	"github.com/gorilla/websocket"
	"github.com/harrisonhesslink/flatend"
	config "github.com/karai/go-karai/configuration"
	"github.com/karai/go-karai/database"
	"github.com/karai/go-karai/transaction"
	"github.com/lithdew/kademlia"
)

const (
	version       = 1
	commandLength = 12
)

var (
	nodeAddress   string
	KnownNodes    = []string{"127.0.0.1:3001"}
)

type Addr struct {
	AddrList []string
}

type GobTx struct {
	TX []byte
}

type GobBatchTx struct {
	Batch     [][]byte
	TotalSent int
}

type GetTxes struct {
	TopHash string
}

type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type Version struct {
	Version  int
	TxSize   int
	AddrFrom string
}

type NewPeer struct {
	AddrFrom string
	NewPeer  string
}

type Server struct {
	Protocol     *Protocol
	cf           *config.Config
	Node         *flatend.Node
	PeerList     *PeerList
	ExternalIP   string
	ExternalPort int
	Sockets      []*websocket.Conn
	Sync         bool
	txNeed       int
}

type Protocol struct {
	Dat *database.Database
}

type Peer struct {
	ID       *kademlia.ID
	Provider *flatend.Provider
}

type PeerList struct {
	Peers []Peer
	Count int
}

type ArrayTX struct {
	Txes []transaction.Transaction `json:"txes"`
}

type ErrorJson struct {
	Message string `json:"message"`
	Error   bool   `json:"is_error"`
}
