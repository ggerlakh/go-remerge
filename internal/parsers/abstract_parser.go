package parsers

type DependencyExtractor interface {
	// ExtractDependencies extracts imports dependencies from file
	ExtractDependencies(nodeName string) []string
	ExtractEntities(filepath string) []string
	ExtractPackage(filepath string) string
}

type InheritanceExtractor interface {
	// ExtractInheritance extracts inheritance entities for entity from given file
	ExtractInheritance(filepath, entityName string) []string
	ExtractEntities(filepath string) []string
}

type CompleteGraphExtractor interface {
	DependencyExtractor
	InheritanceExtractor
}
