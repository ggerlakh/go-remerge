package parsers

type PythonParser struct{}

func (parser *PythonParser) ExtractDependencies(filepath string) []string {
	return []string{}
}

func (parser *PythonParser) ExtractPackage(filepath string) string {
	return ""
}
