package analyzer

import (
	"fmt"
	"go-remerge/internal/config"
	"go-remerge/internal/graphs"
	"go-remerge/internal/parsers"
)

type Analyzer struct {
	Conf config.Config
}

func (a *Analyzer) Start() {
	var isFilesystemGraphCreated bool
	var filesystemGraph *graphs.FileSystemGraph
	for _, gConf := range a.Conf.Graphs {
		if gConf.Name == "filesystem" {
			fsG := graphs.NewFileSystemGraph(gConf.Type, gConf.Name, []graphs.Node{}, []graphs.Edge{}, a.Conf.SourceDirectory,
				a.Conf.IgnoreDirectories, a.Conf.IgnoreFiles)
			fmt.Println(fsG)
			isFilesystemGraphCreated = true
			filesystemGraph = fsG
			//fmt.Println("break")
			//break
			a.Export(fsG)
		} else {
			if !isFilesystemGraphCreated {
				filesystemGraph = graphs.NewFileSystemGraph("directed", "filesystem", []graphs.Node{}, []graphs.Edge{},
					a.Conf.SourceDirectory, a.Conf.IgnoreDirectories, a.Conf.IgnoreFiles)
				isFilesystemGraphCreated = true
			}
			for _, lang := range a.Conf.Languages {
				var parser parsers.DependencyExtractor
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
				switch gConf.Name {
				case "file_dependency":
					fmt.Println("TODO %v", gConf.Name)
					depGraph := graphs.NewDependencyGraph(*filesystemGraph, parser, a.Conf.Extensions)
					a.Export(depGraph)
				case "entity_dependency":
					fmt.Println("TODO %v", gConf.Name)
				case "entity_inheritance":
					fmt.Println("TODO %v", gConf.Name)
				case "entity_complete":
					fmt.Println("TODO %v", gConf.Name)
				}
			}
		}
	}
	// graphtest.ParseConfigTest(conf)
	//graphtest.UndirectedGraphCreationTest()
	//graphtest.DirectedGraphCreationTest()
	//graphtest.FileSystemGraphCreationTest(conf)
	//graphtest.Neo4jHelloWorldTest(conf)
	//graphtest.Neo4jLoadingGraphTest(conf)
	//graphtest.GetArangoGraphTest(conf)
	//graphtest.ArangoLoadingGraphTest(conf)
}

func (a *Analyzer) Export(g graphs.Exporter) {} //TODO: export
