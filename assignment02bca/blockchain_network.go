package assignment02bca

import (
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strconv"
	"sync"

	"github.com/zohaibadnan137/assignment01bca"
)

// ******** ******** CONSTANTS ******** ******** //

const (
	HOST       = "localhost"
	TYPE       = "tcp"
	DIFFICULTY = 2
)

// ******** ******** GLOBAL VARIABLES ******** ******** //
var availableNetworks []*Network

var mutex sync.Mutex
var nextAssignableNetworkId int = 0 // Stores the next ID that can be assigned to a new network
var nextAssignableNodeId int = 0    // Stores the next ID that can be assigned to a new node
var nextAssignablePort int = 5050   // Stores the next port number that can be assigned to a new node

// ******** ******** STRUCTURES ******** ******** //

type LightNode struct {
	IP   string
	PORT string
}

type Bootstrapper struct {
	ip   string
	port string // The bootstrap node has a separate port number to accept incoming bootstrap requests

	nodes []LightNode // The bootstrap node maintains a list of all connected nodes in the network
}

type Node struct {
	id   int
	name string
	ip   string
	port string

	blockchain   assignment01bca.Blockchain // Each node has its own copy of the blockchain
	transactions []assignment01bca.Transaction

	neighbours   []LightNode
	bootstrapper Bootstrapper // Only the bootstrap node will store any bootstrapping data
}

type Network struct {
	id            int
	name          string
	bootstrapNode *Node
}

// ******** ******** UTILITIES ******** ******** //

type UniqueRandInt struct {
	generated map[int]bool
}

// Used to generate unique random numbers without duplication
func (u *UniqueRandInt) GenerateUniqueRandInt(n int) int {
	for {
		i := rand.Intn(n)
		if !u.generated[i] {
			u.generated[i] = true
			return i
		}
	}
}

// ******** ******** PRIVATE FUNCTIONS ******** ******** //

func getAssignableNetworkId() int {
	mutex.Lock()
	defer mutex.Unlock()

	assignableNetworkId := nextAssignableNetworkId
	nextAssignableNetworkId++

	return assignableNetworkId
}

func getAssignableNodeId() int {
	mutex.Lock()
	defer mutex.Unlock()

	assignableNodeId := nextAssignableNodeId
	nextAssignableNodeId++

	return assignableNodeId
}

func getAssignablePort() string {
	mutex.Lock()
	defer mutex.Unlock()

	assignablePort := nextAssignablePort
	nextAssignablePort++

	return strconv.Itoa(assignablePort)
}

func getNetwork(id int) *Network {
	for i := 0; i < len(availableNetworks); i++ {
		if availableNetworks[i].id == id {
			return availableNetworks[i]
		}
	}
	return nil
}

// Handle incoming bootstrapping requests
func (node *Node) Bootstrapping() {
	bootstrapper, err := net.Listen(TYPE, node.bootstrapper.ip+":"+node.bootstrapper.port)

	if err != nil {
		// TODO
		fmt.Println("ERROR")
	}

	for {
		conn, err := bootstrapper.Accept()

		if err != nil {
			// TODO
			fmt.Println("ERROR")
		}

		go node.handleBootstrapping(conn)
	}
}
func (node *Node) handleBootstrapping(conn net.Conn) {
	defer conn.Close()

	// Receive the new connecting node's IP address and port number
	decoder := gob.NewDecoder(conn)
	var ip, port string
	decoder.Decode(&ip)
	decoder.Decode(&port)

	numNodes := len(node.bootstrapper.nodes)

	// Calculate the number of nodes that should be sent to the new connecting node
	var numNeighbours int = int(math.Ceil(float64(numNodes)/2)) + int(math.Log(float64(numNodes)))

	numSelectedNeighbours := 0
	selectedNeighbours := make([]LightNode, numNeighbours) // A list of selected nodes that will be the new connecting node's neighbours

	u := UniqueRandInt{make(map[int]bool)} // Initialize the unique random integer generator

	for i := 0; numSelectedNeighbours < numNeighbours; i++ {
		// Generate a unique integer between zero and the total number of nodes in the network
		index := u.GenerateUniqueRandInt(numNodes)

		// Add the selected node into the list of neighbours
		neighbour := LightNode{node.bootstrapper.nodes[index].IP, node.bootstrapper.nodes[index].PORT}
		selectedNeighbours[i] = neighbour

		numSelectedNeighbours++
	}

	// Send the list of neighbours along with the number to the new connecting node
	encoder := gob.NewEncoder(conn)
	encoder.Encode(numNeighbours)
	encoder.Encode(selectedNeighbours)

	// Add the new connecting node to the list of all nodes in the network
	newNode := LightNode{ip, port}
	node.bootstrapper.nodes = append(node.bootstrapper.nodes, newNode)
}

