package parsers

type DependencyExtractor interface {
	// ExtractDependencies extracts imports dependencies from file
	ExtractDependencies(nodeName string) []string
	ExtractEntities(filepath string) []string
	ExtractPackage(filepath string) string
	ExtractExternalEntities(externalDependencyName, fromNodePath string) []string
	HasEntityDependency(fromEntityName, fromEntityPath, toEntityName, toEntityPath string) bool
}

type InheritanceExtractor interface {
	// ExtractInheritance extracts inheritance entities for entity from given file
	ExtractInheritance(entityFilePath, entityName string) []map[string]string
	ExtractPackage(filepath string) string
}

type CompleteGraphExtractor interface {
	DependencyExtractor
	InheritanceExtractor
}
