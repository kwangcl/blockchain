package main

import (
	_"bytes"
	"math"
	"math/big"
)


var MAX_NONCE = math.MaxInt64

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-8))

	return &ProofOfWork{b, target}
}
/*
func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join
}*/
