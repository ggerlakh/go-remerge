package main

import (
	"fmt"
	"go-remerge/internal/graph"
	"math/rand"
	"strconv"
)

func example1() {
	var nodes []graph.Node
	//labels := make(map[string]any)
	letters := [5]string{"A", "B", "C", "D", "E"}
	g := graph.NewGraph("undirected", []graph.Node{}, []graph.Edge{})
	for i := 0; i < 5; i++ {
		//labels["name"] = letters[i]
		n := graph.Node{Id: strconv.FormatInt(rand.Int63(), 10), Labels: map[string]any{"name": letters[i]}}
		nodes = append(nodes, n)
		g.AddNode(n)
	}
	g.AddEdge(nodes[0], nodes[4])
	g.AddEdge(nodes[2], nodes[3])
	g.AddEdge(nodes[1], nodes[3])
	fmt.Println(g)
}

func main() {
	example1()
}
