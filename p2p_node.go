package main

import (
	"net"
	"fmt"
)

type P2PNetState int

const (
	P2P_INIT = iota
	P2P_DEAD_CONN
	P2P_WAIT_CONN
)


type P2PNode struct {
	address string
	port string
	conn net.Conn
	incoming chan []byte
	outgoing chan []byte
	state	 chan P2PNetState
}

func NewNode(conn net.Conn) *P2PNode {
	var address, port string

	if conn != nil {
		client_info := conn.RemoteAddr().String()
		client_address, client_port, err := net.SplitHostPort(client_info)
		ErrorHandler(err)

		address = client_address
		port = client_port

	} else {
		address = ""
		port = ""
	}
	incoming := make(chan []byte)
	outgoing := make(chan []byte)
	state := make(chan P2PNetState)

	return &P2PNode{address, port, conn, incoming, outgoing, state}
}

func (node *P2PNode)Read() {
	for {
		buf := make([]byte, 1024)
		_, err := node.conn.Read(buf)
		if !node.ConnectionCheck(err) {
			node.state <- P2P_DEAD_CONN
			break
		} else {
			if MsgManager.CheckNewMsg(buf[1:41]) {
				node.incoming <- buf
			}
		}
	}
}

func (node *P2PNode)Write() {
	for data := range node.outgoing {
		_, err := node.conn.Write(data)
		ErrorHandler(err)
	}
}


func (node *P2PNode) Close() {
	err := node.conn.Close()
	ErrorHandler(err)
}

func (node *P2PNode) ConnectionCheck(err error) bool {
	if err != nil {
		return false
	}
	return true
}

func (node *P2PNode) HashKey() string {
	return node.address
}
