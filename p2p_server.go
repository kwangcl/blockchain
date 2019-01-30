package main

import (
	"net"
	"strconv"
	"sync"
	"fmt"
)

var server_map_lock = sync.RWMutex{}

const SERVER_MAX_CONNECTION = 2
const SERVER_PORT = 7777

type P2PServer struct {
	clients map[*P2PNode] bool
	node	*P2PNode
	p2p_client *P2PClient
}


func NewP2PServer(port int) *P2PServer {

	server_node := NewNode(nil)
	server_node.port = strconv.Itoa(port)
	return &P2PServer{map[*P2PNode]bool{}, server_node, nil}
}

func (server *P2PServer)NewConnection(conn net.Conn) {

	p2p_client := NewNode(conn)
	//server.WriteClientMap(p2p_client)
	go server.ClientHandler(p2p_client)
}


func (server *P2PServer)StartServer() {
	listener, err := net.Listen("tcp", ":" + server.node.port)
	ErrorHandler(err)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		ErrorHandler(err)
		server.NewConnection(conn)
	}
}


func (server *P2PServer)ClientHandler(client *P2PNode) {

	go client.Read()
	go client.Write()
	defer client.Close()

	loop : for {
		select {
		case msg := <-client.incoming:
			//fmt.Println("size : " + strconv.Itoa(len(server.clients)))
			//fmt.Println(client.address + " : " + client.port + " : " + msg[0])
			//fmt.Println(msg)
			msg_type := msg[0]
			//src := msg[1:41]

			switch msg_type {
			case MSG_IP_BROADCAST :
				fmt.Println("IP BROADCAST : " + client.address + ", " +  client.port)
				client.outgoing <- MsgManager.ReceiveIPMsg()

				if server.p2p_client.CheckNewConnection(client) {
					tmp_server := server.p2p_client.ConnectServer(client.address, SERVER_PORT)
					server.p2p_client.RequestConn(tmp_server)
				} else {
					fmt.Println("Already")
				}

			case MSG_REQUEST_CONN :
				if server.CheckNewConnection(client) {
					client.outgoing <- MsgManager.ApproveConnMsg()
					server.WriteClientMap(client)
				} else {
					client.outgoing <- MsgManager.RefuseConnMsg()
				}
			}
		case state := <-client.state:
			if state == P2P_DEAD_CONN {
				fmt.Println(client.address + " : " + client.port +  " = dead")
				server.DeleteClientMap(client)
				break loop
			}
		}
	}
}


func (server *P2PServer)CheckNewConnection(client *P2PNode) bool {
	if !server.CheckClientMap(client.address) && !server.p2p_client.CheckServerMap(client.address) {
		fmt.Println("!---")
		return server.CheckClientMapSize()
  	}
	fmt.Println("33333")
	return false
}

func (server *P2PServer)CheckClientMapSize() bool {
	server_map_lock.RLock()
	defer server_map_lock.RUnlock()

	if len(server.clients) < SERVER_MAX_CONNECTION {
		return true
	}
	return false
}

func (server *P2PServer)BroadCastMsg(msg []byte) {
	for client, _ := range server.clients {
			client.outgoing <- msg
	}
}


func (server *P2PServer)CheckClientMap(address string) bool {
	server_map_lock.RLock()
	defer server_map_lock.RUnlock()

	for client, _ := range server.clients {
		if client.address == address {
			return false
		}
	}
	return true
}

func (server *P2PServer)WriteClientMap(client *P2PNode) {
	server_map_lock.Lock()
	defer server_map_lock.Unlock()
	server.clients[client] = true
}

func (server *P2PServer)DeleteClientMap(client *P2PNode) {
	server_map_lock.Lock()
	defer server_map_lock.Unlock()
	delete(server.clients, client)
}
