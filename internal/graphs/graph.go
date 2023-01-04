package graphs

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

func (n *Node) ToCypher() string {
	var props []string
	props = append(props, fmt.Sprintf("Id: '%v'", n.Id))
	for k, v := range n.Labels {
		if strings.Contains(fmt.Sprintf("%v", v), `\t`) {
			v = strings.Replace(fmt.Sprintf("%v", v), `\t`, `\\t`, -1)
		}
		props = append(props, fmt.Sprintf("%s: '%v'", k, v))
	}
	return "{" + strings.Join(props, ", ") + "}"
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
		g.SetNodes(Nodes)
		g.SetEdges(Edges)
		return g
	} else {
		panic(fmt.Sprintf("\"%v\" wrong graphs type value, graphs can be only directed or undirected", Type))
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
		panic(fmt.Sprintf("Edge creation error, node %v does not exist in graphs.", e.From))
	} else if !toInMap {
		panic(fmt.Sprintf("Edge creation error, node %v does not exist in graphs.", e.To))
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

func (g *Graph) SetNodes(Nodes []Node) {
	for _, n := range Nodes {
		g.Nodes[n.Id] = n
	}
}

func (g *Graph) GetEdges() []Edge {
	var edges []Edge
	for _, value := range g.Edges {
		edges = append(edges, value)
	}
	return edges
}

func (g *Graph) SetEdges(Edges []Edge) {
	for _, e := range Edges {
		e.Key = e.From.Id + "->" + e.To.Id
		g.Edges[e.Key] = e
	}
}

func (g *Graph) GetPrettyGraph() map[string]any {
	nodes := g.GetNodes()
	edges := g.GetEdges()
	return map[string]any{"type": g.Type, "nodes": nodes, "edges": edges}
}

func (g *Graph) GetPrettyJson() string {
	prettyGraph := g.GetPrettyGraph()
	b, err := jsontool.ExtendedMarshal(prettyGraph, "", "\t", false)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func (g *Graph) GetCypher() []string {
	var cypherArr []string
	// node creation in cypher
	for _, node := range g.Nodes {
		cypherArr = append(cypherArr, fmt.Sprintf("CREATE (n: Node %s);", node.ToCypher()))
	}
	//edge creation in cypher
	for _, edge := range g.Edges {
		cypherArr = append(cypherArr,
			fmt.Sprintf("MATCH (from: Node {Id: '%s'}),  (to: Node {Id: '%s'}) MERGE (from)-[r: CONNECTED_WITH]->(to);",
				edge.From.Id, edge.To.Id))
	}
	return cypherArr
}

func (g *Graph) String() string {
	return g.GetPrettyJson()
}
