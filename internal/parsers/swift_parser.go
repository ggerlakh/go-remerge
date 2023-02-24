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

func (Parser *SwiftParser) ExtractExternalEntities(externalDependencyName, fromNodePath string) []string {
	var externalEntityDependencies []string
	return externalEntityDependencies
}

func (Parser *SwiftParser) HasEntityDependency(fromEntityName, fromEntityPath, toEntityName, toEntityPath string) bool {
	var hasEntityDependency bool
	// TODO
	return hasEntityDependency
}

func (Parser *SwiftParser) ExtractPackage(filePath string) string {
	return ""
}
