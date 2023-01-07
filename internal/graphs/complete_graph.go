package graphs

type CompleteGraph struct {
	InheritanceGraph
	DependencyGraph
}

func NewEntityCompleteGraph(entityDepG DependencyGraph, entityInhG InheritanceGraph) *CompleteGraph {
	completeG := &CompleteGraph{
		InheritanceGraph: entityInhG,
		DependencyGraph:  entityDepG,
	}
	completeG.CreateCompleteGraph()
	return completeG
}

func (completeG *CompleteGraph) CreateCompleteGraph() {}
