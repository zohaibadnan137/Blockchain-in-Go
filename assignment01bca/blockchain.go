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
	ID        byte
	DATA      map[string]interface{}
	TIMESTAMP time.Time
}

type MerkleTreeNode struct {
	TRANSACTION *Transaction
	HASH        [32]byte

	LEFT   *MerkleTreeNode
	RIGHT  *MerkleTreeNode
	PARENT *MerkleTreeNode

	ISLEAFNODE bool
}

type MerkleTree struct {
	ROOT         MerkleTreeNode
	TRANSACTIONS []Transaction
}

type Block struct {
	HASH         [32]byte // Merely represents the hash of the block's Merkle tree's root
	PREVIOUSHASH [32]byte

	PREVIOUSBLOCK *Block
	NEXTBLOCK     *Block

	MERKLETREE MerkleTree

	TIMESTAMP time.Time
	NONCE     [4]byte // A random four-byte number that is calculated when the block is mined
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block

	difficulty int // Determines the number of required leading zeroes while calculating nonce and the subsequent hash
}

// ******** ******** PRIVATE FUNCTIONS ******** ******** //

// ******** ******** PUBLIC FUNCTIONS ******** ******** //

// Creates a new blockchain along with the respective genesis block
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		TIMESTAMP: time.Now(),
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

