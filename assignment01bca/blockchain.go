package assignment01bca

import (
	"bytes"
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

type MerkleTreeNode struct {
	transaction *Transaction
	hash        [32]byte

	left   *MerkleTreeNode
	right  *MerkleTreeNode
	parent *MerkleTreeNode

	isLeaf bool
}

type MerkleTree struct {
	root         MerkleTreeNode
	transactions []Transaction
}

type Block struct {
	hash         [32]byte // Merely represents the hash of the block's Merkle tree's root
	previousHash [32]byte

	previousBlock *Block
	nextBlock     *Block

	merkleTree MerkleTree

	timestamp time.Time
	nonce     [4]byte // A random four-byte number that is calculated when the block is mined
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block

	difficulty int // Determines the number of required leading zeroes while calculating nonce and the subsequent hash
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

	// Create a new leaf MerkleTreeNode for each transaction
	leafMerkleTreeNodes := make([]MerkleTreeNode, newNumberOfTransactions)
	for i := 0; i < newNumberOfTransactions; i++ {
		leafMerkleTreeNodes[i].transaction = &finalTransactions[i]
		leafMerkleTreeNodes[i].isLeaf = true

		data, error := json.Marshal(finalTransactions[i].data)
		if error != nil {
			fmt.Println("Error! The hash for transaction ", finalTransactions[i].id, " cannot be calculated.")
		} else {
			hash := sha256.Sum256(data)
			leafMerkleTreeNodes[i].hash = hash
		}
	}

	// Create the Merkle tree bottom-up
	numberOfParentMerkleTreeNodes := newNumberOfTransactions / 2 // MerkleTreeNodes on the new higher level that is being created
	//numberOfChildMerkleTreeNodes := newNumberOfTransactions      // MerkleTreeNodes on the previous lower level
	childMerkleTreeNodes := leafMerkleTreeNodes

	var rootMerkleTreeNode MerkleTreeNode

	for numberOfParentMerkleTreeNodes >= 1 {
		parentMerkleTreeNodes := make([]MerkleTreeNode, numberOfParentMerkleTreeNodes)
		childMerkleTreeNodeIndex := 0

		for i := 0; i < numberOfParentMerkleTreeNodes; i++ {
			// Set the current MerkleTreeNode as the parent for the two respective children
			childMerkleTreeNodes[childMerkleTreeNodeIndex].parent = &parentMerkleTreeNodes[i]
			childMerkleTreeNodes[childMerkleTreeNodeIndex+1].parent = &parentMerkleTreeNodes[i]

			// Set the current MerkleTreeNode's left and right child respectively
			parentMerkleTreeNodes[i].left = &childMerkleTreeNodes[childMerkleTreeNodeIndex]
			parentMerkleTreeNodes[i].right = &childMerkleTreeNodes[childMerkleTreeNodeIndex+1]

			// Calculate the current MerkleTreeNode's hash
			hash := sha256.Sum256(append(parentMerkleTreeNodes[i].left.hash[:], parentMerkleTreeNodes[i].right.hash[:]...))
			parentMerkleTreeNodes[i].hash = hash

			// Move to the next two children
			childMerkleTreeNodeIndex += 2

			if numberOfParentMerkleTreeNodes == 1 {
				rootMerkleTreeNode = parentMerkleTreeNodes[i]
			}
		}

		childMerkleTreeNodes = parentMerkleTreeNodes
		numberOfParentMerkleTreeNodes /= 2
	}

	return MerkleTree{
		rootMerkleTreeNode,
		finalTransactions,
	}
}

// Adds the given block into the blockchain
func (blockchain *Blockchain) addBlock(block *Block) {
	var previousBlock *Block = &blockchain.chain[len(blockchain.chain)-1]
	previousBlock.nextBlock = block
	block.previousBlock = previousBlock
	blockchain.chain = append(blockchain.chain, *block)
}

// ******** ******** PUBLIC FUNCTIONS ******** ******** //

// Creates a new blockchain along with the respective genesis block
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		timestamp: time.Now(),
	}

	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
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
func (blockchain *Blockchain) MineBlock(block *Block) {
	// Concatenate the previous block's hash at the end of the current block's hash
	concatenated_hashes := append(block.hash[:], block.previousHash[:]...)
	flag := false     // The flag will be set to true once the nonce has been found
	var nonce [4]byte // The nonce is a byte array of length four

	for !flag {
		rand.Read(nonce[:])                                                // Generate a random byte
		concatenated_hashes := append(concatenated_hashes[:], nonce[:]...) // Add the nonce at the end of the concatenated hashes
		hash := sha256.Sum256(concatenated_hashes[:])                      // Calculate the new hash

		flag = true
		// Check whether the new hash has the required number of leading zeroes
		for count := 0; count < blockchain.difficulty; count++ {
			if hash[count] > 0 {
				flag = false
				break
			}
		}
	}

	fmt.Println(nonce)
	block.nonce = nonce        // Add the nonce to the mined block
	blockchain.addBlock(block) // Add the block to the blockchain
}

