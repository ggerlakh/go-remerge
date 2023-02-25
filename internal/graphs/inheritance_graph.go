package graphs

import (
	"go-remerge/internal/parsers"
	"go-remerge/tools/hashtool"
	"strings"
)

type InheritanceGraph struct {
	Graph
	Parser            parsers.InheritanceExtractor
	AllowedExtensions []string
}

func NewEntityInheritanceGraph(entityDependencyGraph DependencyGraph, parser parsers.InheritanceExtractor, extensions []string, direction string) *InheritanceGraph {
	inhG := &InheritanceGraph{
		Graph: Graph{
			Nodes:     make(map[string]Node),
			Edges:     make(map[string]Edge),
			Direction: direction,
			Name:      "entity_inheritance",
		},
		Parser:            parser,
		AllowedExtensions: extensions,
	}
	inhG.CreateInheritanceGraph(entityDependencyGraph)
	return inhG
}

func (inhG *InheritanceGraph) CreateInheritanceGraph(entityDependencyGraph DependencyGraph) {
	for _, entityNode := range entityDependencyGraph.Nodes {
		if !strings.HasPrefix(entityNode.Labels["path"].(string), "external_dependency") {
			// creating "from" nodes
			fromNode := entityNode
			inhG.AddNode(fromNode)
			for _, parentEntityMap := range inhG.Parser.ExtractInheritance(fromNode.Labels["path"].(string), fromNode.Labels["name"].(string)) {
				toId := hashtool.Sha256(parentEntityMap["name"] + parentEntityMap["path"])
				toNode := Node{Id: toId, Labels: map[string]any{
					"name":        parentEntityMap["name"],
					"path":        parentEntityMap["path"],
					"package":     inhG.Parser.ExtractPackage(parentEntityMap["path"]),
					"isDirectory": false}}
				inhG.AddNode(toNode)
				if fromNode.Id != toNode.Id {
					inhG.AddEdge(Edge{
						From: fromNode,
						To:   toNode,
					})
				}
			}
		}
	}
}
