package parsers

type KotlinParser struct{}

func (parser *KotlinParser) ExtractInheritance(filepath, entityName string) []string {
	//TODO implement me
	panic("implement me")
}

func (parser *KotlinParser) ExtractDependencies(nodeName string) []string {
	return []string{}
}

func (parser *KotlinParser) ExtractEntities(filepath string) []string {
	return []string{}
}

func (parser *KotlinParser) ExtractPackage(filepath string) string {
	return ""
}
