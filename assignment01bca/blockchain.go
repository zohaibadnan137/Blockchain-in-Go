package assignment01bca

import (
	"math/rand"
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
	root         Node
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

// Adds the given node at the end of the specified Merkle tree
func addNode(merkleTree *MerkleTree, node Node) {

}

// ******** ******** PUBLIC FUNCTIONS ******** ******** //

// Creates a new blockchain along with the respective genesis block
func CreateBlockchain() Blockchain {
	genesisBlock := Block{
		timestamp: time.Now(),
	}

	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
	}
}

// Creates a new transaction
func CreateTransaction(from, to string, amount byte) Transaction {

	data := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}

	return Transaction{
		byte(rand.Int()),
		data,
	}
}

// Creates a new block*
func NewBlock(blockchain *Blockchain, transactions []Transaction) Block {
	// Create a Merkle Tree
	//merkleTree := MerkleTree{
	//transactions: transactions,
	//}

	// Add a new node in the Merkle tree for each transaction

	block := Block{}
	return block
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
