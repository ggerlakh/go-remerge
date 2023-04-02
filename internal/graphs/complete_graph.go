package graphs

type CompleteGraph struct {
	Graph
}

func NewEntityCompleteGraph(entityDepG DependencyGraph, entityInhG InheritanceGraph, direction string, AnalysisProjectName string) *CompleteGraph {
	completeG := &CompleteGraph{
		Graph: Graph{
			Nodes:               make(map[string]Node),
			Edges:               make(map[string]Edge),
			Direction:           direction,
			Name:                "entity_complete",
			AnalysisProjectName: AnalysisProjectName,
		},
	}
	completeG.CreateCompleteGraph(entityDepG, entityInhG)
	return completeG
}

func (completeG *CompleteGraph) CreateCompleteGraph(entityDepG DependencyGraph, entityInhG InheritanceGraph) {
	completeG.Nodes = entityDepG.Nodes
	completeG.Edges = entityDepG.Edges
	for idInh, inhNode := range entityInhG.Nodes {
		if _, exists := completeG.Nodes[idInh]; !exists {
			completeG.AddNode(inhNode)
		}
	}
	for edgeIdInh, edgeInh := range entityInhG.Edges {
		if _, exists := completeG.Edges[edgeIdInh]; !exists {
			completeG.AddEdge(edgeInh)
		}
	}
}
