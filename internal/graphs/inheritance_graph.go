package graphs

import "go-remerge/internal/parsers"

type InheritanceGraph struct {
	DependencyGraph
	Nodes             map[string]Node
	Edges             map[string]Edge
	Parser            parsers.InheritanceExtractor
	AllowedExtensions []string
	Type              string
}

func NewEntityInheritanceGraph(fileDependencyGraph DependencyGraph, parser parsers.InheritanceExtractor, extensions []string) *InheritanceGraph {
	inhG := &InheritanceGraph{
		DependencyGraph:   fileDependencyGraph,
		Nodes:             make(map[string]Node),
		Edges:             make(map[string]Edge),
		Parser:            parser,
		AllowedExtensions: extensions,
		Type:              "EntityInheritance",
	}
	inhG.CreateInheritanceGraph()
	return inhG
}

func (inhG *InheritanceGraph) CreateInheritanceGraph() {}
