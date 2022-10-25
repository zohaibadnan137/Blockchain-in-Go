package main

import "fmt"

type Block struct {
	next         *Block
	transactions interface{}
	prevHash     string
}

func main() {
	fmt.Println("hello world")
}
