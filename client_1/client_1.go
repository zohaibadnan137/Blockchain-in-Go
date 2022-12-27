package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/zohaibadnan137/assignment02bca"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("\nCreating a bootstrapping node...")

	node := assignment02bca.CreateNode("BTSTP_NODE")
	network := assignment02bca.CreateNetwork("BLKC", &node)
	_ = network

	fmt.Println("\nDisplaying available networks...")
	assignment02bca.DisplayAvailableNetworks()

	fmt.Println("\nDisplaying node data...")
	node.DisplayNodeData()

	node.Run()
}
