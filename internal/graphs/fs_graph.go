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

func NewFileSystemGraph(Type string, Nodes []Node, Edges []Edge, Root string, SkipDirs []string, SkipFiles []string) *FileSystemGraph {
	if strings.ToLower(Type) == "undirected" || strings.ToLower(Type) == "directed" {
		fsG := &FileSystemGraph{Graph: Graph{
			Type:  Type,
			Nodes: make(map[string]Node),
			Edges: make(map[string]Edge)},
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
		panic(fmt.Sprintf("\"%v\" wrong graphs type value, graphs can be only directed or undirected", Type))
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
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if _, skip := fsG.SkipDirs[info.Name()]; skip && info.IsDir() {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		} else if _, skip := fsG.SkipFiles[info.Name()]; skip && !info.IsDir() {
			fmt.Printf("skipping a file without errors: %+v \n", info.Name())
		} else {
			var fromPath, toPath string
			if path == "." {
				fsG.AddNode(Node{Id: hashtool.Sha256(fsG.GetRootDir(fsG.Root)), Labels: map[string]any{
					"path":        fsG.GetRootDir(fsG.Root),
					"isDirectory": info.IsDir()}})
			} else {
				// normalize path: all must start with Root dir
				if !strings.HasPrefix(path, fsG.Root) {
					path = filepath.Join(fsG.GetRootDir(fsG.Root), path)
				}
				if filepath.Dir(path) == "." {
					fromPath = fsG.GetRootDir(fsG.Root)
				} else {
					fromPath = filepath.Dir(path)
				}
				toPath = path
				// adding "from" Node if not exists
				if _, nodeExists := fsG.Nodes[hashtool.Sha256(fromPath)]; !nodeExists {
					fsG.AddNode(Node{Id: hashtool.Sha256(fromPath), Labels: map[string]any{
						"path":        fromPath,
						"isDirectory": info.IsDir()}})
				}
				// adding "to" Node if not exists
				if _, nodeExists := fsG.Nodes[hashtool.Sha256(toPath)]; !nodeExists {
					fsG.AddNode(Node{Id: hashtool.Sha256(toPath), Labels: map[string]any{
						"path":        toPath,
						"isDirectory": info.IsDir()}})
				}
				if fromPath != toPath {
					fsG.AddEdge(Edge{From: fsG.Nodes[hashtool.Sha256(fromPath)], To: fsG.Nodes[hashtool.Sha256(toPath)]})
					fmt.Println(fromPath + "->" + toPath)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
