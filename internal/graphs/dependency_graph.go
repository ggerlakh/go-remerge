package graphs

import (
	"go-remerge/internal/parsers"
	"go-remerge/tools/hashtool"
	"path/filepath"
	"strings"
)

type DependencyGraph struct {
	Graph
	Parser            parsers.DependencyExtractor
	AllowedExtensions []string
}

func NewFileDependencyGraph(filesystemGraph FileSystemGraph, parser parsers.DependencyExtractor, extensions []string, direction string) *DependencyGraph {
	depGraph := &DependencyGraph{
		Graph: Graph{
			Nodes:     make(map[string]Node),
			Edges:     make(map[string]Edge),
			Direction: direction,
			Name:      "file_dependency",
		},
		Parser:            parser,
		AllowedExtensions: extensions,
	}
	depGraph.CreateFileDependencyGraph(filesystemGraph)
	return depGraph
}

func NewEntityDependencyGraph(fileDependencyGraph DependencyGraph, parser parsers.DependencyExtractor, extensions []string, direction string) *DependencyGraph {
	depGraph := &DependencyGraph{
		Graph: Graph{
			Nodes:     make(map[string]Node),
			Edges:     make(map[string]Edge),
			Direction: direction,
			Name:      "entity_dependency",
		},
		Parser:            parser,
		AllowedExtensions: extensions,
	}
	depGraph.CreateEntityDependencyGraph(fileDependencyGraph)
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

func (depG *DependencyGraph) CreateFileDependencyGraph(filesystemGraph FileSystemGraph) {
	for _, filesystemNode := range filesystemGraph.Nodes {
		if depG.CheckNode(filesystemNode) {
			// append nodes and edges for file dependency
			// adding "from" node
			fileDependencyNode := filesystemNode
			fileDependencyNode.Labels["package"] = depG.Parser.ExtractPackage(fileDependencyNode.Labels["path"].(string))
			fileDependencyNode.Labels["dependencies"] = depG.Parser.ExtractDependencies(fileDependencyNode.Labels["path"].(string))
			depG.AddNode(fileDependencyNode)
			// adding "to" nodes (import dependencies)
			for _, dependency := range fileDependencyNode.Labels["dependencies"].([]string) {
				toId := hashtool.Sha256(dependency)
				toNode := Node{Id: toId, Labels: map[string]any{
					"name":        filepath.Base(dependency),
					"path":        dependency,
					"package":     depG.Parser.ExtractPackage(dependency),
					"isDirectory": false}}
				depG.AddNode(toNode)
				depG.AddEdge(Edge{
					From: fileDependencyNode,
					To:   toNode,
				})
			}
		}
	}
}

func (depG *DependencyGraph) CreateEntityDependencyGraph(fileDependencyGraph DependencyGraph) {
	for _, fileDependencyNode := range fileDependencyGraph.Nodes {
		for _, entity := range depG.Parser.ExtractEntities(fileDependencyNode.Labels["path"].(string)) {
			// creating "from" nodes
			fromId := hashtool.Sha256(entity)
			fromNode := Node{Id: fromId, Labels: map[string]any{
				"name":        entity,
				"path":        fileDependencyNode.Labels["path"].(string),
				"package":     depG.Parser.ExtractPackage(fileDependencyNode.Labels["path"].(string)),
				"isDirectory": false}}
			depG.AddNode(fromNode)
			// creating "to" nodes
			for _, dependency := range fileDependencyNode.Labels["dependencies"].([]string) {
				for _, depEntity := range depG.Parser.ExtractEntities(dependency) {
					toId := hashtool.Sha256(depEntity)
					toNode := Node{Id: toId, Labels: map[string]any{
						"name":        depEntity,
						"path":        dependency,
						"package":     depG.Parser.ExtractPackage(dependency),
						"isDirectory": false}}
					depG.AddNode(toNode)
					depG.AddEdge(Edge{
						From: fromNode,
						To:   toNode,
					})
				}
			}
		}
	}
}