// Prints all the blocks in the blockchain*
func (blockchain Blockchain) DisplayBlocks() {
	for i := 0; i < len(blockchain.chain); i++ {
		fmt.Println("\n********************************")
		fmt.Println("// BLOCK HASH: ", blockchain.chain[i].hash)
		fmt.Println("// PREVIOUS BLOCK HASH: ", blockchain.chain[i].previousHash)
		fmt.Println("// TIMESTAMP: ", blockchain.chain[i].timestamp.Format(time.RFC822))
		fmt.Println("// NONCE: ", blockchain.chain[i].nonce)
		fmt.Println("********************************")
	}
}

// Prints all the transactions in a block*
func DisplayMerkelTree(block Block) {
	if len(block.merkleTree.transactions) == 0 {
		fmt.Println("This block does not have any transactions")
	} else {
		for i := 0; i < len(block.merkleTree.transactions); i++ {
			fmt.Println("")
			fmt.Println("////////////////////////////////////////////////////////////////")
			fmt.Println("// ID: ", block.merkleTree.transactions[i].id)
			fmt.Println("// TIMESTAMP: ", block.merkleTree.transactions[i].timestamp.Format(time.RFC822))

			data, error := json.Marshal(block.merkleTree.transactions[i].data)
			if error != nil {
				fmt.Println("Error! The data for this transaction cannot be displayed.")
			} else {
				fmt.Println("// DATA: ", string(data))
			}
			fmt.Println("////////////////////////////////////////////////////////////////")
		}
	}
	fmt.Println("")
}

// Changes one or multiple transactions in a block*
func (block *Block) ChangeBlock() {
}

// Verifies whether any changes were made in the blockchain*
func (blockchain Blockchain) VerifyChain() {
	modified := false
	currentBlock := blockchain.genesisBlock

	var modifiedBlock Block

	// Start iterating from the genesis block and compare the current hash of each block with its stored hash
	for currentBlock.nextBlock != nil {
		currentHash := currentBlock.CalculateHash()
		storedHash := currentBlock.hash

		if bytes.Equal(currentHash[:], storedHash[:]) {
			modified = true
			modifiedBlock = currentBlock
			break
		}

		currentBlock = *currentBlock.nextBlock
	}

	if modified {
		fmt.Println("The blockchain has been modified. Block ", modifiedBlock.hash, " has been modified.")
	} else {
		fmt.Println("The blockchain has not been modified.")
	}
}

// Calculates the hash of a transaction or block*
func (transaction Transaction) CalculateHash() [32]byte {
	data, error := json.Marshal(transaction.data)
	if error != nil {
		fmt.Println("Error! The hash for this transaction cannot be calculated.")
		return [32]byte{}
	}
	hash := sha256.Sum256(data)
	return hash
}

func (block Block) CalculateHash() [32]byte {
	// Create a temporary Merkle tree for the transactions in the block and return the hash of the root
	tempMerkleTree := createMerkleTree(block.merkleTree.transactions)
	hash := tempMerkleTree.root.hash
	return hash
}

func (blockchain *Blockchain) Test(block Block) {
}
