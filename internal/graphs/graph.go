package graphs

import (
	"context"
	"fmt"
	"go-remerge/internal/arango"
	"go-remerge/internal/neo4j"
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
	Direction string
	Name      string
	Nodes     map[string]Node // map[Node.Id]Node
	Edges     map[string]Edge // map[Edge.Key]Edge
}

func NewGraph(Direction, Name string, Nodes []Node, Edges []Edge) *Graph {
	if strings.ToLower(Direction) == "undirected" || strings.ToLower(Direction) == "directed" {
		g := &Graph{Direction: strings.ToLower(Direction), Name: Name, Nodes: make(map[string]Node), Edges: make(map[string]Edge)} // need to init Edges map
		g.SetNodes(Nodes)
		g.SetEdges(Edges)
		return g
	} else {
		panic(fmt.Sprintf("\"%v\" wrong graphs type value, graphs can be only directed or undirected", Direction))
	}
}

func (g *Graph) AddNode(n Node) {
	/*if _, inNodes := g.Nodes[n.Id]; inNodes {
		panic(fmt.Sprintf("Duplicate error, node with such id=%s already exists. Node: %v", n.Id, n))
	}*/
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
		panic(fmt.Sprintf("Edge creation error, node %v does not exist in graphs.\n%v", e.From, e))
	} else if !toInMap {
		panic(fmt.Sprintf("Edge creation error, node %v does not exist in graphs.\n%v", e.To, e))
	}
	if g.Direction == "undirected" {
		e.Key = e.From.Id + "->" + e.To.Id
		g.Edges[e.Key] = e
		revEdge := Edge{From: e.To, To: e.From}
		revEdge.Key = revEdge.From.Id + "->" + revEdge.To.Id
		g.Edges[revEdge.Key] = revEdge
	} else if g.Direction == "directed" {
		e.Key = e.From.Id + "->" + e.To.Id
		g.Edges[e.Key] = e
	}
}

func (g *Graph) DeleteEdge(e Edge) {
	if g.Direction == "undirected" {
		delete(g.Edges, e.From.Id+"->"+e.To.Id)
		delete(g.Edges, e.To.Id+"->"+e.From.Id)
	} else if g.Direction == "directed" {
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
	return map[string]any{"type": g.Direction, "nodes": nodes, "edges": edges}
}

func (g *Graph) GetPrettyJson() string {
	prettyGraph := g.GetPrettyGraph()
	b, err := jsontool.ExtendedMarshal(prettyGraph, "", "\t", false)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func (g *Graph) String() string {
	return g.GetPrettyJson()
}

func (g *Graph) GetArangoNodes() []map[string]any {
	var arangoNodes []map[string]any
	for _, node := range g.Nodes {
		arangoNode := node.Labels
		arangoNode["_key"] = node.Id
		arangoNode["_id"] = fmt.Sprintf("%s_nodes/%s", g.Name, node.Id)
		arangoNodes = append(arangoNodes, arangoNode)
	}
	return arangoNodes
}

func (g *Graph) GetArangoEdges() []map[string]any {
	var arangoEdges []map[string]any
	for _, edge := range g.Edges {
		arangoEdge := make(map[string]any)
		//arangoEdge["_key"] = edge.Key
		arangoEdge["_id"] = fmt.Sprintf("%s_edges/%s", g.Name, edge.Key)
		arangoEdge["_from"] = fmt.Sprintf("%s_nodes/%s", g.Name, edge.From.Id)
		arangoEdge["_to"] = fmt.Sprintf("%s_nodes/%s", g.Name, edge.To.Id)
		arangoEdges = append(arangoEdges, arangoEdge)
	}
	return arangoEdges
}

func (g *Graph) ToArango() string {
	jsonArangoNodes, err := jsontool.ExtendedMarshal(g.GetArangoNodes(), "", "\t", false)
	if err != nil {
		log.Fatal(err)
	}
	jsonArangoEdges, err := jsontool.ExtendedMarshal(g.GetArangoEdges(), "", "\t", false)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonArangoNodes) + "\n" + string(jsonArangoEdges)
}

func (g *Graph) GetCypher() []string {
	var cypherArr []string
	// node creation in cypher
	for _, node := range g.Nodes {
		//cypherArr = append(cypherArr, fmt.Sprintf("CREATE (n: Node %s);", node.ToCypher()))
		cypherArr = append(cypherArr, fmt.Sprintf("CREATE (n: %s %s);", g.Name, node.ToCypher()))
	}
	//edge creation in cypher
	for _, edge := range g.Edges {
		cypherArr = append(cypherArr,
			fmt.Sprintf("MATCH (from: %s {Id: '%s'}),  (to: %s {Id: '%s'}) MERGE (from)-[r: CONNECTED_WITH]->(to);",
				g.Name, edge.From.Id, g.Name, edge.To.Id))
	}
	return cypherArr
}

func (g *Graph) LoadNeo4jGraph(ctx context.Context, uri, username, password string) error {
	cypherQueries := g.GetCypher()
	err := neo4j.ExecCypher(ctx, uri, username, password, cypherQueries)
	return err
}

func (g *Graph) LoadArangoGraph(ctx context.Context, endpoints []string, username, password, dbName string) {
	nodes := g.GetArangoNodes()
	edges := g.GetArangoEdges()
	arango.LoadGraph(ctx, endpoints, username, password, dbName, g.Name, nodes, edges)
}

type Exporter interface {
	GetPrettyJson() string
	ToArango() string
	LoadNeo4jGraph(ctx context.Context, uri, username, password string) error
	LoadArangoGraph(ctx context.Context, endpoints []string, username, password, dbName string)
}
