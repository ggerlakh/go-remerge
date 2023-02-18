package args

import (
	"flag"
	"go-remerge/internal/analyzer"
	"go-remerge/internal/config"
	"go-remerge/internal/graphs"
)

func ParseArgs() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "path to yaml config")
	flag.Parse()
	//conf, _ := config.ParseConfig(configPath)
	//graphtest.ParseConfigTest(conf)
	conf, exportTypesMap := config.ParseConfig(configPath)
	a := analyzer.Analyzer{Conf: conf, GraphMap: make(map[string]graphs.Exporter), ExportTypesMap: exportTypesMap}
	a.Start()
}
