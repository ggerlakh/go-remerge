package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type GraphConfig struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
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
type ExportConfig struct {
	AsFile []string     `yaml:"as_file"`
	Arango ArangoConfig `yaml:"arango"`
	Neo4j  Neo4jConfig  `yaml:"neo4j"`
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
	// TODO: validate config
	// TODO: omitempty
	var conf Config
	yamlFile, err := os.ReadFile(configPath)
	//fmt.Println(string(yamlFile))
	if err != nil {
		log.Fatal("Error reading yaml config: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatal("Error parsing yaml config %v", err)
	}
	return conf
}
