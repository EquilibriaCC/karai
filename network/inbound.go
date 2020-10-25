package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/harrisonhesslink/flatend"
	"log"
	"strconv"
	//"io/ioutil"
	"github.com/karai/go-karai/transaction"
	"github.com/karai/go-karai/util"
	// "time"
	//"encoding/json"
)

/*
This function handles request for transactions. It takes a top consensus tx hash. 100 txes per batch
*/
func (s *Server) HandleGetTxes(ctx *flatend.Context, request []byte) {
	command := BytesToCmd(request[:commandLength])

	var buff bytes.Buffer
	var payload GetTxes

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Println("ERROR HandleGetTxs: Failed to decode payload", err)
		return
	}
	log.Println(util.Rcv + " [" + command + "] Get Tx from: " + payload.TopHash)
	lastHash := payload.TopHash

	if !s.Protocol.Dat.HaveTx(lastHash) {
		db, err := s.Protocol.Dat.Connect()
		if err != nil {
			log.Println("ERROR HandleGetTxes: creating a DB connection: ", err)
			return
		}
		defer db.Close()

		var transactions []transaction.Transaction
		hit := false
		rows, err := db.Queryx("SELECT * FROM " + s.Protocol.Dat.Cf.GetTableName() + " WHERE tx_type='1' ORDER BY tx_time ASC")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			var thisTx transaction.Transaction
			err = rows.StructScan(&thisTx)
			if err != nil {
				log.Println("ERROR HandleGetTxes (1): Failed to get db data: ", err)
				return
			}

			if lastHash == thisTx.Hash {
				hit = true
			}

			if hit == true {
				transactions = append(transactions, thisTx)

				//loop through to find contracts
				row2, err := db.Queryx("SELECT * FROM "+s.Protocol.Dat.Cf.GetTableName()+" WHERE tx_type='3' AND tx_prnt=$1 ORDER BY tx_time ASC", thisTx.Hash)
				if err != nil {
					panic(err)
				}
				for row2.Next() {
					var thisTx2 transaction.Transaction
					err = row2.StructScan(&thisTx2)
					if err != nil {
						log.Println("ERROR HandleGetTxes (2): Failed to get db data: ", err)
						return
					}
					transactions = append(transactions, thisTx2)
					//loop through to find oracle data
					row3, err := db.Queryx("SELECT * FROM "+s.Protocol.Dat.Cf.GetTableName()+" WHERE tx_type='2' AND tx_epoc=$1 ORDER BY tx_time ASC", thisTx2.Epoc)
					if err != nil {
						panic(err)
					}
					for row3.Next() {
						var thisTx3 transaction.Transaction
						err = row3.StructScan(&thisTx3)
						if err != nil {
							log.Println("ERROR HandleGetTxes (3): Failed to get db data: ", err)
							return
						}
						transactions = append(transactions, thisTx3)
						row3.Close()
					}
					err = row3.Err()
					if err != nil {
						log.Panic(err)
					}
					row2.Close()
				}
				err = rows.Err()
				if err != nil {
					log.Panic(err)
				}
			}
		}

		// get any error encountered during iteration
		err = rows.Err()
		if err != nil {
			log.Println("ERROR HandleGetTxes:", err)
			return
		}
		var txes [][]byte
		for i, tx := range transactions {

			txes = append(txes, tx.Serialize())
			if (i % 100) == 0 {
				data := GobBatchTx{txes, len(transactions)}
				payload := GobEncode(data)
				request := append(CmdToBytes("batchtx"), payload...)

				go s.SendData(ctx, request)
				txes = nil
			}
		}
	}
}

func (s *Server) HandleGetData(ctx *flatend.Context, request []byte) {
	command := BytesToCmd(request[:commandLength])

	var buff bytes.Buffer
	var payload GetData

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)
	if err != nil {
		log.Println("ERROR HandleGetData:", err)
		return
	}

	log.Println(util.Rcv + " [" + command + "] Data Request from: " + ctx.ID.Pub.String())
}

func (s *Server) HandleBatchTx(request []byte) {
	if s.Sync == false {
		command := BytesToCmd(request[:commandLength])

		var buff bytes.Buffer
		var payload GobBatchTx

		buff.Write(request[commandLength:])
		err := gob.NewDecoder(&buff).Decode(&payload)
		if err != nil {
			log.Println("ERROR HandleBatchTx:", err)
			return
		}

		for _, tx_ := range payload.Batch {
			tx := transaction.DeserializeTransaction(tx_)
			if s.Protocol.Dat.HaveTx(tx.Prev) {
				if !s.Protocol.Dat.HaveTx(tx.Hash) {
					s.Protocol.Dat.CommitDBTx(tx)
				}
			}
		}
		percentageString := fmt.Sprintf("%.2f", float64(payload.TotalSent) / float64(s.txNeed) * 100)
		log.Println(util.Rcv + " [" + command + "] Received Transactions. Sync %:" + percentageString + "[" + strconv.Itoa(payload.TotalSent) + "/" + strconv.Itoa(s.txNeed) + "]")
		if payload.TotalSent == s.txNeed {
			s.txNeed = 0
			s.Sync = false
		}
	}
}

func (s *Server) HandleTx(request []byte) {
	command := BytesToCmd(request[:commandLength])

	var buff bytes.Buffer
	var payload GobTx

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)
	if err != nil {
		log.Println("ERROR HandleTx:", err)
		return
	}
	txData := payload.TX
	tx := transaction.DeserializeTransaction(txData)

	log.Println(util.Rcv + " [" + command + "] Transaction: " + tx.Hash)

	if s.Protocol.Dat.HaveTx(tx.Prev) {
		if !s.Protocol.Dat.HaveTx(tx.Hash) {
			s.Protocol.Dat.CommitDBTx(tx)
		}
	}
}

func (s *Server) HandleVersion(ctx *flatend.Context, request []byte) {
	command := BytesToCmd(request[:commandLength])
	var buff bytes.Buffer
	var payload Version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.TxSize > s.Protocol.Dat.GetDAGSize() {
		if s.Sync == false {
			go s.SendGetTxes(ctx)
			s.Sync = true
		}
	}

	log.Println(util.Rcv + " [" + command + "] Node has Num Tx: " + strconv.Itoa(payload.TxSize))
}

func (s *Server) HandleConnection(req []byte, ctx *flatend.Context) {

	command := BytesToCmd(req[:commandLength])
	switch command {
	case "gettxes":
		go s.HandleGetTxes(ctx, req)
	case "getdata":
		go s.HandleGetData(ctx, req)
	case "tx":
		go s.HandleTx(req)
	case "batchtx":
		go s.HandleBatchTx(req)
	case "version":
		go s.HandleVersion(ctx, req)
	}

}
