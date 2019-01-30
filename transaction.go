package main

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
