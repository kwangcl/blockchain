package main

type block struct {
	Tx []Transaction
	PrevBlockHash []byte
	Hash	[]byte
	Nonce	int
}
