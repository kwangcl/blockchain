package main

import (
	"os"
	"fmt"
	"time"
	"sync"
	"strconv"
	"math/rand"
)


var msg_manager *P2PMsgManager

func main() {
	var wg sync.WaitGroup
	var str string

	MsgManager = NewP2PMsgManager()


	var a[30]byte
	copy(a[:], str)
	fmt.Println(a)

	if len(os.Args) > 1 {
		p2p_client := NewP2PClient()
		p2p_server := NewP2PServer(SERVER_PORT)

		p2p_client.p2p_server = p2p_server
		p2p_server.p2p_client = p2p_client
    wg.Add(1)
		go func() {
			p2p_server.StartServer()
			defer wg.Done()
		}()
		tmp_server := p2p_client.ConnectServer(os.Args[1],SERVER_PORT)
		p2p_client.BroadCastNewNode(tmp_server)


		wg.Wait()
	} else {
    	p2p_client := NewP2PClient()
	    p2p_server := NewP2PServer(SERVER_PORT)

		p2p_client.p2p_server = p2p_server
		p2p_server.p2p_client = p2p_client
		wg.Add(1)
		go func() {
			p2p_server.StartServer()
			defer wg.Done()
		}()
		wg.Add(1)
		go func() {
			for {
				msg := []byte("temperature")
				data := []byte(strconv.Itoa(rand.Intn(40 + 5) - 5))
				tx := CreateTransaction()
				time.Sleep(5 * time.Second)
				defer wg.Done()
			}
		}()

		wg.Wait()
	}
}
