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
	var bootstrap_data [2]string
	bootstrap_data[0] = assignment02bca.HOST
	bootstrap_data[1] = "5051"

	var port, name string
	fmt.Println("\nEnter your port number:")
	fmt.Scanln(&port)

	fmt.Println("\nEnter your name:")
	fmt.Scanln(&name)

	fmt.Println("\nCreating a node...")
	node := assignment02bca.CreateNode(name)
	node.SetPort(port)

	fmt.Println("\nDisplaying node data...")
	node.DisplayNodeData()

	fmt.Println("\nJoining the network...")
	node.JoinNetwork(bootstrap_data)

	go node.Run()

	fmt.Println("\nCreating a block with three transactions using the NODE...")
	var transactions [3]assignment01bca.Transaction
	transactions[0] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 10)
	transactions[1] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 15)
	transactions[2] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 15)
	block := assignment01bca.NewBlock(node.GetBlockchain(), transactions[:])

	node.PropogateBlock(block)

	fmt.Println("\nCreating a block with five transactions using the NODE...")
	var transactions_2 [5]assignment01bca.Transaction
	transactions_2[0] = assignment01bca.CreateTransaction("Zohaib", "Huzaifa", 10)
	transactions_2[1] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 15)
	transactions_2[2] = assignment01bca.CreateTransaction("Zohaib", "Huzaifa", 25)
	transactions_2[3] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 30)
	transactions_2[4] = assignment01bca.CreateTransaction("Zohaib", "Hussain", 50)
	block_2 := assignment01bca.NewBlock(node.GetBlockchain(), transactions_2[:])

	node.PropogateBlock(block_2)

	runtime.Goexit()
}
