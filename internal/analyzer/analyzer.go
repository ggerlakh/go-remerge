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
	Verbose        bool
}

func (a *Analyzer) Start() {
	err := os.Chdir(a.Conf.SourceDirectory)
	if err != nil {
		log.Fatalf("Error to change directory: %v\n", err)
	}
	for _, gConf := range a.Conf.Graphs {
		if gConf.Graph == "filesystem" {
			a.CreateFilesystemGraphIfNotCreated(gConf.Direction, a.Conf.ProjectName)
			a.Export(a.GraphMap["filesystem"], gConf.Graph)
		} else {
			for _, lang := range a.Conf.Languages {
				var parser parsers.CompleteGraphExtractor
				switch strings.ToLower(lang) {
				case "python":
					parser = &parsers.PythonParser{}
				case "golang", "go":
					parser = &parsers.GoParser{ProjectDir: a.Conf.SourceDirectory}
				}
				switch gConf.Graph {
				case "file_dependency":
					a.CreateFileDependencyGraphIfNotCreated(gConf.Direction, parser, a.Conf.ProjectName)
					a.Export(a.GraphMap[gConf.Graph], gConf.Graph)
				case "entity_dependency":
					a.CreateEntityDependencyGraphIfNotCreated(gConf.Direction, parser, a.Conf.ProjectName)
					a.Export(a.GraphMap[gConf.Graph], gConf.Graph)
				case "entity_inheritance":
					a.CreateEntityInheritanceGraphIfNotCreated(gConf.Direction, parser, a.Conf.ProjectName)
					a.Export(a.GraphMap[gConf.Graph], gConf.Graph)
				case "entity_complete":
					a.CreateEntityCompleteGraphIfNotCreated(gConf.Direction, parser, a.Conf.ProjectName)
					a.Export(a.GraphMap[gConf.Graph], gConf.Graph)
				}
			}
		}
	}
}

// CreateFilesystemGraphIfNotCreated fill GraphMap
func (a *Analyzer) CreateFilesystemGraphIfNotCreated(Direction string, AnalysisProjectName string) {
	if _, isFilesystemGraphCreated := a.GraphMap["filesystem"]; !isFilesystemGraphCreated { // create graph if not created
		a.GraphMap["filesystem"] = graphs.NewFileSystemGraph(Direction, []graphs.Node{}, []graphs.Edge{},
			a.Conf.SourceDirectory, a.Conf.IgnoreDirectories, a.Conf.IgnoreFiles, AnalysisProjectName, a.Verbose)
	} else if a.GraphMap["filesystem"].(*graphs.FileSystemGraph).Direction != Direction { // if type not equal type from config, create new graph with correct type
		a.GraphMap["filesystem"] = graphs.NewFileSystemGraph(Direction, []graphs.Node{}, []graphs.Edge{},
			a.Conf.SourceDirectory, a.Conf.IgnoreDirectories, a.Conf.IgnoreFiles, AnalysisProjectName, a.Verbose)
	}
}

// CreateFileDependencyGraphIfNotCreated fill GraphMap
func (a *Analyzer) CreateFileDependencyGraphIfNotCreated(Direction string, parser parsers.DependencyExtractor, AnalysisProjectName string) {
	a.CreateFilesystemGraphIfNotCreated(Direction, AnalysisProjectName)
	if _, isFileDependencyGraphCreated := a.GraphMap["file_dependency"]; !isFileDependencyGraphCreated {
		a.GraphMap["file_dependency"] = graphs.NewFileDependencyGraph(*a.GraphMap["filesystem"].(*graphs.FileSystemGraph),
			parser, a.Conf.Extensions, Direction, AnalysisProjectName)
	} else if a.GraphMap["file_dependency"].(*graphs.DependencyGraph).Direction != Direction {
		a.GraphMap["file_dependency"] = graphs.NewFileDependencyGraph(*a.GraphMap["filesystem"].(*graphs.FileSystemGraph),
			parser, a.Conf.Extensions, Direction, AnalysisProjectName)
	}
}

