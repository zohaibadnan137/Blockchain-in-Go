package assignment01bca

// Structs to be used
type Block struct {
	// Hash of the current block
	Hash []byte
	// Data i.e our transaction(s)
	Data []byte
	// Hash of the previous block
	PrevHash []byte
	// Random number used to solve the crypto puzzle
	Nounce int
}

// Represents the Blockchain
type Blockchain struct {
	Chain []*Block
}

// ******** ******** PUBLIC FUNCTIONS ******** ******** //

// Creates a new block
func NewBlock() {

}

// Finds the nonce value for a block
func MineBlock() {

}

// Prints all the blocks in the blockchain
func DisplayBlocks() {

}

// Prints all the transactions in a block
func DisplayMerkelTree() {

}

// Changes one or multiple transactions in a block
func ChangeBlock() {

}

// Verifies whether any changes were made in the blockchain
func VerifyChain() {

}

// Calculates the hash of a transaction or block
func CalculateHash() {

}
