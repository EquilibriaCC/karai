package network

import (
	"github.com/karai/go-karai/transaction"
	"github.com/harrisonhesslink/flatend"
	"bytes"
	"log"
	//"math/rand"
	//"time"
	"strconv"
	"io/ioutil"
	"github.com/lithdew/kademlia"
)

func(s Server) RequestTxes() {
	//for _, node := range KnownNodes {
		///s.SendGetTxes(node)
	//}
}

// func (s Server) SendAddr(provider flatend.Provider) {
// 	nodes := Addr{KnownNodes}
// 	nodes.AddrList = append(nodes.AddrList, nodeAddress)
// 	payload := GobEncode(nodes)
// 	request := append(CmdToBytes("addr"), payload...)

// 	s.SendData(address, request)
// }

func(s *Server)  SendTx(p *flatend.Provider, tx transaction.Transaction) {
	data := GOB_TX{tx.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	stream, err := p.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
	if err != nil {
		//fmt.Printf("Unable to broadcast to %s: %s\n", provider.Addr(), err)
	}
	log.Println("[SEND] [TXT] Sending Transaction to " + p.GetID().Pub.String() + " ip: " + p.GetID().Host.String())
	s.HandleCall(stream)
}

func (s *Server) SendData(peer *kademlia.ID, data []byte) {

}

func (s *Server) BroadCastData(data []byte) {
	providers := s.node.ProvidersFor("karai-xeq")
	log.Println(strconv.Itoa(len(providers)))
	for _, provider := range providers {
		_, err := provider.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(data)))
		if err != nil {
			//fmt.Printf("Unable to broadcast to %s: %s\n", provider.Addr(), err)
		}
	}
}

// func (s Server) SendBroadcastTX(tx transaction.Transaction) {
// 	data := GOB_TX{tx.Serialize()}
// 	payload := GobEncode(data)
// 	request := append(CmdToBytes("broadtx"), payload...)

// 	rand.Seed(time.Now().UnixNano())
// 	// providers := s.node.ProvidersFor("karai-xeq")
// 	// for _, provider := range providers {
// 	// 	s.SendData(*provider, request)
// 	// 	// if err != nil {
// 	// 	// 	fmt.Printf("Unable to broadcast to %s: %s\n", provider.Addr(), err)
// 	// 	// }
// 	// 	log.Println("[SEND] [BRD] Broadcasting Transaction: " + tx.Hash)
// 	// }
// }

// func (s Server) SendBroadcastNewPeer(provider flatend.Provider) {
// 	data := NewPeer{nodeAddress, addr}
// 	payload := GobEncode(data)
// 	request := append(CmdToBytes("newpeer"), payload...)

// 	rand.Seed(time.Now().UnixNano())

// 	ok := false
// 	//loop to make sure it broadcasts not self
// 	for ok == false {
// 		index := rand.Intn(len(KnownNodes))
// 		if KnownNodes[index] != nodeAddress {
// 			//s.SendData(KnownNodes[index], request)
// 			log.Println("[SEND] [BRD] Broadcasting Connected Peer: " + addr)
// 			ok = true
// 		}
// 	}
// }

func (s *Server) SendInv( kind string, items [][]byte) {
	inventory := Inv{nodeAddress, kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("inv"), payload...)

	for _, p := range s.pl.Peers {
			stream, err := p.Provider.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
			if err != nil {
				//fmt.Printf("Unable to broadcast to %s: %s\n", provider.Addr(), err)
			}
			log.Println("[SEND] [INV] Sending INV: " + strconv.Itoa(len(items)))
			s.HandleCall(stream)
	}
}

// func (s Server)SendGetTxes(provider flatend.Provider) {
// 	payload := GobEncode(GetTxes{nodeAddress, s.Prtl.Dat.GetDAGSize()})
// 	request := append(CmdToBytes("gettxes"), payload...)

// 	//s.SendData(address, request)
// 	//log.Println("[SEND] [GTXS] Requesting Transactions to: " + address)
// }

// func (s Server) SendGetData(provider flatend.Provider, kind string, id []byte) {
// 	payload := GobEncode(GetData{nodeAddress, kind, id})
// 	request := append(CmdToBytes("getdata"), payload...)

// 	//s.SendData(address, request)
// 	//log.Println("[SEND] [GDTA][" + kind + "] Sending Data to: " + address)
// }

func (s *Server) SendVersion(p *flatend.Provider) *flatend.Stream {
	numTx := s.Prtl.Dat.GetDAGSize()

	payload := GobEncode(Version{version, numTx, s.ExternalIP})

	request := append(CmdToBytes("version"), payload...)

	stream, err := p.Push([]string{"karai-xeq"}, nil, ioutil.NopCloser(bytes.NewReader(request)))
	if err != nil {
		//fmt.Printf("Unable to broadcast to %s: %s\n", provider.Addr(), err)
	}
	log.Println("[SEND] [VERSION] Version Call")
	return stream
}