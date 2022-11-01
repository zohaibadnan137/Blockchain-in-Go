package assignment01bca

import (
	"time"
)

// ******** ******** STRUCTURES ******** ******** //

type Transaction struct {
	id   byte
	data map[string]interface{}
}

type Node struct {
	transaction *Transaction
	hash        byte

	left  *Node
	right *Node
}

type MerkleTree struct {
	root         *Node
	transactions []Transaction
}

type Block struct {
	hash         byte
	previousHash byte

	merkleTree MerkleTree

	timestamp time.Time
	nonce     byte
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
}

// ******** ******** PRIVATE FUNCTIONS ******** ******** //

// ******** ******** PUBLIC FUNCTIONS ******** ******** //

// Creates a genesis block
func CreateBlockchain() Blockchain {
	genesisBlock := Block{
		timestamp: time.Now(),
	}

	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
	}

}

// Creates a new block*
func NewBlock() {

}

// Finds the nonce value for a block*
func MineBlock() {

}

// Prints all the blocks in the blockchain*
func DisplayBlocks() {

}

// Prints all the transactions in a block*
func DisplayMerkelTree() {

}

// Changes one or multiple transactions in a block*
func ChangeBlock() {

}

// Verifies whether any changes were made in the blockchain*
func VerifyChain() {

}

// Calculates the hash of a transaction or block*
func CalculateHash() {

}
