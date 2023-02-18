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

func (parser *PythonParser) ExtractInheritance(filepath, entityName string) []string {
	//TODO implement me
	panic("implement me")
}

func (parser *PythonParser) ExtractDependencies(path string) []string {
	// TODO use *Node instead of string and implement extracting packages there
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var dependencies []string
	var validImport = regexp.MustCompile(`(?m)^(?:from[ ]+(\S+)[ ]+)?import[ ]+(\S+)(?:[ ]+as[ ]+\S+)?[ ]*$`)
	currDir := filepath.Dir(path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file %s: %v\n", path, err)
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
						log.Println("Dependency not exist: ", cleanImportPath, "path: ", path)
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
		log.Fatalf("Error while reading file %s: %v\n", path, err)
	}
	return dependencies
}

func (parser *PythonParser) ExtractEntities(filepath string) []string {
	return []string{}
}

func (parser *PythonParser) ExtractPackage(filepath string) string {
	// TODO
	return ""
}
