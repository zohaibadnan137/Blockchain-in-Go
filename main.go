package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/zohaibadnan137/assignment01bca"
	"github.com/zohaibadnan137/assignment02bca"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	node_1 := assignment02bca.CreateNode("BLKC_1")
	network := assignment02bca.CreateNetwork("BLKC", &node_1)
	_ = network

	node_2 := assignment02bca.CreateNode("BLKC_2")
	node_2.JoinNetwork(assignment02bca.GetBootstrapData(0))

	fmt.Println("Displaying node data...")
	node_1.DisplayNodeData()
	node_2.DisplayNodeData()

	go node_1.Run()
	go node_2.Run()

	transaction := assignment01bca.CreateTransaction("Zohaib", "Hussain", 10)
	node_2.PropogateTransaction(transaction)

	var transactions [3]assignment01bca.Transaction
	transactions[0] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 10)
	transactions[1] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 15)
	transactions[2] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 15)
	block := assignment01bca.NewBlock(node_2.GetBlockchain(), transactions[:])

	node_2.PropogateBlock(block)

	runtime.Goexit()
}
