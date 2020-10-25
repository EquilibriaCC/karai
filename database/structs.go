package database

import (
	config "github.com/karai/go-karai/configuration"
	"github.com/karai/go-karai/transaction"
)

type Database struct {
	Cf *config.Config
	thisSubgraph          string
	thisSubgraphShortName string
	poolInterval          int
	txCount               int
}

// Graph is a collection of transactions
type Graph struct {
	Transactions []transaction.Transaction `json:"transactions"`
}
