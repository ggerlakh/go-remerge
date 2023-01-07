package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"reflect"
)

type GraphConfig struct {
	Graph string `yaml:"graph"`
	Type  string `yaml:"type"`
}

type ArangoConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	URL      string `yaml:"url"`
	Database string `yaml:"database"`
}
type Neo4jConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	URI      string `yaml:"uri"`
}
type FileExportConfig struct {
	OutputDir string   `yaml:"output_dir"`
	Formats   []string `yaml:"formats"`
}
type ExportConfig struct {
	AsFile FileExportConfig `yaml:"as_file"`
	Arango ArangoConfig     `yaml:"arango"`
	Neo4j  Neo4jConfig      `yaml:"neo4j"`
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

func ParseConfig(configPath string) Config {
	var conf Config
	var defaultExport ExportConfig
	var defautGraphConfig []GraphConfig
	var defaultLanguages, defaultExtensions []string
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading yaml config: %v\n", err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Error parsing yaml config %v\n", err)
	}
	// check mandatory config fields
	if _, err := os.Stat(conf.SourceDirectory); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("source directory %s does not exist\n", conf.SourceDirectory)
		}
	}
	if reflect.DeepEqual(conf.Export, defaultExport) {
		log.Fatalf("Export parameters are not set in %s\n", configPath)
	}
	if reflect.DeepEqual(conf.Languages, defaultLanguages) {
		log.Fatalf("Languages are not set in %s\n", configPath)
	}
	if reflect.DeepEqual(conf.Extensions, defaultExtensions) {
		log.Fatalf("Extensions are not set in %s\n", configPath)
	}
	if reflect.DeepEqual(conf.Graphs, defautGraphConfig) {
		log.Fatalf("Graphs are not set in %s\n", configPath)
	}
	return conf
}
