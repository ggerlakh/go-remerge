package args

import (
	"flag"
	"go-remerge/internal/config"
	"go-remerge/test/graphtest"
)

func ParseArgs() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "path to yaml config")
	flag.Parse()
	conf := config.ParseConfig(configPath)
	graphtest.ParseConfigTest(conf)
	//a := analyzer.Analyzer{Conf: conf, GraphMap: make(map[string]graphs.Exporter)}
	//a.Start()
}
