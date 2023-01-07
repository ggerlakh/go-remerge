package parsers

type PythonParser struct{}

func (parser *PythonParser) ExtractInheritance(filepath, entityName string) []string {
	//TODO implement me
	panic("implement me")
}

func (parser *PythonParser) ExtractDependencies(nodeName string) []string {
	return []string{}
}

func (parser *PythonParser) ExtractEntities(filepath string) []string {
	return []string{}
}

func (parser *PythonParser) ExtractPackage(filepath string) string {
	return ""
}
