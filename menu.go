package main

import (
	"bufio"
	"github.com/karai/go-karai/network"
	"github.com/karai/go-karai/util"
	"log"
	"os"
	"strconv"
	"strings"
)

// returns to listening to input.
func inputHandler(s *network.Server/*keyCollection *ED25519Keys*/) {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if strings.Compare("help", text) == 0 {
		//	menu()
		} else if strings.Compare("?", text) == 0 {
			//menu()
		} else if strings.Compare("peer", text) == 0 {
		//	fmt.Printf(brightcyan + "Peer ID: ")
		//	fmt.Printf(cyan+"%s\n", keyCollection.publicKey)
		} else if strings.Compare("cleardb", text) == 0 {
			if s.Prtl.Dat.TruncateTable() {
				log.Println(util.Brightwhite + "Cleared Database")
				continue
			}
			log.Println(util.Brightred + "Failed to clear DB")
		} else if strings.Compare("version", text) == 0 {
			//menuVersion()
		} else if strings.Compare("license", text) == 0 {
		//	printLicense()
		} else if strings.Compare("dag", text) == 0 {
			count := s.Prtl.Dat.GetDAGSize()
			log.Println(util.Brightwhite + "Txes: " + strconv.Itoa(count))
		} else if strings.Compare("a", text) == 0 {
			// // start := time.Now()
			// // txint := 50
			// // addBulkTransactions(txint)
			// // elapsed := time.Since(start)
			// fmt.Printf("\nWriting %v objects to memory took %s seconds.\n", txint, elapsed)
		} else if strings.HasPrefix(text, "ban ") {
			// bannedPeer := strings.TrimPrefix(text, "ban ")
			// banPeer(bannedPeer)
		} else if strings.HasPrefix(text, "unban ") {
			// unBannedPeer := strings.TrimPrefix(text, "unban ")
			// unBanPeer(unBannedPeer)
		} else if strings.HasPrefix(text, "blacklist") {
			// blackList()
		} else if strings.Compare("clear blacklist", text) == 0 {
			// clearBlackList()
		} else if strings.Compare("clear peerlist", text) == 0 {
			// clearPeerList()
		} else if strings.Compare("peerlist", text) == 0 {
			// whiteList()
		} else if strings.Compare("exit", text) == 0 {
			// menuExit()
		} else if strings.Compare("generate-pointer", text) == 0 {
			// generatePointer()
		} else if strings.Compare("quit", text) == 0 {
			// menuExit()
		} else if strings.Compare("close", text) == 0 {
			// menuExit()
		} else if strings.Compare("nodes", text) == 0 {
			// nodes := "[ "
			// for _, node := range KnownNodes {
			// 	nodes += node + " "
			// }
			// nodes += "]"
			// log.Println(nodes)
		} else if strings.HasPrefix(text, "create_contract ") {
			strings.TrimPrefix(text, "create_contract ")
			args := strings.Fields(text)
			if args[1] == "XHV" || args[1] == "XEQ" || args[1] == "LOKI" || args[1] == "ETH" || args[1] == "DOGE" {
				if args[2] == "BTC" {
					go s.CreateContract(args[1], args[2])
					} else {
						log.Println("Pair Not Supported! BTC")
					}

			} else {
				log.Println("Pair Not Supported! XEQ, XHV, LOKI, ETH, DOGE")

			}
		}
	}
}