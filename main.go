package main

import (
	"fmt"
)

type Node struct {
	id         int
	left_chan  <-chan Token
	right_chan chan Token
}

type Token struct {
	Data      string
	Recipient int
	Ttl       int
}

func (node *Node) Run() {
	token := <-node.left_chan
	fmt.Println("Node", node.id, "Token: { Recipient:", token.Recipient, "Data:", token.Data, " Ttl:", token.Ttl, "}")
	switch {
	case token.Recipient == node.id:
		fmt.Println("Token reached recipient", node.id, "Token: { Data:", token.Data, "Recipient:", token.Recipient, " Ttl:", token.Ttl, "}")
	case token.Ttl <= 0:
		fmt.Println("Token expired at node with id", node.id, "Token: { Data:", token.Data, "Recipient:", token.Recipient, " Ttl:", token.Ttl, "}")
	case token.Ttl > 0:
		token.Ttl -= 1
		node.right_chan <- token
	}
}

func initializeTokenRing(n int) []*Node {
	ring := make([]*Node, 0)

	ring = append(ring, &Node{id: 0, right_chan: make(chan Token)})

	for i := 1; i < n; i++ {
		ring = append(ring, &Node{id: i, left_chan: ring[i-1].right_chan, right_chan: make(chan Token)})
	}

	ring[0].left_chan = ring[n-1].right_chan

	return ring
}

func sendMesage(message Token) {
	for _, node := range tokenRing {
		go node.Run()
	}

	tokenRing[len(tokenRing)-1].right_chan <- message
}

var tokenRing []*Node

func main() {
	var ringSize int = 10
	var ttl int = 10
	var recipient int = 8

	tokenRing = initializeTokenRing(ringSize)

	var token = Token{"Hello world!", recipient, ttl}

	sendMesage(token)

	fmt.Scanln()

}
