package main

import (
	"math/rand"
	"time"

	"github.com/zohaibadnan137/assignment02bca"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	node := assignment02bca.CreateNode("BLKC_1")
	network := assignment02bca.CreateNetwork("BLKC", &node)
	_ = network

	go node.Bootstrapping()

	node_2 := assignment02bca.CreateNode("BLKC_2")
	node_2.JoinNetwork(assignment02bca.GetBootstrapData(0))
	node_2.DisplayNeighbours()
}
