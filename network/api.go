package network

import (
	"net/http"
	//"strconv"
	"log"
	 "github.com/gorilla/handlers"
	 "github.com/gorilla/mux"
	"github.com/karai/go-karai/transaction"
	"encoding/json"
	// "github.com/gorilla/websocket"
)

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

	// Home
	//api.HandleFunc("/", home).Methods(http.MethodGet)

	// Version
	//api.HandleFunc("/version", returnVersion).Methods(http.MethodGet)

	// Stats
	// api.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
	// 	returnStatsWeb(w, r, keyCollection)
	// }).Methods(http.MethodGet)

	// Transaction by ID
	// api.HandleFunc("/transaction/{hash}", func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	hash := vars["hash"]
	// 	returnSingleTransaction(w, r, hash)
	// }).Methods(http.MethodGet)

	// Transaction by qty
	// api.HandleFunc("/transactions/{number}", func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	number := vars["number"]
	// 	returnNTransactions(w, r, number)
	// }).Methods(http.MethodGet)

	api.HandleFunc("/new_tx", func(w http.ResponseWriter, r *http.Request) {
		var req transaction.Request_Data_TX
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("We are data boy")

		go s.NewDataTxFromCore(req)
	}).Methods("POST")

	api.HandleFunc("/new_consensus_tx", func(w http.ResponseWriter, r *http.Request) {
		var req transaction.Request_Consensus_TX
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		log.Println("We are consensus man")

		go s.NewConsensusTXFromCore(req)
	}).Methods("POST")

	// api.HandleFunc("/tx_api", func(w http.ResponseWriter, r *http.Request) {
	// 	var upgrader = websocket.Upgrader{}
	// 	conn, _ := upgrader.Upgrade(w, r, nil)
	// 	defer conn.Close()
	// 	log.Println("socket open")
	// 	s.HandleAPISocket(conn)
	// })

	// Serve via HTTP
	http.ListenAndServe(":4203", handlers.CORS(headersCORS, originsCORS, methodsCORS)(api))
}
