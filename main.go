package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/zohaibadnan137/assignment01bca"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Create a new blockchain...")
	blockchain := assignment01bca.CreateBlockchain(3)

	fmt.Println("Create three new transactions from Zohaib to Hussain...")

	var transactions [3]assignment01bca.Transaction
	transactions[0] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 10)
	transactions[1] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 15)
	transactions[2] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 15)

	fmt.Println("Create a new block and add the transactions to it...")
	block := assignment01bca.NewBlock(&blockchain, transactions[:])

	fmt.Println("Print the transactions added in the new block...")
	assignment01bca.DisplayMerkelTree(block)
}
