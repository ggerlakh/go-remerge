package graph

import (
	"fmt"
	"go-remerge/tools/jsontool"
	"log"
	"strings"
)

type Node struct {
	Id     string
	Labels map[string]any
}

type Edge struct {
	From Node
	To   Node
	Key  string
}

type Graph struct {
	Type  string
	Nodes map[string]Node // map[Node.Id]Node
	Edges map[string]Edge // map[Edge.Key]Edge
}

func NewGraph(Type string, Nodes []Node, Edges []Edge) *Graph {
	if strings.ToLower(Type) == "undirected" || strings.ToLower(Type) == "directed" {
		g := &Graph{Type: Type, Nodes: make(map[string]Node), Edges: make(map[string]Edge)} // need to init Edges map
		for _, n := range Nodes {
			g.Nodes[n.Id] = n
		}
		for _, e := range Edges {
			e.Key = e.From.Id + "->" + e.To.Id
			g.Edges[e.Key] = e
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
	for key, _ := range g.Edges {
		if strings.Contains(key, n.Id) {
			delete(g.Edges, key) // deleting incident edge
		}
	}
}

func (g *Graph) AddEdge(e Edge) {
	_, fromInMap := g.Nodes[e.From.Id]
	_, toInMap := g.Nodes[e.To.Id]
	if !fromInMap {
		panic(fmt.Sprintf("Edge creation error, node %v does not exist in graph.", e.From))
	} else if !toInMap {
		panic(fmt.Sprintf("Edge creation error, node %v does not exist in graph.", e.To))
	}
	if strings.ToLower(g.Type) == "undirected" {
		e.Key = e.From.Id + "->" + e.To.Id
		g.Edges[e.Key] = e
		revEdge := Edge{From: e.To, To: e.From}
		revEdge.Key = revEdge.From.Id + "->" + revEdge.To.Id
		g.Edges[revEdge.Key] = revEdge
	} else if strings.ToLower(g.Type) == "directed" {
		e.Key = e.From.Id + "->" + e.To.Id
		g.Edges[e.Key] = e
	}
}

func (g *Graph) DeleteEdge(e Edge) {
	if strings.ToLower(g.Type) == "undirected" {
		delete(g.Edges, e.From.Id+"->"+e.To.Id)
		delete(g.Edges, e.To.Id+"->"+e.From.Id)
	} else if strings.ToLower(g.Type) == "directed" {
		delete(g.Edges, e.From.Id+"->"+e.To.Id)
	}
}

func (g *Graph) GetNodes() []Node {
	var nodes []Node
	for _, value := range g.Nodes {
		nodes = append(nodes, value)
	}
	return nodes
}

func (g *Graph) GetEdges() []Edge {
	var edges []Edge
	for _, value := range g.Edges {
		edges = append(edges, value)
	}
	return edges
}

func (g *Graph) PrettyJson() string {
	nodes := g.GetNodes()
	edges := g.GetEdges()
	prettyGraph := map[string]any{"type": g.Type, "nodes": nodes, "edges": edges}
	b, err := jsontool.ExtenedMarshal(prettyGraph, "", "\t", false)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func (g *Graph) String() string {
	return g.PrettyJson()
}
