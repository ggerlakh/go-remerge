package parsers

type DependencyExtractor interface {
	ExtractDependency(filepath string) []string
}

type PackageExtractor interface {
	ExtractPackage(filepath string) string
}

type InheritanceExtractor interface {
	ExtractInheritance(filepath, entityName string) []string
}
