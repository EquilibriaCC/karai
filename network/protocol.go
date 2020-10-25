package network

import (
	// "encoding/hex"
	"log"
	//"github.com/karai/go-karai/database"
	config "github.com/karai/go-karai/configuration"
	"github.com/karai/go-karai/database"
	"github.com/harrisonhesslink/flatend"
	"strconv"
	"github.com/glendc/go-external-ip"
	"github.com/karai/go-karai/transaction"
	"github.com/karai/go-karai/util"
	"io/ioutil"
	"time"
	"github.com/lithdew/kademlia"
	"encoding/json"
	//"github.com/gorilla/websocket"

)

func ProtocolInit(c *config.Config, s *Server) {
	var d database.Database
	var p Protocol
	var peerList PeerList

	s.PeerList = &peerList
	d.Cf = c
	s.cf = c

	p.Dat = &d

	s.Protocol = &p

	d.DB_init()

	go s.RestAPI()

  	consensus := externalip.DefaultConsensus(nil, nil)
    // Get your IP,
    // which is never <nil> when err is <nil>.
    ip, err := consensus.ExternalIP()
    if err != nil {
		log.Panic(ip)
	}
	s.ExternalIP = ip.String()
	s.Node = &flatend.Node{
		PublicAddr: ":" + strconv.Itoa(c.Lport),
		BindAddrs:  []string{":" + strconv.Itoa(c.Lport)},
		SecretKey:  flatend.GenerateSecretKey(),
		Services: map[string]flatend.Handler{
			"karai-xeq": func(ctx *flatend.Context) {
				req, err := ioutil.ReadAll(ctx.Body)
				if err != nil {
					log.Panic(err)
				}
				go s.HandleConnection(req, ctx)
			},
		},
	}

	defer s.Node.Shutdown()

	err = s.Node.Start(s.ExternalIP)

	if s.ExternalIP != "167.172.156.118:4201" {
		go s.Node.Probe("167.172.156.118:4201")
	}

	if err != nil {
		log.Println("Unable to connect")
	}

	go s.LookForNodes()

	select {}
}

func (s *Server) HandleCall(stream *flatend.Stream) {
	req, err := ioutil.ReadAll(stream.Reader)
	if err != nil {
		log.Panic(err)
	}
	go s.HandleConnection(req, nil)
}

func (s *Server) GetProviderFromID(id  *kademlia.ID) *flatend.Provider {
	providers := s.Node.ProvidersFor("karai-xeq")
	for _, provider := range providers {
		if provider.GetID().Pub.String() == id.Pub.String(){
			return provider
		}
	}
	return nil
}

func (s *Server) LookForNodes() {
	for {
		if s.PeerList.Count < 9 {
			newIds := s.Node.Bootstrap()

			//probe new nodes

			for _, peer := range newIds {
				s.Node.Probe(peer.Host.String() + ":" + strconv.Itoa(int(peer.Port)))
			}

			providers := s.Node.ProvidersFor("karai-xeq")
			for _, provider := range providers {
					go s.SendVersion(provider)
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func (s *Server) NewDataTxFromCore(req transaction.RequestOracleData) {
	reqString, _ := json.Marshal(req)

	var txPrev string

	db, connectErr := s.Protocol.Dat.Connect()
	defer db.Close()
	util.Handle("Error creating a DB connection: ", connectErr)

	_ = db.QueryRow("SELECT tx_hash FROM " + s.Protocol.Dat.Cf.GetTableName() + " WHERE tx_type='2' AND tx_epoc=$1 ORDER BY tx_time DESC", req.Epoc).Scan(&txPrev)
	
	newTx := transaction.CreateTransaction("2", txPrev, reqString, []string{}, []string{})

	if !s.Protocol.Dat.HaveTx(newTx.Hash) {
		go s.Protocol.Dat.CommitDBTx(newTx)
		go s.BroadCastTX(newTx)
	}
}

func (s *Server) NewConsensusTXFromCore(req transaction.RequestConsensus) {
	reqString, _ := json.Marshal(req)

	var txPrev string

	db, err := s.Protocol.Dat.Connect()
	if err != nil {
		util.Handle("Error creating a DB connection: ", err)
		return
	}
	defer db.Close()

	_ = db.QueryRow("SELECT tx_hash FROM " + s.Protocol.Dat.Cf.GetTableName() + " WHERE tx_type='1' ORDER BY tx_time DESC").Scan(&txPrev)

	newTx := transaction.CreateTransaction("1", txPrev, reqString, []string{}, []string{})
	if !s.Protocol.Dat.HaveTx(newTx.Hash) {
		go s.Protocol.Dat.CommitDBTx(newTx)
		go s.BroadCastTX(newTx)
	}
}

func (s *Server) CreateContract(asset string, denom string) {
	var txPrev string
	contract := transaction.RequestContract{Asset: asset, Denom: denom}
	jsonContract,_ := json.Marshal(contract)

	db, err := s.Protocol.Dat.Connect()
	if err != nil {
		util.Handle("Error creating a DB connection: ", err)
	}
	defer db.Close()

	_ = db.QueryRow("SELECT tx_hash FROM " + s.Protocol.Dat.Cf.GetTableName() + " WHERE tx_type='1' ORDER BY tx_time DESC").Scan(&txPrev)

	tx := transaction.CreateTransaction("3", txPrev, jsonContract, []string{}, []string{})

	if !s.Protocol.Dat.HaveTx(tx.Hash) {
		go s.Protocol.Dat.CommitDBTx(tx)
		go s.BroadCastTX(tx)
	}
	log.Println("Created Contract " + tx.Hash[:8]+ ": " + asset + "/" + denom)
}

/*

CheckNode checks if a Node should be able to put data on the contract takes a Transaction

*/
func (s *Server) CheckNode(tx transaction.Transaction) bool {

	checksOut := false
	var hash string
	var txData string

	db, connectErr := s.Protocol.Dat.Connect()
	defer db.Close()
	util.Handle("Error creating a DB connection: ", connectErr)

	_ = db.QueryRow("SELECT tx_hash, tx_data FROM " + s.Protocol.Dat.Cf.GetTableName() + " WHERE tx_type='1' && tx_epoc=$1 ORDER BY tx_time DESC", tx.Epoc).Scan(&hash, &txData)

	if hash != "" {
		checksOut = true
	}

	var lastConsensus transaction.RequestConsensus
	err := json.Unmarshal([]byte(txData), &lastConsensus)
	if err != nil {
		//unable to parse last consensus ? this should never happen
		log.Println("Failed to Parse Last Consensus TX on Cehck")
		return false
	}

	//get interface for checks [Request_Consensus, Request_Oracle_Data, Request_Contract]

	result := tx.ParseInterface()
	if result == nil {
		return false
	}

	switch v := result.(type) {
	case transaction.RequestConsensus:
		isFound := false
		for _, key := range lastConsensus.Data {
			if key == v.PubKey {
				isFound = true
				break
			}
		}
		if !isFound {
			return false
		}
		// here v has type T
		break;
	case transaction.RequestOracleData:
		// here v has type S
		break;
	case transaction.RequestContract:
		break;
	default:
		return false;
	}

	return checksOut
}


