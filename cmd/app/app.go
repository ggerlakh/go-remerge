package main

import (
	"fmt"
	"go-remerge/internal/graph"
	"strconv"
)

func example1() {
	var nodes []graph.Node
	letters := [5]string{"A", "B", "C", "D", "E"}
	g := graph.NewGraph("undirected", []graph.Node{}, []graph.Edge{})
	for i := 0; i < 5; i++ {
		n := graph.Node{Id: strconv.FormatInt(int64(i+1), 10), Labels: map[string]any{"name": letters[i]}}
		nodes = append(nodes, n)
		g.AddNode(n)
	}
	g.AddEdge(graph.Edge{From: nodes[0], To: nodes[4]})
	g.AddEdge(graph.Edge{From: nodes[2], To: nodes[3]})
	g.AddEdge(graph.Edge{From: nodes[1], To: nodes[3]})
	fmt.Println(g)
}

func main() {
	example1()
}
