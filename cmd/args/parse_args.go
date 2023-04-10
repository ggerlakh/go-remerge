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
	"sync"
)

func ParseArgs() {
	var configPath string
	var async bool
	var verbose bool
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.StringVar(&configPath, "c", "", "path to yaml config")
	flag.BoolVar(&verbose, "v", false, "produce verbose output")
	flag.BoolVar(&async, "async", false, "asynchronous task execution")
	flag.Usage = func() {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s -c <path> [-h] [-v] [--async]:\n-h, --help\tprint help information\n", os.Args[0])
		if err != nil {
			log.Fatal(err)
		}

		flag.VisitAll(func(f *flag.Flag) {
			var prefix string
			if f.Name == "v" || f.Name == "c" {
				prefix = "-"
			} else if f.Name == "async" {
				prefix = "--"
			}
			_, err := fmt.Fprintf(os.Stderr, "%v%v\t\t%v\n", prefix, f.Name, f.Usage)
			if err != nil {
				log.Fatal(err)
			} // f.Name, f.Value
		})
	}
	flag.Parse()
	wg := new(sync.WaitGroup)
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
		fmt.Printf("Searching for configs in %v...\n", configPath)
		for _, item := range dirItems {
			confPath := filepath.Join(configPath, item.Name())
			if !item.IsDir() && (filepath.Ext(item.Name()) == ".yml" || filepath.Ext(item.Name()) == ".yaml") {
				if async {
					wg.Add(1)
					fmt.Printf("Starting async task from config %v...\n", configPath)
					go runTask(confPath, async, verbose, wg)
				} else {
					runTask(confPath, async, verbose, wg)
				}
			}
		}
	} else {
		if async {
			wg.Add(1)
			go runTask(configPath, async, verbose, wg)
		} else {
			runTask(configPath, async, verbose, wg)
		}
	}
	if async {
		wg.Wait()
	}
}

func runTask(configPath string, async, verbose bool, wg *sync.WaitGroup) {
	fmt.Printf("Starting task from config %v...\n", configPath)
	if async {
		defer wg.Done()
	}
	conf, exportTypesMap := config.ParseConfig(configPath)
	a := analyzer.Analyzer{Conf: conf, GraphMap: make(map[string]graphs.Exporter), ExportTypesMap: exportTypesMap, Verbose: verbose}
	a.Start()
	fmt.Printf("Task from config %v was completed successfully\n", configPath)
}
