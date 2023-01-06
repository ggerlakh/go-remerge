package parsers

type SwiftParser struct{}

func (parser *SwiftParser) ExtractDependencies(filepath string) []string {
	return []string{}
}

func (parser *SwiftParser) ExtractPackage(filepath string) string {
	return ""
}
