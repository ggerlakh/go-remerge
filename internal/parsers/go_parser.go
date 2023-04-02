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

func (Parser *GoParser) ExtractInheritance(entityFilePath, entityName string) []map[string]string {
	// inheritance in Golang implemented by embedding one struct to another
	var parentInheritanceEntities []map[string]string
	var pkgEntities []map[string]string
	fromNodeBytes, err := os.ReadFile(entityFilePath)
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}
	fset := token.NewFileSet()
	// extract all entities that is in the same package with given entity
	file, err := parser.ParseFile(fset, entityFilePath, nil, parser.PackageClauseOnly)
	if err != nil {
		log.Fatalf("Error parsing file: %s\n", err)
	}
	// Get the package name.
	pkgName := file.Name.Name
	// Get the package directory.
	pkgDir := filepath.Dir(entityFilePath)
	// Read all the files in the package directory.
	files, err := os.ReadDir(pkgDir)
	if err != nil {
		log.Fatalf("Error reading package directory: %s\n", err)
	}
	// Loop through all the files in the package directory.
	for _, fileInfo := range files {
		if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".go" {
			// Parse the file and check if it belongs to the same package.
			fileParse, err := parser.ParseFile(fset, filepath.Join(pkgDir, fileInfo.Name()), nil, parser.PackageClauseOnly)
			if err == nil && fileParse.Name.Name == pkgName {
				fileNode, err := parser.ParseFile(fset, filepath.Join(pkgDir, fileInfo.Name()), nil, parser.ParseComments)
				if err != nil {
					panic(err)
				}
				// Iterate through the top-level declarations and find the structures
				for _, decl := range fileNode.Decls {
					switch decl.(type) {
					case *ast.GenDecl:
						genDecl := decl.(*ast.GenDecl)
						if genDecl.Tok == token.TYPE || genDecl.Tok == token.FUNC || genDecl.Tok == token.CONST {
							for _, spec := range genDecl.Specs {
								typeSpec := spec.(*ast.TypeSpec)
								pkgEntities = append(pkgEntities, map[string]string{
									"name": typeSpec.Name.Name,
									"path": filepath.Join(pkgDir, fileInfo.Name()),
								})
								//fmt.Printf("Found struct %q in file %q\n", typeSpec.Name.Name, fileInfo.Name())
							}
						}
					}
				}
			}
		}
	}
	fset = token.NewFileSet()
	node, err := parser.ParseFile(fset, entityFilePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Error parsing file: %v\n", err)
	}
	// case if no imports in file
	if len(node.Imports) == 0 {
		// extract entity source code
		for _, decl := range node.Decls {
			switch decl.(type) {
			case *ast.GenDecl:
				genDecl := decl.(*ast.GenDecl)
				if genDecl.Tok == token.TYPE {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							// parse specific struct/interface with name that set in param `entityName`
							structType, stOk := typeSpec.Type.(*ast.StructType)
							interfaceType, interfaceOk := typeSpec.Type.(*ast.InterfaceType)
							if (stOk || interfaceOk) && typeSpec.Name.Name == entityName {
								var entityType ast.Expr
								if stOk {
									entityType = structType
								} else {
									entityType = interfaceType
								}
								// extract struct source code
								startLine := fset.Position(entityType.Pos()).Line
								endLine := fset.Position(entityType.End()).Line
								for _, structLine := range strings.Split(string(fromNodeBytes), "\n")[startLine-1 : endLine] {
									if !strings.HasPrefix(strings.TrimSpace(structLine), "//") && !strings.HasPrefix(strings.TrimSpace(structLine), "/*") {
										// checking for entities in the same package
										for _, pkgEntity := range pkgEntities {
											cleanStructLine := strings.TrimSpace(structLine)
											var fieldType string
											// process struct comments
											if strings.Contains(cleanStructLine, `\\`) {
												cleanStructLine = strings.Split(cleanStructLine, `\\`)[0]
											} else if strings.Contains(cleanStructLine, `\*`) {
												cleanStructLine = strings.Split(cleanStructLine, `\*`)[0]
											}
											fieldList := strings.Fields(cleanStructLine)
											// process tags in struct
											if strings.Contains(cleanStructLine, "`") {
												fieldType = fieldList[len(fieldList)-2]
											} else {
												fieldType = fieldList[len(fieldList)-1]
											}
											// process slice, map or selector type
											if strings.Contains(fieldType, "map[") {
												for _, mapType := range []string{strings.Split(fieldType, "]")[1], strings.Split(strings.Split(fieldType, "]")[1], "[")[1]} {
													if !strings.Contains(structLine, "{") && !strings.Contains(structLine, "}") && mapType == pkgEntity["name"] {
														parentInheritanceEntities = append(parentInheritanceEntities, pkgEntity)
													}
												}
											} else if strings.Contains(fieldType, "]") {
												if !strings.Contains(structLine, "{") && !strings.Contains(structLine, "}") && strings.Split(fieldType, "]")[1] == pkgEntity["name"] {
													parentInheritanceEntities = append(parentInheritanceEntities, pkgEntity)
												}
											} else {
												if !strings.Contains(structLine, "{") && !strings.Contains(structLine, "}") && fieldType == pkgEntity["name"] {
													parentInheritanceEntities = append(parentInheritanceEntities, pkgEntity)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	// extract all imports from file
	for _, imp := range node.Imports {
		impPath := imp.Path.Value[1 : len(imp.Path.Value)-1]
		// Extract struct source code
		for _, decl := range node.Decls {
			switch decl.(type) {
			case *ast.GenDecl:
				genDecl := decl.(*ast.GenDecl)
				if genDecl.Tok == token.TYPE {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							// parse specific struct/interface with name that set in param `entityName`
							structType, stOk := typeSpec.Type.(*ast.StructType)
							interfaceType, interfaceOk := typeSpec.Type.(*ast.InterfaceType)
							if (stOk || interfaceOk) && typeSpec.Name.Name == entityName {
								var entityType ast.Expr
								if stOk {
									entityType = structType
								} else {
									entityType = interfaceType
								}
								// extract struct source code
								startLine := fset.Position(entityType.Pos()).Line
								endLine := fset.Position(entityType.End()).Line
								regex := filepath.Base(impPath) + `\.\w+`
								parentInheritanceEntityRegex := regexp.MustCompile(regex)
								for _, structLine := range strings.Split(string(fromNodeBytes), "\n")[startLine-1 : endLine] {
									if !strings.HasPrefix(strings.TrimSpace(structLine), "//") && !strings.HasPrefix(strings.TrimSpace(structLine), "/*") {
										inhEntity := parentInheritanceEntityRegex.Find([]byte(structLine))
										if string(inhEntity) != "" {
											inhMap := map[string]string{"name": strings.Split(string(inhEntity), ".")[1], "path": impPath}
											parentInheritanceEntities = append(parentInheritanceEntities, inhMap)
										} else {
											// checking for entities in the same package
											for _, pkgEntity := range pkgEntities {
												cleanStructLine := strings.TrimSpace(structLine)
												var fieldType string
												// process struct comments
												if strings.Contains(cleanStructLine, `\\`) {
													cleanStructLine = strings.Split(cleanStructLine, `\\`)[0]
												} else if strings.Contains(cleanStructLine, `\*`) {
													cleanStructLine = strings.Split(cleanStructLine, `\*`)[0]
												}
												fieldList := strings.Fields(cleanStructLine)
												// process struct tags case
												if strings.Contains(cleanStructLine, "`") {
													fieldType = fieldList[len(fieldList)-2]
												} else {
													fieldType = fieldList[len(fieldList)-1]
												}
												// process slice, map or selector type
												if strings.Contains(fieldType, "map[") {
													for _, mapType := range []string{strings.Split(fieldType, "]")[1], strings.Split(strings.Split(fieldType, "]")[0], "[")[1]} {
														if !strings.Contains(structLine, "{") && !strings.Contains(structLine, "}") && mapType == pkgEntity["name"] {
															parentInheritanceEntities = append(parentInheritanceEntities, pkgEntity)
														}
													}
												} else if strings.Contains(fieldType, "]") {
													if !strings.Contains(structLine, "{") && !strings.Contains(structLine, "}") && strings.Split(fieldType, "]")[1] == pkgEntity["name"] {
														parentInheritanceEntities = append(parentInheritanceEntities, pkgEntity)
													}
												} else {
													if !strings.Contains(structLine, "{") && !strings.Contains(structLine, "}") && fieldType == pkgEntity["name"] {
														parentInheritanceEntities = append(parentInheritanceEntities, pkgEntity)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return parentInheritanceEntities
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
		// Parse the imports in Golang file
		if !ostool.Exists(path) {
			path = strings.TrimPrefix(strings.ReplaceAll(filepath.Join(filepath.Join(Parser.ProjectDir, ".."), filepath.Clean(path)), Parser.ProjectDir, ""), string(filepath.Separator))
			//fmt.Println("file: ", filePath, "import: ", path)
			if !ostool.Exists(path) {
				path = strings.TrimLeft(strings.ReplaceAll(filepath.Join(filepath.Join(Parser.ProjectDir, ".."), filepath.Clean(path)), filepath.Dir(Parser.ProjectDir), ""), string(filepath.Separator))
				fileDependenciesMap[filepath.Join("external_dependency", path)] = struct{}{}
				continue
			}
		}
		var packageGoFiles []string
		// iterating over all .go files in package
		files, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".go" {
				packageGoFiles = append(packageGoFiles, filepath.Join(path, file.Name()))
			}
		}
		for _, packageGoFile := range packageGoFiles {
			importedNode, err := parser.ParseFile(fset, packageGoFile, nil, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
			}
			// Iterate through the top-level declarations and find the structures
			for _, decl := range importedNode.Decls {
				switch decl.(type) {
				case *ast.GenDecl:
					genDecl := decl.(*ast.GenDecl)
					if genDecl.Tok == token.TYPE || genDecl.Tok == token.CONST {
						for _, spec := range genDecl.Specs {
							typeSpec := spec.(*ast.TypeSpec)
							lines := Parser.FilterComments(filePath)
							regex := `^.*` + filepath.Base(path) + `\.` + typeSpec.Name.Name + `.*$`
							imported := regexp.MustCompile(regex)
							for _, line := range lines {
								if imported.MatchString(line) {
									fileDependenciesMap[packageGoFile] = struct{}{}
								}
							}
						}
					}
				case *ast.FuncDecl:
					funcDecl := decl.(*ast.FuncDecl)
					lines := Parser.FilterComments(filePath)
					regex := `^.*` + filepath.Base(path) + `\.` + funcDecl.Name.Name + `.*$`
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
	for file, _ := range fileDependenciesMap {
		fileResults = append(fileResults, file)
	}
	return fileResults
}

func (Parser *GoParser) ExtractEntities(filePath string) []string {
	var entityResult []string
	if !ostool.Exists(filePath) {
		return []string{}
	}
	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	// Iterate through the top-level declarations and find the structures and functions
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

func (Parser *GoParser) ExtractExternalEntities(externalDependencyName, fromNodePath, fromNodeEntityName string) []string {
	var externalEntityDependencies []string
	file, err := os.Open(fromNodePath)
	if err != nil {
		log.Fatal(err)
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
	// Extract function/struct source code
	for _, decl := range node.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Recv != nil {
					switch fn.Recv.List[0].Type.(type) {
					case *ast.StarExpr:
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
								if importedObj.MatchString(structLine) && !strings.HasPrefix(strings.TrimSpace(structLine), "//") && !strings.HasPrefix(strings.TrimSpace(structLine), "/*") {
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
			log.Fatal(err)
		}
	}
	return node.Name.Name
}

func (Parser *GoParser) FilterComments(filename string) []string {
	var linesWithoutComments []string
	// Open file for reading
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Read file line by line and filter out comments
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !Parser.isComment(line) {
			linesWithoutComments = append(linesWithoutComments, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return linesWithoutComments
}

func (Parser *GoParser) isComment(line string) bool {
	line = Parser.stripWhitespace(line)
	if len(line) == 0 {
		return false
	}

	if line[0] == '/' && len(line) > 1 {
		if line[1] == '/' || line[1] == '*' {
			return true
		}
	}
	return false
}

func (Parser *GoParser) stripWhitespace(line string) string {
	var result []rune
	for _, ch := range line {
		if ch != ' ' && ch != '\t' {
			result = append(result, ch)
		}
	}
	return string(result)
}
