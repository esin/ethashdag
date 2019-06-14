package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereumproject/go-ethereum/crypto/sha3"
)

const (
	epochLength     uint64 = 30000
	maxBlocksNumber uint64 = 2048
)

func getSeedHash(blockNum uint64) ([]byte, error) {
	if blockNum >= epochLength*maxBlocksNumber {
		return nil, fmt.Errorf("block number too high, limit is %d", epochLength*2048)
	}
	sh := makeSeedHash(blockNum / epochLength)
	return sh[:], nil
}

func keccak256Hash(data ...[]byte) (h common.Hash) {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

func makeSeedHash(epoch uint64) (sh common.Hash) {
	for ; epoch > 0; epoch-- {
		sh = keccak256Hash(sh[:])
	}
	return sh
}

func getDagFileName(currentBlock uint64) (filename string) {
	seedhash, err := getSeedHash(currentBlock)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("full-R23-%s", hex.EncodeToString(seedhash)[:16])
}

func main() {

	var currentBlock uint64
	currentBlock = 0

	fmt.Println("| Epoch # | Epoch first block #| Epoch last block #| Filename |")
	fmt.Println("| --- | --- | --- | --- |")
	var i uint64
	epoch := 0
	for i = 0; i < epochLength*(maxBlocksNumber-1)+1; i = i + epochLength {
		fmt.Printf("|   %d   |   %d   |   %d   |   %s   |\n", epoch, i, i+epochLength-1, getDagFileName(currentBlock+i))
		epoch = epoch + 1
	}
}
