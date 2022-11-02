package assignment01bca

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// ******** ******** STRUCTURES ******** ******** //

type Transaction struct {
	id        byte
	data      map[string]interface{}
	timestamp time.Time
}

type Node struct {
	transaction *Transaction
	hash        string

	left   *Node
	right  *Node
	parent *Node

	isLeaf bool
}

type MerkleTree struct {
	root         Node
	transactions []Transaction
}

type Block struct {
	hash         string
	previousHash string

	previousBlock *Block
	nextBlock     *Block

	merkleTree MerkleTree

	timestamp time.Time
	nonce     int
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
}

// ******** ******** PRIVATE FUNCTIONS ******** ******** //

func createMerkleTree(transactions []Transaction) MerkleTree {

	numberOfTransactions := len(transactions) // Get the number of transactions
	newNumberOfTransactions := numberOfTransactions

	// If the number of transactions is odd, the size of the new array must be increased by one
	if numberOfTransactions%2 != 0 {
		newNumberOfTransactions = numberOfTransactions + 1
	}

	finalTransactions := make([]Transaction, newNumberOfTransactions) // Create a new array to store the transactions

	// Copy the transactions from the original array to the new one
	for i := 0; i < numberOfTransactions; i++ {
		finalTransactions[i] = transactions[i]
	}

	// If the number of transactions is odd, duplicate the last transaction
	if numberOfTransactions%2 != 0 {
		finalTransactions[len(finalTransactions)-1] = finalTransactions[len(finalTransactions)-2]
	}

	// Create a new leaf node for each transaction
	leafNodes := make([]Node, newNumberOfTransactions)
	for i := 0; i < newNumberOfTransactions; i++ {
		leafNodes[i].transaction = &finalTransactions[i]
		leafNodes[i].isLeaf = true

		data, error := json.Marshal(finalTransactions[i].data)
		if error != nil {
		}
		hash := sha256.Sum256(data)
		leafNodes[i].hash = string(hash[:])
	}

	// Create the Merkle tree bottom-up
	numberOfParentNodes := newNumberOfTransactions / 2 // Nodes on the new higher level that is being created
	//numberOfChildNodes := newNumberOfTransactions      // Nodes on the previous lower level
	childNodes := leafNodes

	var rootNode Node

	for numberOfParentNodes >= 1 {
		parentNodes := make([]Node, numberOfParentNodes)
		childNodeIndex := 0

		for i := 0; i < numberOfParentNodes; i++ {
			// Set the current node as the parent for the two respective children
			childNodes[childNodeIndex].parent = &parentNodes[i]
			childNodes[childNodeIndex+1].parent = &parentNodes[i]

			// Set the current node's left and right child respectively
			parentNodes[i].left = &childNodes[childNodeIndex]
			parentNodes[i].right = &childNodes[childNodeIndex+1]

			// Calculate the current node's hash
			hash := sha256.Sum256([]byte(parentNodes[i].left.hash + parentNodes[i].right.hash))
			parentNodes[i].hash = string(hash[:])

			// Move to the next two children
			childNodeIndex += 2

			if numberOfParentNodes == 1 {
				rootNode = parentNodes[i]
			}
		}

		childNodes = parentNodes
		numberOfParentNodes /= 2
	}

	return MerkleTree{
		rootNode,
		finalTransactions,
	}
}

/*func breadthFirstSearch(merkleTree MerkleTree) {

}*/

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
		time.Now(),
	}
}

// Creates a new block*
func NewBlock(blockchain *Blockchain, transactions []Transaction) Block {
	// Create a Merkle Tree
	merkleTree := createMerkleTree(transactions)

	return Block{
		hash:          merkleTree.root.hash,
		previousHash:  blockchain.chain[len(blockchain.chain)-1].hash,
		previousBlock: &blockchain.chain[len(blockchain.chain)-1],
		merkleTree:    merkleTree,
		timestamp:     time.Now(),
	}
}

// Finds the nonce value for a block*
func MineBlock(blockchain *Blockchain, block *Block) {

}

// Prints all the blocks in the blockchain*
func DisplayBlocks(blockChain Blockchain) {
	for i := 0; i < len(blockChain.chain); i++ {
		fmt.Println("")
		fmt.Println("--------------------------------")
		fmt.Println("BLOCK HASH: ", blockChain.chain[i].hash)
		fmt.Println("PREVIOUS BLOCK HASH: ", blockChain.chain[i].hash)
		fmt.Println("TIMESTAMP: ", blockChain.chain[i].timestamp.Format(time.RFC822))
		fmt.Println("NONCE: ", blockChain.chain[i].nonce)
		fmt.Println("--------------------------------")
	}
}

// Prints all the transactions in a block*
func DisplayMerkelTree(block Block) {
	if len(block.merkleTree.transactions) == 0 {
		fmt.Println("This block does not have any transactions")
	} else {
		for i := 0; i < len(block.merkleTree.transactions); i++ {
			fmt.Println("")
			fmt.Println("--------------------------------")
			fmt.Println("ID: ", block.merkleTree.transactions[i].id)
			fmt.Println("TIMESTAMP: ", block.merkleTree.transactions[i].timestamp.Format(time.RFC822))

			data, error := json.Marshal(block.merkleTree.transactions[i].data)
			if error != nil {
			}
			fmt.Println("DATA: ", string(data))
			fmt.Println("--------------------------------")
		}
	}
}

// Changes one or multiple transactions in a block*
func ChangeBlock() {

}

// Verifies whether any changes were made in the blockchain*
func VerifyChain() {

}

// Calculates the hash of a transaction or block*
func CalculateHash(block Block) string {
	return "NULL"
}
