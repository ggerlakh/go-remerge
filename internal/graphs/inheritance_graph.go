package graphs

import "go-remerge/internal/parsers"

type InheritanceGraph struct {
	Graph
	Nodes             map[string]Node
	Edges             map[string]Edge
	Parser            parsers.InheritanceExtractor
	AllowedExtensions []string
}

func NewEntityInheritanceGraph(fileDependencyGraph DependencyGraph, parser parsers.InheritanceExtractor, extensions []string, direction string) *InheritanceGraph {
	inhG := &InheritanceGraph{
		Graph: Graph{
			Nodes:     make(map[string]Node),
			Edges:     make(map[string]Edge),
			Direction: direction,
			Name:      "file_dependency",
		},
		AllowedExtensions: extensions,
	}
	inhG.CreateInheritanceGraph()
	return inhG
}

func (inhG *InheritanceGraph) CreateInheritanceGraph() {}
