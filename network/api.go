package network

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/karai/go-karai/logger"
	"github.com/karai/go-karai/transaction"
	"github.com/karai/go-karai/util"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	//"strconv"
	"log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// restAPI() This is the main API that is activated when isCoord == true
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

	corsOrigins := []string{
		"*",
		"127.0.0.1"}

	corsMethods := []string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"OPTIONS"}

	headersCORS := handlers.AllowedHeaders(corsAllowedHeaders)
	originsCORS := handlers.AllowedOrigins(corsOrigins)
	methodsCORS := handlers.AllowedMethods(corsMethods)

	// Init API
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	api.Use(logRequestsMiddleware)
	api.Use(s.checkSyncStateMiddleware)

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response, err := json.Marshal(map[string]bool{"status": true})
		if err != nil {
			log.Println(err.Error())
		}
		_, _ = w.Write(response)
	})

	api.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		response, err := json.Marshal(map[string]interface{}{"status": true, "message": "v1"})
		if err != nil {
			badRequest(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			badRequest(w, err)
			return
		}
	})

	api.HandleFunc("/apihits", func(w http.ResponseWriter, r *http.Request) {
		response, err := json.Marshal(count)
		if err != nil {
			badRequest(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			badRequest(w, err)
			return
		}
	})

	api.HandleFunc("/transactions/{type}/{txs}", func(w http.ResponseWriter, r *http.Request) {
		var txQuery, queryExtension string
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

		db, err := s.Prtl.Dat.Connect()
		if err != nil {
			badRequest(w, err)
			return
		}
		defer db.Close()

		if txQuery != "" && txQuery != "all" {
			queryExtension = fmt.Sprintf(` WHERE tx_hash = '%s'`, txQuery)
		}
		if txQuery == "nondatatxs" {
			queryExtension = " WHERE tx_type = '1' OR tx_type = '3'"
		}
		if _type == "contract" {
			queryExtension = fmt.Sprintf(" WHERE tx_subg = '%s'", qry)
			order = ""
			if txQuery == "contractsonly" {
				queryExtension = " WHERE tx_type='3' ORDER BY tx_time DESC"
			}
		}

		var transactions []transaction.Transaction
		rows, err := db.Queryx("SELECT * FROM " + s.Prtl.Dat.Cf.GetTableName() + queryExtension + order)
		if err != nil {
			badRequest(w, err)
			return
		}
		defer rows.Close()
		for x := 0; rows.Next() && x < numOfTxs; x++ {
			var thisTx transaction.Transaction
			err = rows.StructScan(&thisTx)
			if err != nil {
				log.Panic(err)
			}
			transactions = append(transactions, thisTx)
		}
		response, err := json.Marshal(transactions)
		if err != nil {
			badRequest(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			badRequest(w, err)
			return
		}
	}).Methods("GET")

	api.HandleFunc("/new_tx", func(w http.ResponseWriter, r *http.Request) {

		var req transaction.Request_Oracle_Data
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			badRequest(w, err)
			return
		}
		for i, values := 0, reflect.ValueOf(req); i < values.NumField(); i++ {
			if values.Type().Field(i).Type.String() == "string" {
				response, err := json.Marshal(map[string]interface{}{"status": false, "message": "Invalid value " + strconv.Itoa(i)})
				if err != nil {
					badRequest(w, err)
					return
				}
				_, err = w.Write(response)
				if err != nil {
					badRequest(w, err)
				}
				return
			}
		}
		log.Println("TEST")
		go s.NewDataTxFromCore(req)
	}).Methods("POST")

	api.HandleFunc("/get_contracts", func(w http.ResponseWriter, r *http.Request) {
		db, connectErr := s.Prtl.Dat.Connect()
		defer db.Close()
		util.Handle("Error creating a DB connection: ", connectErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// reportRequest("transactions/"+hash, w, r)
		var transactions []transaction.Transaction

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
				log.Panic(err)
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

		json, _ := json.Marshal(ArrayTX{transactions})
		w.Write(json)

	}).Methods("GET")

	api.HandleFunc("/new_consensus_tx", func(w http.ResponseWriter, r *http.Request) {
		var req transaction.Request_Consensus
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			badRequest(w, err)
			return
		}
		if req.PubKey != "" && req.Signature != "" && req.Hash != "" && req.Task != "" && len(req.Data) > 0 && req.Height != "" {
			go s.NewConsensusTXFromCore(req)
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
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			badRequest(w, err)
			return
		}
		r.Body.Close()
	}).Methods("POST")

	// Serve via HTTP
	logger.Info(" [API] TX API listening on [::]:4203")
	http.ListenAndServe(":4203", handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}

func badRequest(w http.ResponseWriter, err error) {
	res, _ := json.Marshal(map[string]interface{}{"status": false, "message": err.Error()})
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(res)
	return
}
