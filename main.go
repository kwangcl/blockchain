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

		wg.Add(1)
		go func() {
			for {
				msg := []byte("weight")
				data := []byte(strconv.Itoa(rand.Intn(90 + 50) - 50))
				tx := CreateTransaction(data, msg)
				src_buf := MsgManager.GenSrcData()
				p2p_server.BroadCastMsg(MsgManager.SendTransactionMsg(tx.Serialize()), string(src_buf))
				p2p_client.BroadCastMsg(MsgManager.SendTransactionMsg(tx.Serialize()), string(src_buf))

				time.Sleep(15* time.Second)
				defer wg.Done()
			}
		}()

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
				tx := CreateTransaction(data, msg)
				src_buf := MsgManager.GenSrcData()
				p2p_server.BroadCastMsg(MsgManager.SendTransactionMsg(tx.Serialize()), string(src_buf))
				p2p_client.BroadCastMsg(MsgManager.SendTransactionMsg(tx.Serialize()), string(src_buf))

				time.Sleep(10 * time.Second)
				defer wg.Done()
			}
		}()

		wg.Wait()
	}
}
