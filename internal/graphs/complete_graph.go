package graphs

type CompleteGraph struct {
	Graph
}

func NewEntityCompleteGraph(entityDepG DependencyGraph, entityInhG InheritanceGraph, direction string) *CompleteGraph {
	completeG := &CompleteGraph{
		Graph: Graph{
			Nodes:     make(map[string]Node),
			Edges:     make(map[string]Edge),
			Direction: direction,
			Name:      "file_dependency",
		},
	}
	completeG.CreateCompleteGraph()
	return completeG
}

func (completeG *CompleteGraph) CreateCompleteGraph() {}
