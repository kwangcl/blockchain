package main

import (
	"net"
	"os"
	"fmt"
	"time"
	"sync"
	"strconv"
)


var msg_manager *P2PMsgManager

func main() {
	var wg sync.WaitGroup
	var str string

	MsgManager = NewP2PMsgManager()


	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ip_net, ok := addr.(*net.IPNet); ok && !ip_net.IP.IsLoopback() {
			if ip_net.IP.To4() != nil {
				//str = ip_net.IP.String() + "-"
				str = "255.255.255.255"
			}
		}
	}
	str += strconv.FormatInt(time.Now().Unix(),10)
	
	var a[30]byte
	copy(a[:], str)
	fmt.Println(a)

	if len(os.Args) > 1 {
		p2p_client := NewP2PClient()
		p2p_server := NewP2PServer(6667)    

		p2p_client.p2p_server = p2p_server
		p2p_server.p2p_client = p2p_client
        wg.Add(1)
		go func() {           
			p2p_server.StartServer()
			defer wg.Done()
		}()
		tmp_server := p2p_client.ConnectServer("127.0.0.1",6666)
		p2p_client.IPBroadcast(tmp_server)
		wg.Wait()
	} else {
    	p2p_client := NewP2PClient()
	    p2p_server := NewP2PServer(6666)    

		p2p_client.p2p_server = p2p_server
		p2p_server.p2p_client = p2p_client
		fmt.Println("!!!!!!!!!!!!")
		wg.Add(1)
		go func() {
			p2p_server.StartServer()
			defer wg.Done()
		}()

		wg.Wait()
	}
}
