package parsers

import (
	"go-remerge/tools/ostool"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
)

type GoParser struct {
	ProjectDir string
}

func (Parser *GoParser) ExtractInheritance(filePath, entityName string) []string {
	//TODO implement me
	panic("implement me")
}

func (Parser *GoParser) ExtractDependencies(filePath string) []string {
	var fileResults []string
	var fileDependenciesMap = make(map[string]struct{})
	// Specify the path of the Go file to analyze
	// Parse the file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	// Iterate through the imports and extract the structure objects
	for _, i := range node.Imports {
		// Get the import path
		path := i.Path.Value[1 : len(i.Path.Value)-1]
		// Parse the imported file
		if !ostool.Exists(path) {
			absPath := filepath.Join(filepath.Join(Parser.ProjectDir, ".."), filepath.Clean(path))
			if !ostool.Exists(absPath) {
				//fmt.Printf("Skipping external dependency %v\n", path)
				fileDependenciesMap["external_dependency/"+path] = struct{}{}
				continue
			}
		}
		packageGoFiles, err := filepath.Glob(filepath.Join(path, "*.go"))
		if err != nil {
			panic(err)
		}
		for _, packageGoFile := range packageGoFiles {
			importedNode, err := parser.ParseFile(fset, packageGoFile, nil, parser.ParseComments)
			if err != nil {
				panic(err)
			}
			// Iterate through the top-level declarations and find the structures
			for _, decl := range importedNode.Decls {
				switch decl.(type) {
				case *ast.GenDecl:
					genDecl := decl.(*ast.GenDecl)
					if genDecl.Tok == token.TYPE {
						for _, spec := range genDecl.Specs {
							typeSpec := spec.(*ast.TypeSpec)
							lines := ostool.FilterComments(filePath)
							regex := `^.*` + filepath.Base(path) + `\.` + typeSpec.Name.Name + `.*$`
							importedStruct := regexp.MustCompile(regex)
							for _, line := range lines {
								if importedStruct.MatchString(line) {
									fileDependenciesMap[packageGoFile] = struct{}{}
								}
							}
						}
					}
				}
			}
		}
	}
	for file, _ := range fileDependenciesMap {
		fileResults = append(fileResults, file)
	}
	return fileResults
}

func (Parser *GoParser) ExtractEntities(filePath string) []string {
	return []string{}
}

func (Parser *GoParser) ExtractPackage(filePath string) string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		if !ostool.Exists(filePath) {
			return strings.ReplaceAll(filePath, "external_dependency/", "")
		} else {
			panic(err)
		}
	}
	return node.Name.Name
}
