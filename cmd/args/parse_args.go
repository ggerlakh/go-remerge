package args

import (
	"flag"
	"fmt"
	"go-remerge/internal/analyzer"
	"go-remerge/internal/config"
	"go-remerge/internal/graphs"
	"log"
	"os"
	"path/filepath"
)

func ParseArgs() {
	var configPath string
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.StringVar(&configPath, "c", "", "path to yaml config")
	flag.Parse()
	//conf, _ := config.ParseConfig(configPath)
	//graphtest.ParseConfigTest(conf)
	// get info about config path
	fi, err := os.Stat(configPath)
	if err != nil {
		log.Fatalf("Error while opening %v: %v", configPath, err)
	}
	if fi.IsDir() {
		dirItems, err := os.ReadDir(configPath)
		if err != nil {
			log.Fatalf("Error while opening directory %v: %v", configPath, err)
		}
		for _, item := range dirItems {
			confPath := filepath.Join(configPath, item.Name())
			if !item.IsDir() && (filepath.Ext(item.Name()) == ".yml" || filepath.Ext(item.Name()) == ".yaml") {
				runTask(confPath)
				fmt.Printf("Task with config %v done\n", confPath)
			}
		}
	} else {
		runTask(configPath)
	}
}

func runTask(configPath string) {
	conf, exportTypesMap := config.ParseConfig(configPath)
	a := analyzer.Analyzer{Conf: conf, GraphMap: make(map[string]graphs.Exporter), ExportTypesMap: exportTypesMap}
	a.Start()
}
