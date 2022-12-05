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

	fmt.Println("\nCreating two nodes...")

	node_1 := assignment02bca.CreateNode("NODE_1")
	network := assignment02bca.CreateNetwork("BLKC", &node_1)
	_ = network

	node_2 := assignment02bca.CreateNode("NODE_2")
	node_2.JoinNetwork(assignment02bca.GetBootstrapData(0))

	fmt.Println("\nDisplaying node data...")
	node_1.DisplayNodeData()
	node_2.DisplayNodeData()

	go node_1.Run()
	go node_2.Run()

	fmt.Println("\nCreating a block with three transactions using NODE_2...")
	var transactions [3]assignment01bca.Transaction
	transactions[0] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 10)
	transactions[1] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 15)
	transactions[2] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 15)
	block := assignment01bca.NewBlock(node_2.GetBlockchain(), transactions[:])

	go node_2.PropogateBlock(block)

	fmt.Println("\nCreating two more nodes...")
	node_3 := assignment02bca.CreateNode("NODE_3")
	node_3.JoinNetwork(assignment02bca.GetBootstrapData(0))

	node_4 := assignment02bca.CreateNode("NODE_4")
	node_4.JoinNetwork(assignment02bca.GetBootstrapData(0))

	fmt.Println("\nDisplaying node data...")
	node_3.DisplayNodeData()
	node_4.DisplayNodeData()

	fmt.Println("\nCreating a block with five transactions using NODE_3...")
	var transactions_2 [5]assignment01bca.Transaction
	transactions_2[0] = assignment01bca.CreateTransaction("Zohaib", "Huzaifa", 10)
	transactions_2[1] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 15)
	transactions_2[2] = assignment01bca.CreateTransaction("Zohaib", "Huzaifa", 25)
	transactions_2[3] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 30)
	transactions_2[4] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 50)
	block_2 := assignment01bca.NewBlock(node_2.GetBlockchain(), transactions_2[:])

	go node_3.PropogateBlock(block_2)

	runtime.Goexit()
}
