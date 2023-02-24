package parsers

import (
	"bufio"
	"go-remerge/tools/ostool"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PythonParser struct{}

func (Parser *PythonParser) ExtractInheritance(filePath, entityName string) []string {
	//TODO implement me
	panic("implement me")
}

func (Parser *PythonParser) ExtractDependencies(filePath string) []string {
	// TODO use *Node instead of string and implement extracting packages there
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var dependencies []string
	var validImport = regexp.MustCompile(`(?m)^(?:from[ ]+(\S+)[ ]+)?import[ ]+(\S+)(?:[ ]+as[ ]+\S+)?[ ]*$`)
	currDir := filepath.Dir(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %v\n", filePath, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if validImport.MatchString(line) {
			importPath := strings.Fields(line)[1]
			// case "import abcd.efg"
			if strings.Contains(importPath, ".") && importPath != "." && importPath != ".." {
				// case "import ..abcd.efg"
				if strings.HasPrefix(importPath, "..") {
					var dependency string
					cleanImportPath := filepath.Clean(filepath.Join("..", strings.Replace(strings.Split(importPath, "..")[1], ".", string(filepath.Separator), -1)))
					if ostool.Exists(cleanImportPath + ".py") {
						dependency = cleanImportPath + ".py"
					} else if ostool.Exists(cleanImportPath) {
						dependency = filepath.Join(cleanImportPath, "__init__.py")
					} else if ostool.Exists(filepath.Clean(filepath.Join(currDir, cleanImportPath))) {
						dependency = filepath.Join(filepath.Clean(filepath.Join(currDir, cleanImportPath)), "__init__.py")
					} else if ostool.Exists(filepath.Clean(filepath.Join(currDir, cleanImportPath)) + ".py") {
						dependency = filepath.Clean(filepath.Join(currDir, cleanImportPath)) + ".py"
					} else {
						log.Println("Dependency not exist: ", cleanImportPath, "path: ", filePath)
						dependency = filepath.Join("external_dependency", filepath.Base(cleanImportPath))
					}
					dependencies = append(dependencies, dependency)
				} else {
					var dependency string
					cleanImportPath := strings.Replace(importPath, ".", string(filepath.Separator), -1)
					if ostool.Exists(cleanImportPath + ".py") {
						dependency = cleanImportPath + ".py"
					} else if ostool.Exists(cleanImportPath) {
						dependency = filepath.Join(cleanImportPath, "__init__.py")
					} else if ostool.Exists(filepath.Clean(filepath.Join(currDir, cleanImportPath))) {
						dependency = filepath.Join(filepath.Clean(filepath.Join(currDir, cleanImportPath)), "__init__.py")
					} else if ostool.Exists(filepath.Clean(filepath.Join(currDir, cleanImportPath)) + ".py") {
						dependency = filepath.Clean(filepath.Join(currDir, cleanImportPath)) + ".py"
					} else {
						dependency = filepath.Join("external_dependency", filepath.Base(cleanImportPath))
					}
					dependencies = append(dependencies, dependency)
				}
			} else if importPath == "." || importPath == ".." { // case .. and .
				dependency := filepath.Clean(filepath.Join(currDir, importPath, "__init__.py"))
				dependencies = append(dependencies, dependency)
			} else { // case "import module"
				var dependency string
				if ostool.Exists(importPath + ".py") {
					dependency = importPath + ".py"
				} else if ostool.Exists(importPath) {
					dependency = filepath.Join(importPath, "__init__.py")
				} else {
					dependency = filepath.Join("external_dependency", importPath)
				}
				dependencies = append(dependencies, dependency)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error while reading file %s: %v\n", filePath, err)
	}
	return dependencies
}

func (Parser *PythonParser) ExtractEntities(filePath string) []string {
	return []string{}
}

func (Parser *PythonParser) ExtractExternalEntities(externalDependencyName, fromNodePath string) []string {
	var externalEntityDependencies []string
	return externalEntityDependencies
}

func (Parser *PythonParser) HasEntityDependency(fromEntityName, fromEntityPath, toEntityName, toEntityPath string) bool {
	var hasEntityDependency bool
	// TODO
	return hasEntityDependency
}

func (Parser *PythonParser) ExtractPackage(filePath string) string {
	// TODO
	return ""
}
