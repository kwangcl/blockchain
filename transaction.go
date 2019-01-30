package main

import (
  "bytes" 
  "encoding/gob"
)
const TRANSACTION_MAX_DATA_SIZE = 256
const TRANSACTION_MAX_MSG_SIZE = 128

type Transaction struct {
  data [TRANSACTION_MAX_DATA_SIZE]byte
	msg  [TRANSACTION_MAX_MSG_SIZE]byte
}

func CreateTransaction(data []byte, msg []byte) *Transaction{
	tx := Transaction{[TRANSACTION_MAX_DATA_SIZE]byte{}, [TRANSACTION_MAX_MSG_SIZE]byte{}}
	copy(tx.data[:], data)
	copy(tx.msg[:], msg)
	return &tx
}

func (tx *Transaction) Serialize() []byte {
  var result bytes.Buffer

  encoder := gob.NewEncoder(&result)
  err := encoder.Encode(tx)
  ErrorHandler(err)

  return tx.Bytes()
}

func DeserializeTx(data []byte) *Transaction {
  var tx Transaction

  decoder := gob.NewDecoder(bytes.NewReader(data))
  err := decoder.Decode(&tx)
  ErrorHandler(err)

  return *tx
}
