package main

import (
	"net"
	"strconv"
	"fmt"
	"sync"
)

var client_map_lock = sync.RWMutex{}
var CLIENT_MAX_CONNECTION = 2


type P2PClient struct {
	servers map[string]bool
	node	*P2PNode
	p2p_server *P2PServer
}

func NewP2PClient() *P2PClient {
	
	client_node := NewNode(nil)
	return &P2PClient{map[string]bool{}, client_node, nil}
}

func (client *P2PClient)ConnectServer(address string, port int) *P2PNode {
	server_port := strconv.Itoa(port)
	conn, err := net.Dial("tcp", address + ":" + server_port)
	ErrorHandler(err)
	return client.NewConnection(conn)
}

func (client *P2PClient)NewConnection(conn net.Conn) *P2PNode {
	p2p_server := NewNode(conn)
	go client.ConnectionHandler(p2p_server)

	return p2p_server
}

func (client *P2PClient)ConnectionHandler(server *P2PNode) {
	go server.Read()
	go server.Write()
	defer server.Close()

	loop : for {
		select {
		case msg := <-server.incoming :
			switch msg[0] {
			case MSG_RECEIVE_IP :
				fmt.Println("BROADCAST OK")
				break loop
			case MSG_APPROVE_CONN :
				fmt.Println("APPROVE !")
				client.WriteServerMap(server)
				fmt.Println("Gooo : "  + strconv.Itoa(len(client.servers)))
			case MSG_REFUSE_CONN :
				fmt.Println("Refuse....") 
				break loop 
			}

		case state := <-server.state:
			switch state {
			case P2P_DEAD_CONN :
				fmt.Println("DEAD====")
				client.DeleteServerMap(server)
			}
		}
	}
}

func (client *P2PClient)CheckNewConnection(address string) bool{
	if !client.CheckServerMap(address) && !client.p2p_server.CheckClientMap(address) {
		return client.CheckServerMapSize()
	}
	return false
}

func (client *P2PClient)CheckServerMapSize() bool {
    client_map_lock.RLock()
	defer client_map_lock.RUnlock()

	if len(client.servers) < CLIENT_MAX_CONNECTION {
		return true
	}
	
	return false
}


func (client *P2PClient)CheckServerMap(address string) bool {
    client_map_lock.RLock()
	defer client_map_lock.RUnlock()

	_, ok := client.servers[address]
	return ok
}

func (client *P2PClient)WriteServerMap(server *P2PNode) {
	client_map_lock.Lock()
	defer client_map_lock.Unlock()
	client.servers[server.HashKey()] = true
}

func (client *P2PClient)DeleteServerMap(server *P2PNode) {
	client_map_lock.Lock()
	defer client_map_lock.Unlock()
    delete(client.servers, server.HashKey())
}

func (client *P2PClient)IPBroadcast(server *P2PNode) {
	buf := IPBroadcastMsg()
	server.outgoing <- buf
}

func (client *P2PClient)RequestConn(server *P2PNode) {
	buf := RequestConnMsg()
	server.outgoing <- buf
}