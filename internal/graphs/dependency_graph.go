package graphs

import (
	"go-remerge/internal/parsers"
	"strings"
)

type DependencyGraph struct {
	FileSystemGraph
	Nodes             map[string]Node
	Edges             map[string]Edge
	Parser            parsers.DependencyExtractor
	AllowedExtensions []string
}

func NewDependencyGraph(filesystemGraph FileSystemGraph, parser parsers.DependencyExtractor, extensions []string) *DependencyGraph {
	depGraph := &DependencyGraph{
		FileSystemGraph:   filesystemGraph,
		Nodes:             make(map[string]Node),
		Edges:             make(map[string]Edge),
		Parser:            parser,
		AllowedExtensions: extensions,
	}
	depGraph.CreateDependencyGraph(filesystemGraph)
	return depGraph
}

func (depG *DependencyGraph) CheckNode(node Node) bool {
	var hasAllowedExtension bool
	for _, ext := range depG.AllowedExtensions {
		if !node.Labels["isDirectory"].(bool) && strings.HasSuffix(node.Labels["name"].(string), ext) {
			hasAllowedExtension = true
			break
		}
	}
	return hasAllowedExtension
}

func (depG *DependencyGraph) CreateDependencyGraph(filesystemGraph FileSystemGraph) {
	for _, filesystemNode := range filesystemGraph.Nodes {
		if depG.CheckNode(filesystemNode) {
			// append nodes and edges
		}
	}
}
