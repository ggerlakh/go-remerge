package graphtest

import (
	"fmt"
	"go-remerge/internal/graphs"
	"strconv"
)

func UndirectedGraphCreationTest() {
	var nodes []graphs.Node
	fmt.Println("Creating undirected test graphs:")
	letters := [5]string{"A", "B", "C", "D", "E"}
	g := graphs.NewGraph("directed", []graphs.Node{}, []graphs.Edge{})
	for i := 0; i < 5; i++ {
		n := graphs.Node{Id: strconv.FormatInt(int64(i+1), 10), Labels: map[string]any{"name": letters[i]}}
		nodes = append(nodes, n)
		g.AddNode(n)
	}
	g.AddEdge(graphs.Edge{From: nodes[0], To: nodes[4]})
	g.AddEdge(graphs.Edge{From: nodes[2], To: nodes[3]})
	g.AddEdge(graphs.Edge{From: nodes[1], To: nodes[3]})
	fmt.Println(g)
}

func DirectedGraphCreationTest() {
	var nodes []graphs.Node
	fmt.Println("Creating directed test graphs:")
	letters := [5]string{"A", "B", "C", "D", "E"}
	g := graphs.NewGraph("directed", []graphs.Node{}, []graphs.Edge{})
	for i := 0; i < 5; i++ {
		n := graphs.Node{Id: strconv.FormatInt(int64(i+1), 10), Labels: map[string]any{"name": letters[i]}}
		nodes = append(nodes, n)
		g.AddNode(n)
	}
	g.AddEdge(graphs.Edge{From: nodes[0], To: nodes[4]})
	g.AddEdge(graphs.Edge{From: nodes[2], To: nodes[3]})
	g.AddEdge(graphs.Edge{From: nodes[1], To: nodes[3]})
	fmt.Println(g)
}

func FileSystemGraphCreationTest() {
	fmt.Println("Creating filesystem graphs:")
	skipDirs := []string{".git", ".idea", "neo4j"}
	skipFiles := []string{".gitignore", "go_build_go_remerge_linux"}
	fsG := graphs.NewFileSystemGraph("directed", []graphs.Node{}, []graphs.Edge{}, ".", skipDirs, skipFiles)
	fmt.Println(fsG)
}
