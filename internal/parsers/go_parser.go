package parsers

import (
	"bufio"
	"fmt"
	"go-remerge/tools/ostool"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type GoParser struct {
	ProjectDir string
}

func (Parser *GoParser) ExtractInheritance(filePath, entityName string) []map[string]string {
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
					if genDecl.Tok == token.TYPE || genDecl.Tok == token.FUNC || genDecl.Tok == token.CONST {
						for _, spec := range genDecl.Specs {
							typeSpec := spec.(*ast.TypeSpec)
							lines := ostool.FilterComments(filePath)
							regex := `^.*` + filepath.Base(path) + `\.` + typeSpec.Name.Name + `.*$`
							imported := regexp.MustCompile(regex)
							for _, line := range lines {
								if imported.MatchString(line) {
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

func (Parser *GoParser) ExtractExternalEntities(externalDependencyName, fromNodePath string) []string {
	var externalEntityDependencies []string
	file, err := os.Open(fromNodePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	splittedName := strings.Split(externalDependencyName, string(filepath.Separator))
	regex := strings.Split(splittedName[len(splittedName)-1], ".")[0] + `\.\w+`
	importedExternalDependency := regexp.MustCompile(regex)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "//") && !strings.Contains(line, "/*") {
			tmpRes := importedExternalDependency.Find([]byte(line))
			if string(tmpRes) != "" {
				externalEntityDependencies = append(externalEntityDependencies, string(tmpRes))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return externalEntityDependencies
}

func (Parser *GoParser) HasEntityDependency(fromEntityName, fromEntityPath, toEntityName, toEntityPackage string) bool {
	var hasEntityDependency bool
	fromNodeBytes, err := os.ReadFile(fromEntityPath)
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fromEntityPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Error parsing file: %v\n", err)
	}
	// Extract function source code
	for _, decl := range node.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Recv != nil {
					switch fn.Recv.List[0].Type.(type) {
					case *ast.StarExpr:
						//fmt.Printf("Recevier: %v\n", fn.Recv.List[0].Type.(*ast.StarExpr).X)
						if fmt.Sprintf("%v", fn.Recv.List[0].Type.(*ast.StarExpr).X) == fromEntityName {
							startLine := fset.Position(fn.Pos()).Line
							endLine := fset.Position(fn.End()).Line
							regex := `^.*` + toEntityPackage + `\.` + toEntityName + `.*$`
							importedObj := regexp.MustCompile(regex)
							// if field struct method has imported obj return true
							for _, funcLine := range strings.Split(string(fromNodeBytes), "\n")[startLine-1 : endLine] {
								if importedObj.MatchString(funcLine) && !strings.Contains(funcLine, "//") && !strings.Contains(funcLine, "/*") {
									return true
								}
							}
						}
					case *ast.Ident:
						//fmt.Printf("Recevier: %v\n", fn.Recv.List[0].Type.(*ast.Ident).Name)
						if fn.Recv.List[0].Type.(*ast.Ident).Name == fromEntityName {
							startLine := fset.Position(fn.Pos()).Line
							endLine := fset.Position(fn.End()).Line
							regex := `^.*` + toEntityPackage + `\.` + toEntityName + `.*$`
							importedObj := regexp.MustCompile(regex)
							// if field struct method has imported obj return true
							for _, funcLine := range strings.Split(string(fromNodeBytes), "\n")[startLine-1 : endLine] {
								if importedObj.MatchString(funcLine) && !strings.Contains(funcLine, "//") && !strings.Contains(funcLine, "/*") {
									return true
								}
							}
						}
					}
				}
			}
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			if genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							// extract struct source code
							startLine := fset.Position(structType.Pos()).Line
							endLine := fset.Position(structType.End()).Line
							regex := `^.*` + toEntityPackage + `\.` + toEntityName + `.*$`
							importedObj := regexp.MustCompile(regex)
							// if field struct field has imported obj return true
							for _, structLine := range strings.Split(string(fromNodeBytes), "\n")[startLine-1 : endLine] {
								if importedObj.MatchString(structLine) && !strings.Contains(structLine, "//") && !strings.Contains(structLine, "/*") {
									return true
								}
							}
						}
					}
				}
			}
		}
	}
	return hasEntityDependency
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
