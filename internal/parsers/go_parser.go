package parsers

import (
	"go-remerge/tools/ostool"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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
		//openPath := path
		// Parse the imported file
		if !ostool.Exists(path) {
			//fmt.Println("Replace: ", strings.ReplaceAll(filepath.Join(filepath.Join(Parser.ProjectDir, ".."), filepath.Clean(path)), Parser.ProjectDir, ""))
			path = strings.TrimPrefix(strings.ReplaceAll(filepath.Join(filepath.Join(Parser.ProjectDir, ".."), filepath.Clean(path)), Parser.ProjectDir, ""), string(filepath.Separator))
			//path = strings.ReplaceAll(filepath.Join(filepath.Join(Parser.ProjectDir, ".."), filepath.Clean(path)), Parser.ProjectDir, "")
			if !ostool.Exists(path) {
				path = strings.TrimLeft(strings.ReplaceAll(filepath.Join(filepath.Join(Parser.ProjectDir, ".."), filepath.Clean(path)), filepath.Dir(Parser.ProjectDir), ""), string(filepath.Separator))
				//fmt.Printf("Skipping external dependency %v\n", path)
				fileDependenciesMap[filepath.Join("external_dependency", path)] = struct{}{}
				continue
			}
		}
		//fmt.Printf("File: %s, Import path: %s\n", filePath, path)
		var packageGoFiles []string
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		files, err := f.Readdir(0)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".go") {
				packageGoFiles = append(packageGoFiles, filepath.Join(path, file.Name()))
			}
		}
		//packageGoFiles, err := filepath.Glob(filepath.Join(path, "*.go"))
		//fmt.Printf("Package go files: %s\n", packageGoFiles)
		for _, packageGoFile := range packageGoFiles {
			//fmt.Printf("Package go file: %s\n", packageGoFile)
			importedNode, err := parser.ParseFile(fset, packageGoFile, nil, parser.ParseComments)
			if err != nil {
				panic(err)
			}
			//packageGoFile = strings.TrimPrefix(strings.ReplaceAll(packageGoFile, Parser.ProjectDir, ""), string(filepath.Separator))
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
	var entityResult []string
	if !ostool.Exists(filePath) {
		//fmt.Printf("Path is not exist %s\n", filePath)
		return []string{}
	}
	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	// Iterate through the top-level declarations and find the structures
	for _, decl := range fileNode.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			if genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					typeSpec := spec.(*ast.TypeSpec)
					entityResult = append(entityResult, typeSpec.Name.Name)
				}
			}
		}
	}
	return entityResult
}

func (Parser *GoParser) ExtractPackage(filePath string) string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		if !ostool.Exists(filePath) {
			return strings.ReplaceAll(filePath, "external_dependency"+string(filepath.Separator), "")
		} else {
			panic(err)
		}
	}
	return node.Name.Name
}
