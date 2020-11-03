package network

import (
	"github.com/karai/go-karai/logger"
	// "encoding/hex"
	"log"
	"encoding/json"
	"github.com/glendc/go-external-ip"
	"github.com/harrisonhesslink/flatend"
	//"github.com/karai/go-karai/database"
	config "github.com/karai/go-karai/configuration"
	"github.com/karai/go-karai/database"
	"github.com/karai/go-karai/transaction"
	"github.com/lithdew/kademlia"
	"io/ioutil"
	"strconv"
	"time"
	//"github.com/gorilla/websocket"

)

func Protocol_Init(c *config.Config, s *Server) {
	var d database.Database
	var p Protocol
	var peer_list PeerList

	s.pl = &peer_list
	d.Cf = c
	s.cf = c

	p.Dat = &d

	s.Prtl = &p

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
	s.node = &flatend.Node{
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

	defer s.node.Shutdown()

	err = s.node.Start(s.ExternalIP)

	if s.ExternalIP != "167.172.156.118:4201" {
		go s.node.Probe("167.172.156.118:4201")
	}

	if err != nil {
		logger.Error_log(" [PROTOCOLINIT] Unable to connect to node")
	}

	go s.LookForNodes()

	select {}
}

func (s *Server) HandleCall(stream *flatend.Stream) {
	req, err := ioutil.ReadAll(stream.Reader)
	if err != nil {
		logger.Error_log(" [HANDLECALL] " + err.Error())
		log.Panic(err)
	}
	go s.HandleConnection(req, nil)
}

func (s *Server) GetProviderFromID(id  *kademlia.ID) *flatend.Provider {
	providers := s.node.ProvidersFor("karai-xeq")
	for _, provider := range providers {
		if provider.GetID().Pub.String() == id.Pub.String(){
			return provider
		}
	}
	return nil
}

func (s *Server) LookForNodes() {
	for {
		if s.pl.Count < 9 {
			new_ids := s.node.Bootstrap()

			//probe new nodes

			for _, peer := range new_ids {
				s.node.Probe(peer.Host.String() + ":" + strconv.Itoa(int(peer.Port)))
			}

			providers := s.node.ProvidersFor("karai-xeq")
			logger.Info(" " + strconv.Itoa(len(providers)))
			for _, provider := range providers {
					go s.SendVersion(provider)
			}
		}

		time.Sleep(30 * time.Second)
	}
}

func (s *Server) NewDataTxFromCore(req transaction.Request_Oracle_Data) {
	req_string, _ := json.Marshal(req)

	var txPrev string

	db, connectErr := s.Prtl.Dat.Connect()
	defer db.Close()
	logger.Error_log(" [NEWDATATXFROMCORE] Error creating a DB connection " + connectErr.Error())

	_ = db.QueryRow("SELECT tx_hash FROM " + s.Prtl.Dat.Cf.GetTableName() + " WHERE tx_type='2' AND tx_epoc=$1 ORDER BY tx_time DESC", req.Epoc).Scan(&txPrev)
	
	new_tx := transaction.CreateTransaction("2", txPrev, req_string, []string{}, []string{})

	if !s.Prtl.Dat.HaveTx(new_tx.Hash) {
		go s.Prtl.Dat.CommitDBTx(new_tx)
		go s.BroadCastTX(new_tx)
	}
}

func (s *Server) NewConsensusTXFromCore(req transaction.Request_Consensus) {
	req_string, _ := json.Marshal(req)

	var txPrev string

	db, connectErr := s.Prtl.Dat.Connect()
	defer db.Close()
	logger.Error_log(" [NEWCONSENSUSTXFROMCORE] Error creating a DB connection " + connectErr.Error())

	_ = db.QueryRow("SELECT tx_hash FROM " + s.Prtl.Dat.Cf.GetTableName() + " WHERE tx_type='1' ORDER BY tx_time DESC").Scan(&txPrev)

	new_tx := transaction.CreateTransaction("1", txPrev, req_string, []string{}, []string{})
	if !s.Prtl.Dat.HaveTx(new_tx.Hash) {
		go s.Prtl.Dat.CommitDBTx(new_tx)
		go s.BroadCastTX(new_tx)
	}
}

func (s *Server) CreateContract(asset string, denom string) {
	var txPrev string
	contract := transaction.Request_Contract{asset, denom}
	json_contract,_ := json.Marshal(contract)

	db, connectErr := s.Prtl.Dat.Connect()
	defer db.Close()
	logger.Error_log(" [CREATECONTRACT] Error creating a DB connection " + connectErr.Error())

	_ = db.QueryRow("SELECT tx_hash FROM " + s.Prtl.Dat.Cf.GetTableName() + " WHERE tx_type='1' ORDER BY tx_time DESC").Scan(&txPrev)

	tx := transaction.CreateTransaction("3", txPrev, []byte(json_contract), []string{}, []string{})

	if !s.Prtl.Dat.HaveTx(tx.Hash) {
		go s.Prtl.Dat.CommitDBTx(tx) 
		go s.BroadCastTX(tx)
	}
	logger.Info(" Created Contract " + tx.Hash[:8]+ ": " + asset + "/" + denom)
}

/*

CheckNode checks if a node should be able to put data on the contract takes a Transaction

*/
func (s *Server) CheckNode(tx transaction.Transaction) bool {

	checks_out := false
	var hash string
	var tx_data string

	db, connectErr := s.Prtl.Dat.Connect()
	defer db.Close()
	logger.Error_log(" [CHECKNODE] Error creating a DB connection " + connectErr.Error())

	_ = db.QueryRow("SELECT tx_hash, tx_data FROM " + s.Prtl.Dat.Cf.GetTableName() + " WHERE tx_type='1' && tx_epoc=$1 ORDER BY tx_time DESC", tx.Epoc).Scan(&hash, &tx_data)

	if hash != "" {
		checks_out = true
	}

	var last_consensus transaction.Request_Consensus
	err := json.Unmarshal([]byte(tx_data), &last_consensus)
	if err != nil {
		//unable to parse last consensus ? this should never happen
		logger.Warning_log(" Failed to parse last consensus TX on check")
		return false
	}

	//get interface for checks [Request_Consensus, Request_Oracle_Data, Request_Contract]

	result := tx.ParseInterface()
	if result == nil {
		return false
	}

	switch v := result.(type) {
	case transaction.Request_Consensus:
		isFound := false
		for _, key := range last_consensus.Data {
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
	case transaction.Request_Oracle_Data:
		// here v has type S
		break;
	case transaction.Request_Contract:
		break;
	default:
		return false;
	}

	return checks_out
}


