package parsers

type DependencyExtractor interface {
	// ExtractDependencies ExtractDependency extracts imports dependencies from file
	ExtractDependencies(filepath string) []string
	ExtractPackage(filepath string) string
}

type InheritanceExtractor interface {
	// ExtractInheritance extracts inheritance entities for entity from given file
	ExtractInheritance(filepath, entityName string) []string
}
