package graphs

import (
	"fmt"
	"go-remerge/tools/hashtool"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type FileSystemGraph struct {
	Graph
	Root      string
	SkipDirs  map[string]struct{}
	SkipFiles map[string]struct{}
}

func NewFileSystemGraph(direction string, Nodes []Node, Edges []Edge, Root string, SkipDirs []string, SkipFiles []string) *FileSystemGraph {
	if strings.ToLower(direction) == "undirected" || strings.ToLower(direction) == "directed" {
		fsG := &FileSystemGraph{Graph: Graph{
			Direction: direction,
			Name:      "filesystem",
			Nodes:     make(map[string]Node),
			Edges:     make(map[string]Edge)},
			Root:      Root,
			SkipDirs:  make(map[string]struct{}),
			SkipFiles: make(map[string]struct{}),
		}
		fsG.SetNodes(Nodes)
		fsG.SetEdges(Edges)
		fsG.SetSkipDirs(SkipDirs)
		fsG.SetSkipFiles(SkipFiles)
		fsG.WalkTree()
		return fsG
	} else {
		panic(fmt.Sprintf("\"%v\" wrong graphs type value, graphs can be only directed or undirected", direction))
	}
}

func (fsG *FileSystemGraph) SetSkipDirs(SkipDirs []string) {
	for _, dir := range SkipDirs {
		fsG.SkipDirs[dir] = struct{}{}
	}
}

func (fsG *FileSystemGraph) SetSkipFiles(SkipFiles []string) {
	for _, file := range SkipFiles {
		fsG.SkipFiles[file] = struct{}{}
	}
}

func (fsG *FileSystemGraph) GetRootDir(Root string) string {
	abs, err := filepath.Abs(Root)
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Base(abs)
}

func (fsG *FileSystemGraph) WalkTree() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if _, skip := fsG.SkipDirs[info.Name()]; skip && info.IsDir() {
			log.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		} else if _, skip := fsG.SkipFiles[info.Name()]; skip && !info.IsDir() {
			log.Printf("skipping a file without errors: %+v \n", info.Name())
		} else {
			var fromPath, toPath string
			if path == "." {
				fsG.AddNode(Node{Id: hashtool.Sha256(fsG.GetRootDir(fsG.Root)), Labels: map[string]any{
					"name":        filepath.Base(fsG.GetRootDir(fsG.Root)),
					"path":        fsG.GetRootDir(fsG.Root),
					"isDirectory": info.IsDir()}})
			} else {
				// normalize path: all must start with Root dir
				if filepath.Dir(path) == "." {
					fromPath = fsG.GetRootDir(fsG.Root)
				} else {
					fromPath = filepath.Dir(path)
				}
				toPath = path
				// adding "from" Node if not exists
				if _, nodeExists := fsG.Nodes[hashtool.Sha256(fromPath)]; !nodeExists {
					fsG.AddNode(Node{Id: hashtool.Sha256(fromPath), Labels: map[string]any{
						"name":        filepath.Base(fromPath),
						"path":        fromPath,
						"isDirectory": info.IsDir()}})
				}
				// adding "to" Node if not exists
				if _, nodeExists := fsG.Nodes[hashtool.Sha256(toPath)]; !nodeExists {
					fsG.AddNode(Node{Id: hashtool.Sha256(toPath), Labels: map[string]any{
						"name":        filepath.Base(toPath),
						"path":        toPath,
						"isDirectory": info.IsDir()}})
				}
				if fromPath != toPath {
					fsG.AddEdge(Edge{From: fsG.Nodes[hashtool.Sha256(fromPath)], To: fsG.Nodes[hashtool.Sha256(toPath)]})
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
