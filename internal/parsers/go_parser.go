package parsers

type GoParser struct{}

func (parser *GoParser) ExtractInheritance(filepath, entityName string) []string {
	//TODO implement me
	panic("implement me")
}

func (parser *GoParser) ExtractDependencies(nodeName string) []string {
	return []string{}
}

func (parser *GoParser) ExtractEntities(filepath string) []string {
	return []string{}
}

func (parser *GoParser) ExtractPackage(filepath string) string {
	return ""
}
