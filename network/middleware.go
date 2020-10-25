package network

import (
	"errors"
	"github.com/karai/go-karai/util"
	"log"
	"net/http"
	"strconv"
)

var count = map[string]int{"total":0}
func logRequestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := count[r.RequestURI]; ok {
			count[r.RequestURI]++
		} else {
			count[r.RequestURI] = 1
		}
		count["total"]++
		log.Println(
			util.Brightyellow, "[API]",
			util.Brightwhite + r.RequestURI,
			"("+strconv.Itoa(count[r.RequestURI])+"/"+strconv.Itoa(count["total"])+")",
			util.Brightred,
		)
		next.ServeHTTP(w, r)
		return
	})
}

func (s *Server) checkSyncStateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.Sync {
			badRequest(w, errors.New("node not synced"))
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}
