package network

import (
	"github.com/karai/go-karai/transaction"
	"github.com/karai/go-karai/util"
	"github.com/harrisonhesslink/flatend"
	"bytes"
	"log"
	//"math/rand"
	//"time"
	"strconv"
	"io/ioutil"
	//"github.com/lithdew/kademlia"
)

func(s *Server)  SendTx(p *flatend.Provider, tx transaction.Transaction) {
	data := GobTx{tx.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	stream, err := p.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
	if err == nil {
		log.Println(util.Send + " [TXT] Sending Transaction to " + p.GetID().Pub.String() + " ip: " + p.GetID().Host.String())
		go s.HandleCall(stream)

	}
}

func(s *Server)  BroadCastTX(tx transaction.Transaction) {
	data := GobTx{tx.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	stream, err := s.Node.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
	if err == nil {
		log.Println(util.Send + " [TXT] Broadcasting Transaction Out")
		go s.HandleCall(stream)
	}

}

func (s *Server) SendData(ctx *flatend.Context, data []byte) {

	p := s.GetProviderFromID(&ctx.ID)
	stream, err := p.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(data)))
	if err == nil {
		go s.HandleCall(stream)
	}
}

func (s *Server) BroadCastData(data []byte) {
	providers := s.Node.ProvidersFor("karai-xeq")
	log.Println(strconv.Itoa(len(providers)))
	for _, provider := range providers {
		_, err := provider.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(data)))
		if err != nil {
			util.Handle("BroadCastData", err)
		}
	}
}

func (s *Server) SendInv( kind string, items [][]byte) {
	inventory := Inv{nodeAddress, kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("inv"), payload...)

	for _, p := range s.PeerList.Peers {
			stream, err := p.Provider.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
			if err != nil {
				util.Handle("SendInv", err)
			}
			log.Println(util.Send + " [INV] Sending INV: " + strconv.Itoa(len(items)))
			s.HandleCall(stream)
	}
}

func (s *Server)SendGetTxes(ctx *flatend.Context) {
	
	var txPrev string

	db, connectErr := s.Protocol.Dat.Connect()
	defer db.Close()
	util.Handle("Error creating a DB connection: ", connectErr)

	_ = db.QueryRow("SELECT tx_hash FROM " + s.Protocol.Dat.Cf.GetTableName() + " WHERE tx_type='1' ORDER BY tx_time DESC").Scan(&txPrev)

	payload := GobEncode(GetTxes{txPrev})
	request := append(CmdToBytes("gettxes"), payload...)

	go s.SendData(ctx, request)
	log.Println(util.Send + " [GTXS] Requesting Transactions starting from: " + txPrev)
}

func (s *Server) SendVersion(p *flatend.Provider) {
	numTx := s.Protocol.Dat.GetDAGSize()

	payload := GobEncode(Version{version, numTx, s.ExternalIP})

	request := append(CmdToBytes("version"), payload...)

	stream, err := p.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
	if err == nil {
		go s.HandleCall(stream)
		log.Println(util.Send + " [VERSION] Call")
	}

}