package main

import (
	"sync"
	"strconv"
	"time"
)

type MsgState int

const MAX_MSG_BUF_SIZE = 1024
const MAX_SRC_BUF_SIZE = 1000
const (
	MSG_NEW_NODE = iota
	MSG_RECEVIE_NEW_INFO
	MSG_IP_BROADCAST
	MSG_RECEIVE_IP
	MSG_REQUEST_CONN
	MSG_APPROVE_CONN
	MSG_REFUSE_CONN
	MSG_CONN_READY
	MSG_SEND_TRANSACTION
	//MSG
	//MSG
	//MSG
	//MSG
)


type P2PMsgManager struct {
	src_buffer *P2PSrcBuffer
}

type P2PSrcBuffer struct {
	queue []string
	msg_map map[string]bool
	address string
}

var src_buffer_lock = sync.RWMutex{}

var MsgManager *P2PMsgManager

func NewP2PMsgManager() *P2PMsgManager{
	address := GetIPv4Address()
	src_buffer := &P2PSrcBuffer{[]string{}, map[string]bool{}, address}
	return &P2PMsgManager{src_buffer}
}


func (msg_manager *P2PMsgManager)CheckNewMsg(src []byte) bool {
	str := string(src[:])
	return msg_manager.src_buffer.CheckSrcBuf(str)
}

func (src_buffer *P2PSrcBuffer)CheckSrcDuplicate(src string) bool {
	src_buffer_lock.RLock()
	defer src_buffer_lock.RUnlock()
	_, ok := src_buffer.msg_map[src]
	return ok
}

func (src_buffer *P2PSrcBuffer)CheckSrcBuf(src string) bool {
	if src_buffer.CheckSrcDuplicate(src) {
		return false
	}
	src_buffer_lock.Lock()
	defer src_buffer_lock.Unlock()
	if len(src_buffer.queue) == MAX_SRC_BUF_SIZE {
		src_buffer.queue = src_buffer.queue[100:]
	}
	src_buffer.queue = append(src_buffer.queue, src)
	src_buffer.msg_map[src] = true

	return true
}



func (msg_manager *P2PMsgManager)NewNodeMsg() []byte {
	buf := make([]byte, MAX_MSG_BUF_SIZE)
	buf[0] = byte(MSG_NEW_NODE)
	src_buf := msg_manager.GenSrcData()
	copy(buf[1:], src_buf)
	return buf
}


func (msg_manager *P2PMsgManager)IPBroadcastMsg(src []byte) []byte {
	buf := make([]byte, MAX_MSG_BUF_SIZE)
	buf[0] = byte(MSG_IP_BROADCAST)
	copy(buf[1:41], src)
	return buf
}

func (msg_manager *P2PMsgManager)ReceiveIPMsg() []byte {
	buf := make([]byte, MAX_MSG_BUF_SIZE)
	buf[0] = byte(MSG_RECEIVE_IP)
	src_buf := msg_manager.GenSrcData()
	copy(buf[1:], src_buf)
	return buf
}

func (msg_manager *P2PMsgManager)RequestConnMsg() []byte {
	buf := make([]byte, MAX_MSG_BUF_SIZE)
	buf[0] = byte(MSG_REQUEST_CONN)
	src_buf := msg_manager.GenSrcData()
	copy(buf[1:], src_buf)
	return buf
}

func (msg_manager *P2PMsgManager)ApproveConnMsg() []byte {
	buf := make([]byte, MAX_MSG_BUF_SIZE)
	buf[0] = byte(MSG_APPROVE_CONN)
	src_buf := msg_manager.GenSrcData()
	copy(buf[1:], src_buf)
	return buf
}

func (msg_manager *P2PMsgManager)RefuseConnMsg() []byte {
	buf := make([]byte, MAX_MSG_BUF_SIZE)
	buf[0] = byte(MSG_REFUSE_CONN)
	src_buf := msg_manager.GenSrcData()
	copy(buf[1:], src_buf)
	return buf
}


func (msg_manager *P2PMsgManager)ConnReadyMsg() []byte {
	buf := make([]byte, MAX_MSG_BUF_SIZE)
	buf[0] = byte(MSG_CONN_READY)
	src_buf := msg_manager.GenSrcData()
	copy(buf[1:], src_buf)
	return buf
}

func (msg_manager *P2PMsgManager)SendTransactionMsg(data []byte, new_src []byte) []byte {
	buf := make([]byte, len(data) + 1)
	buf[0] = byte(MSG_SEND_TRANSACTION)
	src_buf := msg_manager.GenSrcData()
	copy(buf[1:], src_buf)

	copy(buf[41:], data[:])
	return buf
}

func (msg_manager *P2PMsgManager)GenSrcData() []byte {
	str := msg_manager.src_buffer.address + "-"
	str += strconv.FormatInt(time.Now().UnixNano(), 10)
	buf := make([]byte, 40)
	copy(buf[:], str)

	if len(msg_manager.src_buffer.queue) == MAX_SRC_BUF_SIZE {
		msg_manager.src_buffer.queue = src_buffer.queue[100:]
	}
	msg_manager.src_buffer.queue = append(src_buffer.queue, src)
	msg_manager.src_buffer.msg_map[src] = true

	return buf[:]
}
