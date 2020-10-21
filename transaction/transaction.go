package transaction

import (
	"bytes"
	"encoding/gob"
	"log"
	"github.com/karai/go-karai/util"
	"golang.org/x/crypto/sha3"
	"encoding/hex"
	"encoding/json"
)

type Request_Data_TX struct {
	Hash string `json:hash`
	PubKey string `json:pub_key`
	Signature string `json:signature`
	Data string`json:data`
	Task string`json:task`
	Height string `json:height`
	Source string `json:source`
	Epoc string `json:epoc`
}

type Request_Consensus_TX struct {
	Hash string `json:hash`
	PubKey string `json:pub_key`
	Signature string `json:signature`
	Data []string`json:data`
	Task string`json:task`
	Height string `json:height`
}

// Transaction This is the structure of the transaction
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

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}

func CreateTransaction(txType, last_epoc_tx string, data []byte, txhash_on_epoc []string, txdata_on_epoc []string) Transaction {
	var newTx Transaction

	newTx.Type = txType
	// if isCoordinator && txType == "2" {
	if newTx.Type == "2" {
		parsePayload := json.Valid(data)
		if !parsePayload {
			newTx.Data = hex.EncodeToString(data)
		} else if parsePayload {
			newTx.Data = string(data)
		}

		rct := Request_Data_TX{}
		_ = json.Unmarshal(data, &rct)
	
		if (last_epoc_tx == "") {
			newTx.Prev = rct.Epoc

		} else {
			newTx.Prev = last_epoc_tx
		}
			newTx.Time = util.UnixTimeStampNano()
			newTx.Epoc = rct.Epoc
			newTx.Mile = false


			newTx.Prnt = newTx.Epoc
				
			newTx.Hash = hashTransaction(newTx.Time, newTx.Type, newTx.Data, newTx.Prev)
			newTx.Subg = newTx.Epoc		

		return newTx
	} else if newTx.Type == "1" {
		
		parsePayload := json.Valid(data)
		if !parsePayload {
			newTx.Data = hex.EncodeToString(data)
		} else if parsePayload {
			newTx.Data = string(data)
		}

		newTx.Prev = last_epoc_tx

		newTx.Time = util.UnixTimeStampNano()
		newTx.Hash = hashTransaction(newTx.Time, newTx.Type, newTx.Data, last_epoc_tx)
		newTx.Subg = newTx.Hash
		newTx.Epoc = "0"
		newTx.Mile = true
		newTx.Lead = true
		log.Println("[SELF] New Transaction: " + newTx.Hash)
		return newTx
	} else if newTx.Type == "3" {
		parsePayload := json.Valid(data)
		if !parsePayload {
			newTx.Data = hex.EncodeToString(data)
		} else if parsePayload {
			newTx.Data = string(data)
		}

		newTx.Prev = last_epoc_tx
		newTx.Time = util.UnixTimeStampNano()
		newTx.Hash = hashTransaction(newTx.Time, newTx.Type, newTx.Data, last_epoc_tx)
		newTx.Subg = newTx.Hash
		newTx.Epoc = newTx.Hash
		newTx.Mile = true
		newTx.Lead = false
		newTx.Prnt = last_epoc_tx
		return newTx
	}
	return newTx
}

// hashTransaction takes elements of a transaction and computes a hash using SHA512
func hashTransaction(txTime, txType, txData, txPrev string) string {
	hashedData := []byte(txTime + txType + txData + txPrev)
	slot := make([]byte, 64)
	sha3.ShakeSum256(slot, hashedData)
	// fmt.Printf("%x\n", slot)
	txHash := hex.EncodeToString(slot[:])
	// legacy sha512
	// hash := sha512.Sum512(hashedData)
	// txHash := hex.EncodeToString(hash[:])

	return txHash
}
