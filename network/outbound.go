package network

import (
	"bytes"
	"fmt"
	"github.com/harrisonhesslink/flatend"
	"github.com/karai/go-karai/logger"
	"github.com/karai/go-karai/transaction"
	"log"
	"io/ioutil"
	//"math/rand"
	//"time"
	"strconv"
	//"github.com/lithdew/kademlia"
)

func(s *Server)  SendTx(p *flatend.Provider, tx transaction.Transaction) {
	data := GOB_TX{tx.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	stream, err := p.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
	if err == nil {
		logger.Send(" [TXT] Sending Transaction to " + p.GetID().Pub.String() + " ip: " + p.GetID().Host.String())
		go s.HandleCall(stream)

	}
}

func(s *Server)  BroadCastTX(tx transaction.Transaction) {
	data := GOB_TX{tx.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	stream, err := s.node.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
	if err == nil {
		logger.Info(" [TXT] Broadcasting Transaction Out")
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
	providers := s.node.ProvidersFor("karai-xeq")
	log.Println(strconv.Itoa(len(providers)))
	for _, provider := range providers {
		_, err := provider.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(data)))
		if err != nil {
			logger.Error_log(fmt.Sprintf("Unable to broadcast to %s: %s\n", provider.Addr(), err))
		}
	}
}

func (s *Server) SendInv( kind string, items [][]byte) {
	inventory := Inv{nodeAddress, kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("inv"), payload...)

	for _, p := range s.pl.Peers {
			stream, err := p.Provider.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
			if err != nil {
				//fmt.Printf("Unable to broadcast to %s: %s\n", provider.Addr(), err)
			}
			logger.Send(" [INV] Sending INV: " + strconv.Itoa(len(items)))
			s.HandleCall(stream)
	}
}

func (s *Server)SendGetTxes(ctx *flatend.Context, fill bool, contracts map[string]string) {
	
	var txPrev string

	db, connectErr := s.Prtl.Dat.Connect()
	defer db.Close()
	logger.Error_log(fmt.Sprintf(" Error creating a DB connection: %s", connectErr))

	_ = db.QueryRow("SELECT tx_hash FROM " + s.Prtl.Dat.Cf.GetTableName() + " WHERE tx_type='1' ORDER BY tx_time DESC").Scan(&txPrev)
	payload := GobEncode(GetTxes{txPrev, fill, contracts})
	request := append(CmdToBytes("gettxes"), payload...)

	go s.SendData(ctx, request)
	logger.Send(fmt.Sprintf(" [GTXS] Requesting Transactions starting from " + txPrev))
}

func (s *Server) SendVersion(p *flatend.Provider) {
	numTx := s.Prtl.Dat.GetDAGSize()

	payload := GobEncode(Version{version, numTx, s.ExternalIP})

	request := append(CmdToBytes("version"), payload...)

	stream, err := p.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
	if err == nil {
		go s.HandleCall(stream)
		logger.Send(" [VERSION] Call")
	}

}