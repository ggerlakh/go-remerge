package args

import (
	"flag"
	"go-remerge/internal/analyzer"
	"go-remerge/internal/config"
)

func ParseArgs() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "path to yaml config")
	flag.Parse()
	conf := config.ParseConfig(configPath)
	a := analyzer.Analyzer{Conf: conf}
	a.Start()
}
