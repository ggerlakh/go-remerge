package parsers

type KotlinParser struct{}

func (parser *KotlinParser) ExtractDependencies(filepath string) []string {
	return []string{}
}

func (parser *KotlinParser) ExtractPackage(filepath string) string {
	return ""
}
