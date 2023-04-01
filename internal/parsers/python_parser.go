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

func (Parser *PythonParser) ExtractInheritance(entityFilePath, entityName string) []map[string]string {
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
	//TODO
	var entities []string
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v\n", err)
	}
	defer file.Close()
	// Create a regular expression to match class definitions
	classRegex := regexp.MustCompile(`class\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`)
	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Match the regular expression against the line
		match := classRegex.FindStringSubmatch(line)
		if match != nil {
			// The first match group is the class name
			className := match[1]
			entities = append(entities, className)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}
	return entities
}

func (Parser *PythonParser) ExtractExternalEntities(externalDependencyName, fromNodePath, fromNodeEntityName string) []string {
	//TODO
	var externalEntityDependencies []string
	// iterating over fromEntityName source code in .py file
	// Define regular expression to match class definition lines
	re := regexp.MustCompile(`class\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(?`)
	// Scan the file line by line and extract class definitions
	file, err := os.Open(fromNodePath)
	if err != nil {
		log.Fatalf("Error opening file: %s\n", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var className string
	var classLines []string
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			// Found a class definition
			// Start a new class
			className = match[1]
			classLines = []string{line}
		} else if className != "" && className == fromNodeEntityName {
			// filter line with comments and without python code
			if strings.HasPrefix(strings.TrimSpace(line), "#") || strings.HasPrefix(strings.TrimSpace(line), `"""`) {
				continue
			}
			// Add the line to the current class
			classBodyre := regexp.MustCompile(`^\s{4}.*`)
			emptyLine := regexp.MustCompile(`^\s*$`)
			if !strings.HasPrefix(line, "class") && !classBodyre.MatchString(line) && !emptyLine.MatchString(line) {
				break
			}
			classLines = append(classLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}
	// if fromEntity has external entity name in code line extract external entity
	regex := externalDependencyName + `\.\w+`
	importedExternalDependency := regexp.MustCompile(regex)
	for _, codeLine := range classLines {
		if !strings.HasPrefix(codeLine, "class") && (strings.Contains(codeLine, externalDependencyName)) {
			tmpRes := importedExternalDependency.Find([]byte(codeLine))
			if string(tmpRes) != "" {
				externalEntityDependencies = append(externalEntityDependencies, string(tmpRes))
			}
		}
	}
	return externalEntityDependencies
}

func (Parser *PythonParser) HasEntityDependency(fromEntityName, fromEntityPath, toEntityName, toEntityPackage string) bool {
	var hasEntityDependency bool
	// iterating over fromEntityName source code in .py file
	// Define regular expression to match class definition lines
	re := regexp.MustCompile(`class\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(?`)
	// Scan the file line by line and extract class definitions
	file, err := os.Open(fromEntityPath)
	if err != nil {
		log.Fatalf("Error opening file: %s\n", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var className string
	var classLines []string
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			// Found a class definition
			// Start a new class
			className = match[1]
			classLines = []string{line}
		} else if className != "" && className == fromEntityName {
			// filter line with comments and without python code
			if strings.HasPrefix(strings.TrimSpace(line), "#") || strings.HasPrefix(strings.TrimSpace(line), `"""`) {
				continue
			}
			// Add the line to the current class
			classBodyre := regexp.MustCompile(`^\s{4}.*`)
			emptyLine := regexp.MustCompile(`^\s*$`)
			if !strings.HasPrefix(line, "class") && !classBodyre.MatchString(line) && !emptyLine.MatchString(line) {
				break
			}
			classLines = append(classLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}
	// if fromEntity has toEntity name in code line hasEntityDependency = true
	for _, codeLine := range classLines {
		if !strings.HasPrefix(codeLine, "class") && (strings.Contains(codeLine, toEntityName) || strings.Contains(codeLine, toEntityName+"."+toEntityPackage)) {
			hasEntityDependency = true
		}
	}
	return hasEntityDependency
}

func (Parser *PythonParser) ExtractPackage(filePath string) string {
	var pyPackage string
	if strings.HasPrefix(filePath, "external_dependency") {
		pyPackage = strings.ReplaceAll(filePath, "external_dependency"+string(filepath.Separator), "")
	} else if ostool.Exists(filepath.Join(filepath.Dir(filePath), "__init__.py")) {
		pyPackage = filepath.Base(filepath.Dir(filePath))
	}
	return pyPackage
}
