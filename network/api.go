package network

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/karai/go-karai/transaction"
	"github.com/karai/go-karai/util"
	"net/http"
	"strconv"
	"strings"
	//"strconv"
	"log"
)

func (s *Server) RestAPI() {

	// CORS
	corsAllowedHeaders := []string{
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"Cache-Control",
		"Content-Security-Policy",
		"Feature-Policy",
		"Referrer-Policy",
		"X-Requested-With"}

	corsOrigins := []string{"*", "127.0.0.1"}

	corsMethods := []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}

	headersCORS := handlers.AllowedHeaders(corsAllowedHeaders)
	originsCORS := handlers.AllowedOrigins(corsOrigins)
	methodsCORS := handlers.AllowedMethods(corsMethods)

	// Init API
	api := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(util.Brightyellow + "[API] /")
		w.WriteHeader(http.StatusOK)
		response, err := json.Marshal(map[string]bool{"status": true})
		if err != nil {
			badRequest(w, err)
			return
		}
		_, err = w.Write(response)
		if err != nil {
			badRequest(w, err)
			return
		}
	})

	api.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		log.Println(util.Brightyellow + "[API] /version")
		w.WriteHeader(http.StatusOK)
		response, err := json.Marshal(map[string]interface{}{"status": true, "message": "v1"})
		if err != nil {
			badRequest(w, err)
			return
		}
		_, err = w.Write(response)
		if err != nil {
			badRequest(w, err)
			return
		}
	})

	// Stats
	// api.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
	// 	returnStatsWeb(w, r, keyCollection)
	// }).Methods(http.MethodGet)

	api.HandleFunc("/transactions/{type}/{txs}", func(w http.ResponseWriter, r *http.Request) {
		log.Println(util.Brightyellow + "[API] /transactions/{type}/{txs}")
		var txQuery string
		qry := mux.Vars(r)["txs"]
		numOfTxs, err := strconv.Atoi(qry)
		if err != nil {
			txQuery = qry
			numOfTxs = 1000000000
		}

		_type := mux.Vars(r)["type"]
		if _type != "asc" && _type != "desc" && _type != "contract" {
			badRequest(w, err)
			return
		}

		order := " ORDER BY tx_time " + strings.ToUpper(_type)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		db, err := s.Prtl.Dat.Connect()
		if err != nil {
			badRequest(w, err)
			return
		}
		defer db.Close()

		var queryExtension string
		if txQuery != "" && txQuery != "all" {
			queryExtension = fmt.Sprintf(` WHERE tx_hash = '%s'`, txQuery)
		}
		if txQuery == "nondatatxs" {
			queryExtension = " WHERE tx_type = '1' OR tx_type = '3'"
		}
		if _type == "contract" {
			queryExtension = fmt.Sprintf(" WHERE tx_subg = '%s'", qry)
			order = ""
		}

		var transactions []transaction.Transaction
		rows, _ := db.Queryx("SELECT * FROM " + s.Prtl.Dat.Cf.GetTableName() + queryExtension + order)
		defer rows.Close()
		x := 1
		for rows.Next() {
			var thisTx transaction.Transaction
			err = rows.StructScan(&thisTx)
			if err != nil {
				log.Panic(err)
			}
			transactions = append(transactions, thisTx)
			if x >= numOfTxs {
				break
			}
			x++
		}
		response, err := json.Marshal(transactions)
		if err != nil {
			badRequest(w, err)
			return
		}
		_, err = w.Write(response)
		if err != nil {
			badRequest(w, err)
			return
		}
	}).Methods("GET")

	api.HandleFunc("/new_tx", func(w http.ResponseWriter, r *http.Request) {
		log.Println(util.Brightyellow + "[API] /new_tx")
		if s.sync == false {
			var req transaction.Request_Oracle_Data
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if req.PubKey != "" && req.Signature != "" && req.Hash != "" && req.Task != "" && req.Data != "" && req.Height != "" && req.Source != "" && req.Epoc != "" {
				go s.NewDataTxFromCore(req)
			}
		}
		r.Body.Close()
	}).Methods("POST")

	api.HandleFunc("/get_contracts", func(w http.ResponseWriter, r *http.Request) {
		log.Println(util.Brightyellow + "[API] /new_tx")
		if s.sync == false {

			db, connectErr := s.Prtl.Dat.Connect()
			defer db.Close()
			util.Handle("Error creating a DB connection: ", connectErr)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			transactions := []transaction.Transaction{}

			rows, err := db.Queryx("SELECT * FROM " + s.Prtl.Dat.Cf.GetTableName() + " WHERE tx_type='3' ORDER BY tx_time DESC")
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			for rows.Next() {
				var thisTx transaction.Transaction
				err = rows.StructScan(&thisTx)
				if err != nil {
					// handle this error
					log.Println(err)
				}
				transactions = append(transactions, thisTx)
			}
			// get any error encountered during iteration
			err = rows.Err()
			if err != nil {
				log.Println(err)
				badRequest(w, err)
			}

			response, err := json.Marshal(ArrayTX{transactions})
			if err != nil {
				log.Println(err.Error())
				badRequest(w, err)
			}
			_, err = w.Write(response)
			if err != nil {
				badRequest(w, err)
				return
			}
		} else {
			errorMSG := ErrorJson{"Not Done Syncing", false}
			errorJson, _ := json.Marshal(errorMSG)

			_, err := w.Write(errorJson)
			if err != nil {
				badRequest(w, err)
				return
			}
		}

	}).Methods("GET")

	api.HandleFunc("/new_consensus_tx", func(w http.ResponseWriter, r *http.Request) {
		log.Println(util.Brightyellow + "[API] /new_tx")
		if s.sync == false {
			var req transaction.Request_Consensus
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				badRequest(w, err)
				return
			}
			log.Println("We are consensus man")
			if req.PubKey != "" && req.Signature != "" && req.Hash != "" && req.Task != "" && len(req.Data) > 0 && req.Height != "" {
				go s.NewConsensusTXFromCore(req)
			}
		}
		response, err := json.Marshal(map[string]bool{"status": true})
		if err != nil {
			badRequest(w, err)
			return
		}
		err = r.Body.Close()
		if err != nil {
			badRequest(w, err)
			return
		}
		_, err = w.Write(response)
		if err != nil {
			badRequest(w, err)
			return
		}

	}).Methods("POST")

	// Serve via HTTP
	log.Println("TX API listening on [::]:4203")
	http.ListenAndServe(":4203", handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func badRequest(w http.ResponseWriter, err error) {
	res, _ := json.Marshal(map[string]interface{}{"status": false, "message": err.Error()})
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(res)
	return
}
