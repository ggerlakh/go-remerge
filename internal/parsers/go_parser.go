package parsers

type GoParser struct{}

func (parser *GoParser) ExtractDependencies(filepath string) []string {
	return []string{}
}

func (parser *GoParser) ExtractPackage(filepath string) string {
	return ""
}
