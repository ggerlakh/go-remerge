package graph

import (
	"fmt"
	"strings"
)

type Node struct {
	Id     string
	Labels map[string]any
}

type Graph struct {
	Type  string
	Nodes map[string]Node            // map[Node.id]Node
	Edges map[string]map[string]Node //map[FirstNode.id]map[SecondNode.id]SecondNode <-> {from: FirstNode, to: SecondNode}
}

type Edge [2]Node

func NewGraph(Type string, Nodes []Node, Edges []Edge) *Graph {
	if strings.ToLower(Type) == "undirected" || strings.ToLower(Type) == "directed" {
		g := &Graph{Type: Type, Nodes: make(map[string]Node), Edges: make(map[string]map[string]Node)} // need to init Edges map
		for _, n := range Nodes {
			g.Nodes[n.Id] = n
		}
		for _, e := range Edges {
			g.Edges[e[0].Id] = map[string]Node{e[1].Id: e[1]}
		}
		return g
	} else {
		panic(fmt.Sprintf("\"%v\" wrong graph type value, graph can be only directed or undirected", Type))
	}
}

func (g *Graph) AddNode(n Node) {
	if _, inNodes := g.Nodes[n.Id]; inNodes {
		panic(fmt.Sprintf("Duplicate error, node with such id=%s already exists.", n.Id))
	}
	g.Nodes[n.Id] = n
}

func (g *Graph) DeleteNode(n Node) {
	// deleting a node and all edges incident to it
	delete(g.Nodes, n.Id) // deleting node
	for dstNodeId, _ := range g.Edges {
		delete(g.Edges[dstNodeId], n.Id) // deleting all inbound edges
	}
	delete(g.Edges, n.Id) // deleting all outbound edges
}

func (g *Graph) AddEdge(from, to Node) {
	_, fromInMap := g.Nodes[from.Id]
	_, toInMap := g.Nodes[to.Id]
	if !fromInMap {
		panic(fmt.Sprintf("Edge creation error, node %v does not exist in graph.", from))
	} else if !toInMap {
		panic(fmt.Sprintf("Edge creation error, node %v does not exist in graph.", to))
	}
	if strings.ToLower(g.Type) == "undirected" {
		g.Edges[from.Id] = map[string]Node{to.Id: to}
		g.Edges[to.Id] = map[string]Node{from.Id: from}
	} else if strings.ToLower(g.Type) == "directed" {
		g.Edges[from.Id] = map[string]Node{to.Id: to}
	}
}

func (g *Graph) DeleteEdge(from, to Node) {
	if strings.ToLower(g.Type) == "undirected" {
		delete(g.Edges[from.Id], to.Id)
		delete(g.Edges[to.Id], from.Id)
	} else if strings.ToLower(g.Type) == "directed" {
		delete(g.Edges[from.Id], to.Id)
	}
}

/*TODO

func (g *Graph) String() string {
}

func (g *Graph) ConvertToJSON() []bytes {
}
*/
