package analyzer

import (
	"fmt"
	"go-remerge/internal/config"
	"go-remerge/internal/graphs"
	"go-remerge/internal/parsers"
)

type Analyzer struct {
	Conf     config.Config
	GraphMap map[string]graphs.Exporter
}

func (a *Analyzer) Start() {
	for _, gConf := range a.Conf.Graphs {
		if gConf.Graph == "filesystem" {
			a.CreateFilesystemGraphIfNotCreated(gConf.Type)
			fmt.Println(a.GraphMap["filesystem"])
			//fmt.Println("break")
			//break
			a.Export(a.GraphMap["filesystem"])
		} else {
			for _, lang := range a.Conf.Languages {
				var parser parsers.CompleteGraphExtractor
				switch lang {
				case "python":
					parser = &parsers.PythonParser{}
				case "golang", "go":
					parser = &parsers.GoParser{}
				case "kotlin":
					parser = &parsers.KotlinParser{}
				case "swift":
					parser = &parsers.SwiftParser{}
				}
				switch gConf.Graph {
				case "file_dependency":
					fmt.Printf("TODO %v\n", gConf.Graph)
					a.CreateFileDependencyGraphIfNotCreated(gConf.Type, parser)
					a.Export(a.GraphMap[gConf.Graph])
				case "entity_dependency":
					fmt.Printf("TODO %v\n", gConf.Graph)
					a.CreateEntityDependencyGraphIfNotCreated(gConf.Type, parser)
					a.Export(a.GraphMap[gConf.Graph])
				case "entity_inheritance":
					fmt.Printf("TODO %v\n", gConf.Graph)
					a.CreateEntityInheritanceGraphIfNotCreated(gConf.Type, parser)
					a.Export(a.GraphMap[gConf.Graph])
				case "entity_complete":
					fmt.Printf("TODO %v\n", gConf.Graph)
					a.CreateEntityCompleteGraphIfNotCreated(gConf.Type, parser)
					a.Export(a.GraphMap[gConf.Graph])
				}
			}
		}
	}
}

func (a *Analyzer) CreateFilesystemGraphIfNotCreated(Type string) {
	if _, isFilesystemGraphCreated := a.GraphMap["filesystem"]; !isFilesystemGraphCreated { // create graph if not created
		a.GraphMap["filesystem"] = graphs.NewFileSystemGraph(Type, "filesystem", []graphs.Node{}, []graphs.Edge{},
			a.Conf.SourceDirectory, a.Conf.IgnoreDirectories, a.Conf.IgnoreFiles)
	} else if a.GraphMap["filesystem"].(*graphs.FileSystemGraph).Type != Type { // if type not equal type from config, create new graph with correct type
		a.GraphMap["filesystem"] = graphs.NewFileSystemGraph(Type, "filesystem", []graphs.Node{}, []graphs.Edge{},
			a.Conf.SourceDirectory, a.Conf.IgnoreDirectories, a.Conf.IgnoreFiles)
	}
}

func (a *Analyzer) CreateFileDependencyGraphIfNotCreated(Type string, parser parsers.DependencyExtractor) {
	a.CreateFilesystemGraphIfNotCreated(Type)
	if _, isFileDependencyGraphCreated := a.GraphMap["file_dependency"]; !isFileDependencyGraphCreated {
		a.GraphMap["file_dependency"] = graphs.NewFileDependencyGraph(*a.GraphMap["filesystem"].(*graphs.FileSystemGraph),
			parser, a.Conf.Extensions)
	} else if a.GraphMap["file_dependency"].(*graphs.DependencyGraph).Type != Type {
		a.GraphMap["file_dependency"] = graphs.NewFileDependencyGraph(*a.GraphMap["filesystem"].(*graphs.FileSystemGraph),
			parser, a.Conf.Extensions)
	}
}

func (a *Analyzer) CreateEntityDependencyGraphIfNotCreated(Type string, parser parsers.DependencyExtractor) {
	a.CreateFilesystemGraphIfNotCreated(Type)
	if _, isEntityDependencyGraphCreated := a.GraphMap["entity_dependency"]; !isEntityDependencyGraphCreated {
		a.GraphMap["entity_dependency"] = graphs.NewEntityDependencyGraph(*a.GraphMap["file_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions)
	} else if a.GraphMap["entity_dependency"].(*graphs.DependencyGraph).Type != Type {
		a.GraphMap["entity_dependency"] = graphs.NewEntityDependencyGraph(*a.GraphMap["file_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions)
	}
}

func (a *Analyzer) CreateEntityInheritanceGraphIfNotCreated(Type string, parser parsers.InheritanceExtractor) {
	a.CreateFilesystemGraphIfNotCreated(Type)
	if _, isEntityInheritanceGraphCreated := a.GraphMap["entity_inheritance"]; !isEntityInheritanceGraphCreated {
		a.GraphMap["entity_inheritance"] = graphs.NewEntityInheritanceGraph(*a.GraphMap["file_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions)
	} else if a.GraphMap["entity_inheritance"].(*graphs.DependencyGraph).Type != Type {
		a.GraphMap["entity_inheritance"] = graphs.NewEntityInheritanceGraph(*a.GraphMap["file_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions)
	}
}

func (a *Analyzer) CreateEntityCompleteGraphIfNotCreated(Type string, parser parsers.CompleteGraphExtractor) {
	a.CreateEntityInheritanceGraphIfNotCreated(Type, parser)
	a.CreateEntityDependencyGraphIfNotCreated(Type, parser)
	if _, isEntityCompleteGraphCreated := a.GraphMap["entity_complete"]; !isEntityCompleteGraphCreated {
		a.GraphMap["entity_complete"] = graphs.NewEntityCompleteGraph(*a.GraphMap["entity_dependency"].(*graphs.DependencyGraph),
			*a.GraphMap["entity_inheritance"].(*graphs.InheritanceGraph))
	} else if a.GraphMap["entity_complete"].(*graphs.CompleteGraph).Type != Type {
		a.GraphMap["entity_complete"] = graphs.NewEntityCompleteGraph(*a.GraphMap["entity_dependency"].(*graphs.DependencyGraph),
			*a.GraphMap["entity_inheritance"].(*graphs.InheritanceGraph))
	}
}

func (a *Analyzer) Export(g graphs.Exporter) {
} //TODO: export