// CreateEntityDependencyGraphIfNotCreated fill GraphMap
func (a *Analyzer) CreateEntityDependencyGraphIfNotCreated(Direction string, parser parsers.DependencyExtractor, AnalysisProjectName string) {
	a.CreateFileDependencyGraphIfNotCreated(Direction, parser, AnalysisProjectName)
	if _, isEntityDependencyGraphCreated := a.GraphMap["entity_dependency"]; !isEntityDependencyGraphCreated {
		a.GraphMap["entity_dependency"] = graphs.NewEntityDependencyGraph(*a.GraphMap["file_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions, Direction, AnalysisProjectName)
	} else if a.GraphMap["entity_dependency"].(*graphs.DependencyGraph).Direction != Direction {
		a.GraphMap["entity_dependency"] = graphs.NewEntityDependencyGraph(*a.GraphMap["file_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions, Direction, AnalysisProjectName)
	}
}

// CreateEntityInheritanceGraphIfNotCreated fill GraphMap
func (a *Analyzer) CreateEntityInheritanceGraphIfNotCreated(Direction string, parser parsers.CompleteGraphExtractor, AnalysisProjectName string) {
	a.CreateEntityDependencyGraphIfNotCreated(Direction, parser, AnalysisProjectName)
	if _, isEntityInheritanceGraphCreated := a.GraphMap["entity_inheritance"]; !isEntityInheritanceGraphCreated {
		a.GraphMap["entity_inheritance"] = graphs.NewEntityInheritanceGraph(*a.GraphMap["entity_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions, Direction, AnalysisProjectName)
	} else if a.GraphMap["entity_inheritance"].(*graphs.InheritanceGraph).Direction != Direction {
		a.GraphMap["entity_inheritance"] = graphs.NewEntityInheritanceGraph(*a.GraphMap["entity_dependency"].(*graphs.DependencyGraph),
			parser, a.Conf.Extensions, Direction, AnalysisProjectName)
	}
}

// CreateEntityCompleteGraphIfNotCreated fill GraphMap
func (a *Analyzer) CreateEntityCompleteGraphIfNotCreated(Direction string, parser parsers.CompleteGraphExtractor, AnalysisProjectName string) {
	a.CreateEntityDependencyGraphIfNotCreated(Direction, parser, AnalysisProjectName)
	a.CreateEntityInheritanceGraphIfNotCreated(Direction, parser, AnalysisProjectName)
	if _, isEntityCompleteGraphCreated := a.GraphMap["entity_complete"]; !isEntityCompleteGraphCreated {
		a.GraphMap["entity_complete"] = graphs.NewEntityCompleteGraph(*a.GraphMap["entity_dependency"].(*graphs.DependencyGraph),
			*a.GraphMap["entity_inheritance"].(*graphs.InheritanceGraph), Direction, AnalysisProjectName)
	} else if a.GraphMap["entity_complete"].(*graphs.CompleteGraph).Direction != Direction {
		a.GraphMap["entity_complete"] = graphs.NewEntityCompleteGraph(*a.GraphMap["entity_dependency"].(*graphs.DependencyGraph),
			*a.GraphMap["entity_inheritance"].(*graphs.InheritanceGraph), Direction, AnalysisProjectName)
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
						outputJSONFile := filepath.Join(a.Conf.Export.AsJSONFile.OutputDir, fmt.Sprintf("%s_%s.json", a.Conf.ProjectName, graphName))
						err := os.WriteFile(outputJSONFile, []byte(g.GetPrettyJson()), 0644)
						if err != nil {
							log.Fatalf("Error writing json output in file %s\n", outputJSONFile)
						}
						if a.Verbose {
							fmt.Printf("%s graph exported as JSON file in %s\n", graphName, outputJSONFile)
						}
					case "arango_format":
						outputJSONArangoFormatFile := filepath.Join(a.Conf.Export.AsJSONFile.OutputDir, fmt.Sprintf("%s_%sArangoFormat.json", a.Conf.ProjectName, graphName))
						err := os.WriteFile(outputJSONArangoFormatFile, []byte(g.ToArango()), 0644)
						if err != nil {
							log.Fatalf("Error writing json in arango format output in file %s\n", outputJSONArangoFormatFile)
						}
						if a.Verbose {
							fmt.Printf("%s graph exported as ArangoDB formatted JSON file in %s\n", graphName, outputJSONArangoFormatFile)
						}
					}
				}
			}
		}
	}
	if a.ExportTypesMap["arango"] {
		if a.Verbose {
			fmt.Printf("starting export %s graph in ArangoDB...\n", graphName)
		}
		arangoCtx := context.Background()
		arangoEndpoints := a.Conf.Export.Arango.Endpoints
		arangoUsername := a.Conf.Export.Arango.Username
		arangoPassword := a.Conf.Export.Arango.Password
		arangoDatabase := a.Conf.Export.Arango.Database
		g.LoadArangoGraph(arangoCtx, arangoEndpoints, arangoUsername, arangoPassword, arangoDatabase)
		if a.Verbose {
			fmt.Printf("graph %s exported in ArangoDB\n", graphName)
		}
	}
	if a.ExportTypesMap["neo4j"] {
		if a.Verbose {
			fmt.Printf("starting export %s graph in Neo4j...\n", graphName)
		}
		neo4jCtx := context.Background()
		neo4jURI := a.Conf.Export.Neo4j.URI
		neo4jUsername := a.Conf.Export.Neo4j.Username
		neo4jPassword := a.Conf.Export.Neo4j.Password
		err := g.LoadNeo4jGraph(neo4jCtx, neo4jURI, neo4jUsername, neo4jPassword)
		if err != nil {
			log.Fatalf("Error loading graph in neo4j: %s\n", err)
		}
		if a.Verbose {
			fmt.Printf("%s graph exported in Neo4j\n", graphName)
		}
	}
}