// ******** ******** PUBLIC FUNCTIONS ******** ******** //

func CreateNetwork(name string, bootstrapNode *Node) Network {
	network := Network{
		getAssignableNetworkId(),
		name,
		bootstrapNode,
	}

	availableNetworks = append(availableNetworks, &network)

	// Assign bootstrap data to the boostrap node
	bootstrapper := Bootstrapper{
		ip:   HOST,
		port: getAssignablePort(),
	}
	bootstrapNode.bootstrapper = bootstrapper

	// Add the bootstrap node to the list of all nodes in the network
	currNode := LightNode{bootstrapNode.ip, bootstrapNode.port}
	bootstrapNode.bootstrapper.nodes = append(bootstrapNode.bootstrapper.nodes, currNode)

	return network
}

// Displays all available networks
func DisplayAvailableNetworks() {
	for i := 0; i < len(availableNetworks); i++ {
		fmt.Println("////////////////////////////////")
		fmt.Println("// ID: ", availableNetworks[i].id)
		fmt.Println("// NAME: ", availableNetworks[i].name)
		fmt.Println("////////////////////////////////")
	}
}

// Returns the IP address and port number of the bootstrap node for the given network
func GetBootstrapData(id int) [2]string {
	var bootstrapData [2]string
	network := getNetwork(id)

	if network != nil {
		bootstrapData[0] = network.bootstrapNode.bootstrapper.ip
		bootstrapData[1] = network.bootstrapNode.bootstrapper.port
	} else {
		bootstrapData[0] = "NULL"
		bootstrapData[1] = "NULL"
	}

	return bootstrapData
}

// Contacts the bootstrap node and returns a list of received neighbours
func (node *Node) JoinNetwork(bootstrapData [2]string) {
	conn, err := net.Dial(TYPE, bootstrapData[0]+":"+bootstrapData[1])

	if err != nil {
		// TODO
		fmt.Println("ERROR")
	}

	defer conn.Close()

	// Send the connecting node's IP address and port number to the bootstrap node
	encoder := gob.NewEncoder(conn)
	encoder.Encode(node.ip)
	encoder.Encode(node.port)

	// Receive the list of neighbours along with the number from the bootstrap node
	decoder := gob.NewDecoder(conn)

	var numNeighbours int
	decoder.Decode(&numNeighbours)

	neighbours := make([]LightNode, numNeighbours)
	decoder.Decode(&neighbours)

	node.neighbours = neighbours
}

func CreateNode(name string) Node {
	return Node{
		id:         getAssignableNodeId(),
		name:       name,
		ip:         HOST,
		port:       getAssignablePort(),
		blockchain: assignment01bca.CreateBlockchain(DIFFICULTY),
	}
}

func (node *Node) UpdateNeighbours() {

}

func (node *Node) Run() {
}

func (node Node) PropogateTransaction(transaction assignment01bca.Transaction) {

}

func (node Node) PropogateBlock(block assignment01bca.Block) {

}

func (node Node) ReceiveTransaction() {

}

func (node Node) ReceiveBlock() {

}

func (node Node) DisplayNeighbours() {
	for i := 0; i < len(node.neighbours); i++ {
		fmt.Println("////////////////////////////////")
		fmt.Println("// IP: ", node.neighbours[i].IP)
		fmt.Println("// PORT: ", node.neighbours[i].PORT)
		fmt.Println("////////////////////////////////")
	}
}
