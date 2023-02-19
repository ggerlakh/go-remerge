package parsers

type SwiftParser struct{}

func (Parser *SwiftParser) ExtractInheritance(filePath, entityName string) []string {
	//TODO implement me
	panic("implement me")
}

func (Parser *SwiftParser) ExtractDependencies(filePath string) []string {
	return []string{}
}

func (Parser *SwiftParser) ExtractEntities(filePath string) []string {
	return []string{}
}

func (Parser *SwiftParser) ExtractPackage(filePath string) string {
	return ""
}
