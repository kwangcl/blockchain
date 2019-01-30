package main

import (
	"net"
	"strconv"
	"log"
	"sync"
)

var client_map_lock = sync.RWMutex{}
var CLIENT_MAX_CONNECTION = 2


type P2PClient struct {
	servers map[*P2PNode]bool
	node	*P2PNode
	p2p_server *P2PServer
}

func NewP2PClient() *P2PClient {

	log.Println("Log - [P2PClient] New P2P client")

	client_node := NewNode(nil)
	return &P2PClient{map[*P2PNode]bool{}, client_node, nil}
}

func (client *P2PClient)ConnectServer(address string, port int) *P2PNode {

	log.Println("Log - [P2PClient] Connect server : " + address)

	server_port := strconv.Itoa(port)
	conn, err := net.Dial("tcp", address + ":" + server_port)
	ErrorHandler(err)
	return client.NewConnection(conn)
}

func (client *P2PClient)NewConnection(conn net.Conn) *P2PNode {

	log.Println("Log - [P2PClient] Request connection to server")

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
			log.Println("Log - [P2PClient] Get data from server : " + server.address)

			msg_type := msg[0]
			src := msg[1:41]
			log.Println("Log - [P2PClient] Data source tag : " + string(src))

			switch msg_type {
			case MSG_RECEIVE_IP :
				log.Println("Log - [P2PClient] Receive IP & break connection")
				break loop
			case MSG_APPROVE_CONN :
				log.Println("Log - [P2PClient] Approve connection from server : " + server.address)
				client.WriteServerMap(server)
			case MSG_REFUSE_CONN :
			  log.Println("Log - [P2PClient] Refuse connection from server : " + server.address)
				break loop
			}

		case state := <-server.state:
			switch state {
			case P2P_DEAD_CONN :
				log.Println("Log - [P2PClient] Server connection dead : " + server.address)
				client.DeleteServerMap(server)
			case P2P_DUP_MSG :
				log.Println("Log - [P2PServer] Duplicated msg")
			}
		}
	}
}

func (client *P2PClient)CheckNewConnection(server *P2PNode) bool{
	if client.CheckServerMap(server.address) && client.p2p_server.CheckClientMap(server.address) {
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


func (client *P2PClient)BroadCastMsg(msg []byte, src string) {
	log.Println("Log -[P2PClient] BroadCastMsg : ")
	for server, _ := range client.servers {
		if server.address != src {
			log.Println("Log -[P2PClient] Server IP - " + server.address)
			server.outgoing <- msg
		}
	}
	log.Println(msg)
}


func (client *P2PClient)CheckServerMap(address string) bool {
    client_map_lock.RLock()
	defer client_map_lock.RUnlock()

	for server, _ := range client.servers {
		if server.address == address {
			return false
		}
	}
	return true
}

func (client *P2PClient)WriteServerMap(server *P2PNode) {
	client_map_lock.Lock()
	defer client_map_lock.Unlock()
	client.servers[server] = true
}

func (client *P2PClient)DeleteServerMap(server *P2PNode) {
	client_map_lock.Lock()
	defer client_map_lock.Unlock()
    delete(client.servers, server)
}


func (client *P2PClient)BroadCastNewNode(server *P2PNode) {
	buf := MsgManager.NewNodeMsg()
	server.outgoing <- buf
}

/*
func (client *P2PClient)IPBroadcast(server *P2PNode) {
	buf := MsgManager.IPBroadcastMsg()
	server.outgoing <- buf
}*/


func (client *P2PClient)RequestConn(server *P2PNode) {
	buf := MsgManager.RequestConnMsg()
	server.outgoing <- buf
}

func (client *P2PClient)ConnReady(server *P2PNode) {
	buf := MsgManager.ConnReadyMsg()
	server.outgoing <- buf
}

func (client *P2PClient)PrintClientMap() {
	for server, _ := range client.servers {
		log.Println("Log - [P2PClinet] Print server map : server IP - " + server.address)
	}
}
