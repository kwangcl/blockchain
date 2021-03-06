package main

import (
	"net"
	"strconv"
	"sync"
	"log"
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

	log.Println("Log - [P2PServer] New P2P server")

	server_node := NewNode(nil)
	server_node.port = strconv.Itoa(port)
	return &P2PServer{map[*P2PNode]bool{}, server_node, nil}
}

func (server *P2PServer)NewConnection(conn net.Conn) {

	log.Println("Log - [P2PServer] New connection req")

	p2p_client := NewNode(conn)
	go server.ClientHandler(p2p_client)
}


func (server *P2PServer)StartServer() {

	log.Println("Log - [P2PServer] Start P2P server")

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
			log.Println("Log - [P2PServer] Get data from client : " + client.address)

			msg_type := msg[0]
			src := msg[1:41]
			log.Println("Log - [P2PServer] Data source tag : " + string(src))
			switch msg_type {
			case MSG_NEW_NODE :
				log.Println("Log - [P2PServer] New Node!")
				client.outgoing <- MsgManager.ReceiveIPMsg()
				fallthrough
			case MSG_IP_BROADCAST :
				log.Println("Log - [P2PServer] IP broadcast : " + client.address)
				if server.p2p_client.CheckNewConnection(client) {
					log.Println("Log - [P2PServer] Connection request to new node : " + client.address)
					tmp_server := server.p2p_client.ConnectServer(client.address, SERVER_PORT)
					server.p2p_client.RequestConn(tmp_server)
				} else if server.CheckNewConnection(client) {
					log.Println("Log - [P2PServer] sand connection ready msg : " + client.address)
					tmp_server := server.p2p_client.ConnectServer(client.address, SERVER_PORT)
					server.p2p_client.ConnReady(tmp_server)
				} else {
					log.Println("Log - [P2PServer] Connection full or duplicated")
				}
				new_msg := MsgManager.IPBroadcastMsg(src)
				server.BroadCastMsg(new_msg, client.address)
				server.p2p_client.BroadCastMsg(new_msg, client.address)

			case MSG_REQUEST_CONN :
				log.Println("Log - [P2PServer] Reqeust connection : " + client.address)
				if server.CheckNewConnection(client) {

					client.outgoing <- MsgManager.ApproveConnMsg()
					server.WriteClientMap(client)
				} else {
					client.outgoing <- MsgManager.RefuseConnMsg()
				}
				case MSG_SEND_TRANSACTION :
					log.Println("Log - [P2PServer] New Transaction from client!")
					MsgManager.TransactionBroadCast(msg)
					fallthrough
				case MSG_TRANSACTION_BROADCAST :
					log.Println("Log - [P2PServer] Transaction broadcast")
					tx := DeserializeTx(msg[41:])
					tx.PrintTxData()
					server.BroadCastMsg(msg, client.address)
					server.p2p_client.BroadCastMsg(msg, client.address)

			}
		case state := <-client.state:
			switch state {
			case P2P_DEAD_CONN :
				log.Println("Log - [P2PServer] Client connection dead : " + client.address)
				server.DeleteClientMap(client)
				break loop
			case P2P_DUP_MSG :
				log.Println("Log - [P2PServer] Duplicated msg")
			}
		}
	}
}


func (server *P2PServer)CheckNewConnection(client *P2PNode) bool {
	if server.CheckClientMap(client.address) && server.p2p_client.CheckServerMap(client.address) {
		return server.CheckClientMapSize()
  }
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

func (server *P2PServer)BroadCastMsg(msg []byte, src string) {
	log.Println("Log -[P2PServer] BroadCastMsg : ")
	for client, _ := range server.clients {
		if client.address != src {
			log.Println("Log -[P2PServer] Client IP - " + client.address)
			client.outgoing <- msg
		}
	}
	//log.Println(msg)
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

func (server *P2PServer)PrintServerMap() {
	for client, _ := range server.clients {
		log.Println("Log - [P2PServer] Print server map : client IP - " + client.address)
	}
}
