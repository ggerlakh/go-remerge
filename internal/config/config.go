package config

import (
	"go-remerge/tools/ostool"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"reflect"
)

type GraphConfig struct {
	Graph     string `yaml:"graph"`
	Direction string `yaml:"direction"`
}

type ArangoConfig struct {
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	Endpoints []string `yaml:"endpoints"`
	Database  string   `yaml:"database"`
}
type Neo4jConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	URI      string `yaml:"uri"`
}
type JsonExportConfig struct {
	OutputDir string   `yaml:"output_dir"`
	Formats   []string `yaml:"formats"`
}
type ExportConfig struct {
	AsJSONFile JsonExportConfig `yaml:"as_file"`
	Arango     ArangoConfig     `yaml:"arango"`
	Neo4j      Neo4jConfig      `yaml:"neo4j"`
}
type Config struct {
	ProjectName       string        `yaml:"project_name"`
	AnalysisName      string        `yaml:"analysis_name"`
	SourceDirectory   string        `yaml:"source_directory"`
	Languages         []string      `yaml:"languages"`
	Extensions        []string      `yaml:"extensions"`
	IgnoreDirectories []string      `yaml:"ignore_directories"`
	IgnoreFiles       []string      `yaml:"ignore_files"`
	Graphs            []GraphConfig `yaml:"graphs"`
	Export            ExportConfig  `yaml:"export"`
}

func ParseConfig(configPath string) (Config, map[string]bool) {
	var conf Config
	var defaultExport ExportConfig
	var defaultJsonExportConfig JsonExportConfig
	var defaultNeo4jConfig Neo4jConfig
	var defaultArangoConfig ArangoConfig
	var defaultGraphConfig []GraphConfig
	var defaultLanguages, defaultExtensions []string
	// map for check if export type is set
	exportTypesMap := map[string]bool{"json": false, "arango": false, "neo4j": false}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading yaml config: %v\n", err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Error parsing yaml config %v\n", err)
	}
	// check mandatory config fields
	// check if path exists
	if !ostool.Exists(conf.SourceDirectory) {
		log.Fatalf("source directory %s does not exist\n", conf.SourceDirectory)
	}
	// check export fields in config
	if reflect.DeepEqual(conf.Export, defaultExport) {
		log.Fatalf("Export parameters are not set in %s\n", configPath)
	}
	// check json file export in config
	if !reflect.DeepEqual(conf.Export.AsJSONFile, defaultJsonExportConfig) {
		exportTypesMap["json"] = true
	}
	// check arangodb export in config
	if !reflect.DeepEqual(conf.Export.Arango, defaultArangoConfig) {
		exportTypesMap["arango"] = true
	}
	// check neo4j export in config
	if !reflect.DeepEqual(conf.Export.Neo4j, defaultNeo4jConfig) {
		exportTypesMap["neo4j"] = true
	}
	if reflect.DeepEqual(conf.Languages, defaultLanguages) {
		log.Fatalf("Languages are not set in %s\n", configPath)
	}
	if reflect.DeepEqual(conf.Extensions, defaultExtensions) {
		log.Fatalf("Extensions are not set in %s\n", configPath)
	}
	if reflect.DeepEqual(conf.Graphs, defaultGraphConfig) {
		log.Fatalf("Graphs are not set in %s\n", configPath)
	}
	return conf, exportTypesMap
}
