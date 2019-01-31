package main

import (
	_"bytes"
	_"crypto/sha256"
)

type Block struct {
	Transactions []Transaction
	PrevBlockHash []byte
	Hash	[]byte
	Nonce	int
}
/*
func (b *Block) HashTransactions() []byte {
	var tx_hashes [][]byte
	var tx_hash [32]byte

	for _, tx := range b.Transactions {
		tx_hashes  = append(tx_hashes, tx.Serialize())
	}
	return sha256.Sum256(bytes.Join(tx_hashes, []byte{}))[:]
}*/
