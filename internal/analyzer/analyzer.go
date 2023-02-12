package analyzer

import (
	"context"
	"errors"
	"fmt"
	"go-remerge/internal/config"
	"go-remerge/internal/graphs"
	"go-remerge/internal/parsers"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Analyzer struct {
	Conf           config.Config
	GraphMap       map[string]graphs.Exporter
	ExportTypesMap map[string]bool
}

func (a *Analyzer) Start() {
	for _, gConf := range a.Conf.Graphs {
		if gConf.Graph == "filesystem" {
			a.CreateFilesystemGraphIfNotCreated(gConf.Type)
			fmt.Println(a.GraphMap["filesystem"])
			//fmt.Println("break")
			//break
			a.Export(a.GraphMap["filesystem"], gConf.Graph)
		} else {
			for _, lang := range a.Conf.Languages {
				var parser parsers.CompleteGraphExtractor
				switch strings.ToLower(lang) {
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
					a.Export(a.GraphMap[gConf.Graph], gConf.Graph)
				case "entity_dependency":
					fmt.Printf("TODO %v\n", gConf.Graph)
					a.CreateEntityDependencyGraphIfNotCreated(gConf.Type, parser)
					a.Export(a.GraphMap[gConf.Graph], gConf.Graph)
				case "entity_inheritance":
					fmt.Printf("TODO %v\n", gConf.Graph)
					a.CreateEntityInheritanceGraphIfNotCreated(gConf.Type, parser)
					a.Export(a.GraphMap[gConf.Graph], gConf.Graph)
				case "entity_complete":
					fmt.Printf("TODO %v\n", gConf.Graph)
					a.CreateEntityCompleteGraphIfNotCreated(gConf.Type, parser)
					a.Export(a.GraphMap[gConf.Graph], gConf.Graph)
				}
			}
		}
	}
}

// CreateFilesystemGraphIfNotCreated fill GraphMap
func (a *Analyzer) CreateFilesystemGraphIfNotCreated(Type string) {
	if _, isFilesystemGraphCreated := a.GraphMap["filesystem"]; !isFilesystemGraphCreated { // create graph if not created
		a.GraphMap["filesystem"] = graphs.NewFileSystemGraph(Type, "filesystem", []graphs.Node{}, []graphs.Edge{},
			a.Conf.SourceDirectory, a.Conf.IgnoreDirectories, a.Conf.IgnoreFiles)
	} else if a.GraphMap["filesystem"].(*graphs.FileSystemGraph).Direction != Type { // if type not equal type from config, create new graph with correct type
		a.GraphMap["filesystem"] = graphs.NewFileSystemGraph(Type, "filesystem", []graphs.Node{}, []graphs.Edge{},
			a.Conf.SourceDirectory, a.Conf.IgnoreDirectories, a.Conf.IgnoreFiles)
	}
}

// CreateFileDependencyGraphIfNotCreated fill GraphMap
func (a *Analyzer) CreateFileDependencyGraphIfNotCreated(Type string, parser parsers.DependencyExtractor) {
	a.CreateFilesystemGraphIfNotCreated(Type)
	if _, isFileDependencyGraphCreated := a.GraphMap["file_dependency"]; !isFileDependencyGraphCreated {
		a.GraphMap["file_dependency"] = graphs.NewFileDependencyGraph(*a.GraphMap["filesystem"].(*graphs.FileSystemGraph),
			parser, a.Conf.Extensions)
	} else if a.GraphMap["file_dependency"].(*graphs.DependencyGraph).Direction != Type {
		a.GraphMap["file_dependency"] = graphs.NewFileDependencyGraph(*a.GraphMap["filesystem"].(*graphs.FileSystemGraph),
			parser, a.Conf.Extensions)
	}
}

// CreateEntityDependencyGraphIfNotCreated fill GraphMap
func (a *Analyzer) CreateEntityDependencyGraphIfNotCreated(Type string, parser parsers.DependencyExtractor) {
	a.CreateFilesystemGraphIfNotCreated(Type)
	if _, isEntityDependencyGraphCreated := a.GraphMap["entity_dependency"]; !isEntityDependencyGraphCreated {
		a.GraphMap["entity_dependency"] = graphs.NewEntityDependencyGraph(*a.GraphMap["file_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions)
	} else if a.GraphMap["entity_dependency"].(*graphs.DependencyGraph).Direction != Type {
		a.GraphMap["entity_dependency"] = graphs.NewEntityDependencyGraph(*a.GraphMap["file_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions)
	}
}

// CreateEntityInheritanceGraphIfNotCreated fill GraphMap
func (a *Analyzer) CreateEntityInheritanceGraphIfNotCreated(Type string, parser parsers.InheritanceExtractor) {
	a.CreateFilesystemGraphIfNotCreated(Type)
	if _, isEntityInheritanceGraphCreated := a.GraphMap["entity_inheritance"]; !isEntityInheritanceGraphCreated {
		a.GraphMap["entity_inheritance"] = graphs.NewEntityInheritanceGraph(*a.GraphMap["file_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions)
	} else if a.GraphMap["entity_inheritance"].(*graphs.DependencyGraph).Direction != Type {
		a.GraphMap["entity_inheritance"] = graphs.NewEntityInheritanceGraph(*a.GraphMap["file_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions)
	}
}

// CreateEntityCompleteGraphIfNotCreated fill GraphMap
func (a *Analyzer) CreateEntityCompleteGraphIfNotCreated(Type string, parser parsers.CompleteGraphExtractor) {
	a.CreateEntityInheritanceGraphIfNotCreated(Type, parser)
	a.CreateEntityDependencyGraphIfNotCreated(Type, parser)
	if _, isEntityCompleteGraphCreated := a.GraphMap["entity_complete"]; !isEntityCompleteGraphCreated {
		a.GraphMap["entity_complete"] = graphs.NewEntityCompleteGraph(*a.GraphMap["entity_dependency"].(*graphs.DependencyGraph),
			*a.GraphMap["entity_inheritance"].(*graphs.InheritanceGraph))
	} else if a.GraphMap["entity_complete"].(*graphs.CompleteGraph).Direction != Type {
		a.GraphMap["entity_complete"] = graphs.NewEntityCompleteGraph(*a.GraphMap["entity_dependency"].(*graphs.DependencyGraph),
			*a.GraphMap["entity_inheritance"].(*graphs.InheritanceGraph))
	}
}

func (a *Analyzer) Export(g graphs.Exporter, graphName string) {
	// export as json file
	if a.ExportTypesMap["json"] {
		// check if path exists
		if fi, err := os.Stat(a.Conf.Export.AsJSONFile.OutputDir); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Fatalf("output directory path %s does not exist\n", a.Conf.Export.AsJSONFile.OutputDir)
			} else if !fi.IsDir() {
				log.Fatalf("output path must be a directory\n")
			}
		} else if len(a.Conf.Export.AsJSONFile.Formats) == 0 || len(a.Conf.Export.AsJSONFile.Formats) > 2 {
			log.Fatalf("Wrong formats quantity at export.as_json_file block. Quantity must be 1 or 2\n")
		} else {
			for _, format := range a.Conf.Export.AsJSONFile.Formats {
				format = strings.ToLower(format)
				if format != "json" && format != "arango_format" {
					log.Fatalf("Wrong format parameter in export.as_json_file.format block. Format may be \"json\" or \"arango_format\"\n")
				} else {
					switch format {
					case "json":
						outputJSONFile := filepath.Join(a.Conf.Export.AsJSONFile.OutputDir, fmt.Sprintf("%s.json", graphName))
						fmt.Println("outputJSONFileName: ", outputJSONFile)
						fmt.Println(g.GetNodes(), len(g.GetNodes()))
						err := os.WriteFile(outputJSONFile, []byte(g.GetPrettyJson()), 0644)
						if err != nil {
							log.Fatalf("Error writing json output in file %s\n")
						}
					case "arango_format":
						outputJSONArangoFormatFile := filepath.Join(a.Conf.Export.AsJSONFile.OutputDir, fmt.Sprintf("%sArangoFormat.json", graphName))
						fmt.Println("outputJSONArangoFormatFileName: ", outputJSONArangoFormatFile)
						err := os.WriteFile(outputJSONArangoFormatFile, []byte(g.ToArango()), 0644)
						if err != nil {
							log.Fatalf("Error writing json in arango format output in file %s\n")
						}
					}
				}
			}
		}
	}
	if a.ExportTypesMap["arango"] {
		arangoCtx := context.Background()
		arangoEndpoints := a.Conf.Export.Arango.Endpoints
		arangoUsername := a.Conf.Export.Arango.Username
		arangoPassword := a.Conf.Export.Arango.Password
		arangoDatabase := a.Conf.Export.Arango.Database
		g.LoadArangoGraph(arangoCtx, arangoEndpoints, arangoUsername, arangoPassword, arangoDatabase)
	}
	if a.ExportTypesMap["neo4j"] {
		neo4jCtx := context.Background()
		neo4jURI := a.Conf.Export.Neo4j.URI
		neo4jUsername := a.Conf.Export.Neo4j.Username
		neo4jPassword := a.Conf.Export.Neo4j.Password
		err := g.LoadNeo4jGraph(neo4jCtx, neo4jURI, neo4jUsername, neo4jPassword)
		if err != nil {
			log.Fatalf("Error loading graph in neo4j: %s\n", err)
		}
	}
}
