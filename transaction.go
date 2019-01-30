package main

import (
  "bytes"
  "encoding/gob"
  "log"
)
const TRANSACTION_MAX_DATA_SIZE = 256
const TRANSACTION_MAX_MSG_SIZE = 128

type Transaction struct {
  Data [TRANSACTION_MAX_DATA_SIZE]byte
	Msg  [TRANSACTION_MAX_MSG_SIZE]byte
}

func CreateTransaction(data []byte, msg []byte) *Transaction{

  log.Println("Log - [Transaction] Create Transaction")
  log.Println("Log - [Transaction] Data : " + string(data[:]))
  log.Println("Log - [Transaction] Msg : " + string(msg[:]))

	tx := Transaction{[TRANSACTION_MAX_DATA_SIZE]byte{}, [TRANSACTION_MAX_MSG_SIZE]byte{}}
	copy(tx.Data[:], data)
	copy(tx.Msg[:], msg)
	return &tx
}

func (tx *Transaction) PrintTxData() {
  log.Println("Log - [Transaction] Print transaction")
  log.Println("Log - [Transaction] Data : " + string(tx.Data[:]))
  log.Println("Log - [Transaction] Msg : " + string(tx.Msg[:]))
}


func (tx *Transaction) Serialize() []byte {
  var result bytes.Buffer

  encoder := gob.NewEncoder(&result)
  err := encoder.Encode(tx)
  ErrorHandler(err)

  return result.Bytes()
}


func DeserializeTx(data []byte) *Transaction {
  var tx Transaction

  decoder := gob.NewDecoder(bytes.NewReader(data))
  err := decoder.Decode(&tx)
  ErrorHandler(err)

  return &tx
}
