package graphs

type DependencyGraph struct {
	FileSystemGraph
	Language          string
	AllowedExtensions string
}

func NewDependencyGraph() *DependencyGraph {}

func (dG *DependencyGraph) CreateDependencyGraph() {}
