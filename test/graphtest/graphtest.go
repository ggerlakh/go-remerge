package graphtest

import (
	"context"
	"fmt"
	"go-remerge/internal/graphs"
	"log"
	"strconv"
)

func UndirectedGraphCreationTest() {
	var nodes []graphs.Node
	fmt.Println("Creating undirected test graphs:")
	letters := [5]string{"A", "B", "C", "D", "E"}
	g := graphs.NewGraph("undirected", "undirected_test", []graphs.Node{}, []graphs.Edge{})
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
	g := graphs.NewGraph("directed", "directed_test", []graphs.Node{}, []graphs.Edge{})
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
	fsG := graphs.NewFileSystemGraph("directed", "filesystem", []graphs.Node{}, []graphs.Edge{}, ".", skipDirs, skipFiles)
	fmt.Println(fsG)
}

func Neo4jLoadingGraphTest() {
	ctx := context.Background()
	skipDirs := []string{".git", ".idea", "neo4jdb"}
	skipFiles := []string{".gitignore", "go_build_go_remerge_linux", "token", "log.json"}
	fsG := graphs.NewFileSystemGraph("directed", "filesystem_neo4j", []graphs.Node{}, []graphs.Edge{}, ".", skipDirs, skipFiles)
	err := fsG.LoadNeo4jGraph(ctx, "neo4j://localhost:7687", "neo4j", "neo4jdevops")
	if err != nil {
		log.Fatal("Neo4j loading graph failed: %v\n", err)
	}
}

func GetArangoGraphTest() {
	skipDirs := []string{".git", ".idea", "neo4jdb"}
	skipFiles := []string{".gitignore", "go_build_go_remerge_linux", "token", "log.json"}
	fsG := graphs.NewFileSystemGraph("directed", "filesystem_arango", []graphs.Node{}, []graphs.Edge{}, ".", skipDirs, skipFiles)
	fmt.Println(fsG.ToArango())
}

func ArangoLoadingGraphTest() {
	ctx := context.Background()
	skipDirs := []string{".git", ".idea", "neo4jdb"}
	skipFiles := []string{".gitignore", "go_build_go_remerge_linux", "token", "log.json"}
	fsG := graphs.NewFileSystemGraph("directed", "filesystem_arango", []graphs.Node{}, []graphs.Edge{}, ".", skipDirs, skipFiles)
	endpoints := []string{"http://localhost:8529"}
	fsG.LoadArangoGraph(ctx, endpoints, "root", "password", "test")
}

func FileDependencyCreationTest() {}

func EntityDependencyCreationTest() {}

func EntityInheritanceCreationTest() {}

func EntityCompleteCreationTest() {}
