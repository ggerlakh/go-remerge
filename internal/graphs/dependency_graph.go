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

func NewFileDependencyGraph(filesystemGraph FileSystemGraph, parser parsers.DependencyExtractor, extensions []string, direction string, AnalysisProjectName string) *DependencyGraph {
	depGraph := &DependencyGraph{
		Graph: Graph{
			Nodes:               make(map[string]Node),
			Edges:               make(map[string]Edge),
			Direction:           direction,
			Name:                "file_dependency",
			AnalysisProjectName: AnalysisProjectName,
		},
		Parser:            parser,
		AllowedExtensions: extensions,
	}
	depGraph.CreateFileDependencyGraph(filesystemGraph)
	return depGraph
}

func NewEntityDependencyGraph(fileDependencyGraph DependencyGraph, parser parsers.DependencyExtractor, extensions []string, direction string, AnalysisProjectName string) *DependencyGraph {
	depGraph := &DependencyGraph{
		Graph: Graph{
			Nodes:               make(map[string]Node),
			Edges:               make(map[string]Edge),
			Direction:           direction,
			Name:                "entity_dependency",
			AnalysisProjectName: AnalysisProjectName,
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
				var toDependencies []string
				toId := hashtool.Sha256(filesystemGraph.GetRootRelPath(dependency))
				// extract dependencies if the local file
				if !strings.HasPrefix(dependency, "external_dependency") {
					toDependencies = depG.Parser.ExtractDependencies(dependency)
				} else {
					// do not extract entities if external dependency
					toDependencies = []string{}
				}
				toNode := Node{Id: toId, Labels: map[string]any{
					"name":         filepath.Base(dependency),
					"path":         dependency,
					"dependencies": toDependencies,
					"package":      depG.Parser.ExtractPackage(dependency),
					"isDirectory":  false}}
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
		if !strings.HasPrefix(fileDependencyNode.Labels["path"].(string), "external_dependency") {
			for _, entity := range depG.Parser.ExtractEntities(fileDependencyNode.Labels["path"].(string)) {
				// creating "from" nodes
				fromId := hashtool.Sha256(filepath.Join(fileDependencyNode.Labels["path"].(string), entity))
				fromNode := Node{Id: fromId, Labels: map[string]any{
					"name":        entity,
					"path":        fileDependencyNode.Labels["path"].(string),
					"package":     depG.Parser.ExtractPackage(fileDependencyNode.Labels["path"].(string)),
					"isDirectory": false}}
				depG.AddNode(fromNode)
				// edges creation based on file_dependency graph
				for _, dependency := range fileDependencyNode.Labels["dependencies"].([]string) {
					// external entity dependencies
					if strings.HasPrefix(dependency, "external_dependency") {
						for _, externalEntity := range depG.Parser.ExtractExternalEntities(strings.ReplaceAll(dependency, "external_dependency"+string(filepath.Separator), ""), fromNode.Labels["path"].(string), fromNode.Labels["name"].(string)) {
							// creating "to" nodes
							toId := hashtool.Sha256(externalEntity)
							toNode := Node{Id: toId, Labels: map[string]any{
								"name":        externalEntity,
								"path":        dependency,
								"package":     depG.Parser.ExtractPackage(dependency),
								"isDirectory": false}}
							depG.AddNode(toNode)
							depG.AddEdge(Edge{
								From: fromNode,
								To:   toNode,
							})
						}
					} else {
						// internal entity dependencies
						for _, depEntity := range depG.Parser.ExtractEntities(dependency) {
							toId := hashtool.Sha256(filepath.Join(dependency, depEntity))
							toNode := Node{Id: toId, Labels: map[string]any{
								"name":        depEntity,
								"path":        dependency,
								"package":     depG.Parser.ExtractPackage(dependency),
								"isDirectory": false}}
							if depG.Parser.HasEntityDependency(fromNode.Labels["name"].(string), fromNode.Labels["path"].(string), toNode.Labels["name"].(string), toNode.Labels["package"].(string)) {
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
		}
	}
}
