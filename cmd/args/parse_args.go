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
	var async bool
	var verbose bool
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.StringVar(&configPath, "c", "", "path to yaml config")
	flag.BoolVar(&verbose, "v", false, "produce verbose output")
	flag.BoolVar(&async, "async", false, "asynchronous task execution")
	flag.Parse()
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
				runTask(confPath, verbose)
				fmt.Printf("Task with config %v was completed successfully\n", configPath)
			}
		}
	} else {
		runTask(configPath, verbose)
		fmt.Printf("Task with config %v was completed successfully\n", configPath)
	}
}

func runTask(configPath string, verbose bool) {
	conf, exportTypesMap := config.ParseConfig(configPath)
	a := analyzer.Analyzer{Conf: conf, GraphMap: make(map[string]graphs.Exporter), ExportTypesMap: exportTypesMap, Verbose: verbose}
	a.Start()
}
