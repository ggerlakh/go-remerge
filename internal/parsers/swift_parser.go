package parsers

type SwiftParser struct{}

func (parser *SwiftParser) ExtractInheritance(filepath, entityName string) []string {
	//TODO implement me
	panic("implement me")
}

func (parser *SwiftParser) ExtractDependencies(nodeName string) []string {
	return []string{}
}

func (parser *SwiftParser) ExtractEntities(filepath string) []string {
	return []string{}
}

func (parser *SwiftParser) ExtractPackage(filepath string) string {
	return ""
}