func CreateMerkleTree(transactions []Transaction) MerkleTree {

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
		leafMerkleTreeNodes[i].TRANSACTION = &finalTransactions[i]
		leafMerkleTreeNodes[i].ISLEAFNODE = true

		data, error := json.Marshal(finalTransactions[i].DATA)
		if error != nil {
			fmt.Println("Error! The hash for transaction ", finalTransactions[i].ID, " cannot be calculated.")
		} else {
			hash := sha256.Sum256(data)
			leafMerkleTreeNodes[i].HASH = hash
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
			childMerkleTreeNodes[childMerkleTreeNodeIndex].PARENT = &parentMerkleTreeNodes[i]
			childMerkleTreeNodes[childMerkleTreeNodeIndex+1].PARENT = &parentMerkleTreeNodes[i]

			// Set the current MerkleTreeNode's left and right child respectively
			parentMerkleTreeNodes[i].LEFT = &childMerkleTreeNodes[childMerkleTreeNodeIndex]
			parentMerkleTreeNodes[i].RIGHT = &childMerkleTreeNodes[childMerkleTreeNodeIndex+1]

			// Calculate the current MerkleTreeNode's hash
			hash := sha256.Sum256(append(parentMerkleTreeNodes[i].LEFT.HASH[:], parentMerkleTreeNodes[i].RIGHT.HASH[:]...))
			parentMerkleTreeNodes[i].HASH = hash

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

// Creates a new block*
func NewBlock(blockchain *Blockchain, transactions []Transaction) Block {
	// Create a Merkle Tree
	merkleTree := CreateMerkleTree(transactions)

	return Block{
		HASH:          merkleTree.ROOT.HASH,
		PREVIOUSHASH:  blockchain.chain[len(blockchain.chain)-1].HASH,
		PREVIOUSBLOCK: &blockchain.chain[len(blockchain.chain)-1],
		MERKLETREE:    merkleTree,
		TIMESTAMP:     time.Now(),
	}
}

// Finds the nonce value for a block*
func (blockchain *Blockchain) MineBlock(block *Block) {
	// Concatenate the previous block's hash at the end of the current block's hash
	concatenated_hashes := append(block.HASH[:], block.PREVIOUSHASH[:]...)
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

	block.NONCE = nonce        // Add the nonce to the mined block
	blockchain.AddBlock(block) // Add the block to the blockchain
}

// Adds the given block into the blockchain
func (blockchain *Blockchain) AddBlock(block *Block) {
	var previousBlock *Block = &blockchain.chain[len(blockchain.chain)-1]
	previousBlock.NEXTBLOCK = block
	block.PREVIOUSBLOCK = previousBlock
	blockchain.chain = append(blockchain.chain, *block)
}

func (blockchain Blockchain) Display() {
	for i := 0; i < len(blockchain.chain); i++ {
		fmt.Println("////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////")
		fmt.Println("// BLOCK HASH: ", blockchain.chain[i].HASH)
		fmt.Println("// PREVIOUS BLOCK HASH: ", blockchain.chain[i].PREVIOUSHASH)
		fmt.Println("// TIMESTAMP: ", blockchain.chain[i].TIMESTAMP.Format(time.RFC822))
		fmt.Println("// NONCE: ", blockchain.chain[i].NONCE)
		fmt.Println("////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////")

		DisplayMerkelTree(blockchain.chain[i])
	}
}

// Prints all the blocks in the blockchain*
func (blockchain Blockchain) DisplayBlocks() {
	for i := 0; i < len(blockchain.chain); i++ {
		fmt.Println("\n////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////")
		fmt.Println("// BLOCK HASH: ", blockchain.chain[i].HASH)
		fmt.Println("// PREVIOUS BLOCK HASH: ", blockchain.chain[i].PREVIOUSHASH)
		fmt.Println("// TIMESTAMP: ", blockchain.chain[i].TIMESTAMP.Format(time.RFC822))
		fmt.Println("// NONCE: ", blockchain.chain[i].NONCE)
		fmt.Println("////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////")
	}
}

// Prints all the transactions in a block*
func DisplayMerkelTree(block Block) {
	if len(block.MERKLETREE.TRANSACTIONS) == 0 {
		fmt.Println("\nThis block does not have any transactions.")
	} else {
		for i := 0; i < len(block.MERKLETREE.TRANSACTIONS); i++ {
			fmt.Println("")
			fmt.Println("////////////////////////////////////////////////////////////////")
			fmt.Println("// ID: ", block.MERKLETREE.TRANSACTIONS[i].ID)
			fmt.Println("// TIMESTAMP: ", block.MERKLETREE.TRANSACTIONS[i].TIMESTAMP.Format(time.RFC822))

			data, error := json.Marshal(block.MERKLETREE.TRANSACTIONS[i].DATA)
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
	for currentBlock.NEXTBLOCK != nil {
		currentHash := currentBlock.CalculateHash()
		storedHash := currentBlock.HASH

		if bytes.Equal(currentHash[:], storedHash[:]) {
			modified = true
			modifiedBlock = currentBlock
			break
		}

		currentBlock = *currentBlock.NEXTBLOCK
	}

	if modified {
		fmt.Println("The blockchain has been modified. Block ", modifiedBlock.HASH, " has been modified.")
	} else {
		fmt.Println("The blockchain has not been modified.")
	}
}

// Calculates the hash of a transaction or block*
func (transaction Transaction) CalculateHash() [32]byte {
	data, error := json.Marshal(transaction.DATA)
	if error != nil {
		fmt.Println("Error! The hash for this transaction cannot be calculated.")
		return [32]byte{}
	}
	hash := sha256.Sum256(data)
	return hash
}

func (block Block) CalculateHash() [32]byte {
	// Create a temporary Merkle tree for the transactions in the block and return the hash of the root
	temporaryMerkleTree := CreateMerkleTree(block.MERKLETREE.TRANSACTIONS)
	hash := temporaryMerkleTree.ROOT.HASH
	return hash
}

func (blockchain *Blockchain) Test(block Block) {
}

// Calculates the length of the blockchain that the given block is in
func (block Block) GetBlockchainLength() int {
	length := 0
	var currentBlock *Block = &block

	// Iterate backwards till the genesis block
	for currentBlock != nil {
		currentBlock = currentBlock.PREVIOUSBLOCK
		length++
	}

	return length
}

// Checks whether the nonce value for the given block is correct
func (blockchain Blockchain) VerifyBlock(block Block) bool {
	concatenated_hashes := append(block.HASH[:], block.PREVIOUSHASH[:]...)
	concatenated_hashes = append(concatenated_hashes[:], block.NONCE[:]...)
	hash := sha256.Sum256(concatenated_hashes[:])

	flag := true
	// Check whether the hash has the required number of leading zeroes
	for count := 0; count < blockchain.difficulty; count++ {
		if hash[count] > 0 {
			flag = false
			break
		}
	}

	return flag
}
