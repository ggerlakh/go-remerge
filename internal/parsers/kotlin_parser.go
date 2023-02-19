package parsers

type KotlinParser struct{}

func (Parser *KotlinParser) ExtractInheritance(filePath, entityName string) []string {
	//TODO implement me
	panic("implement me")
}

func (Parser *KotlinParser) ExtractDependencies(filePath string) []string {
	return []string{}
}

func (Parser *KotlinParser) ExtractEntities(filePath string) []string {
	return []string{}
}

func (Parser *KotlinParser) ExtractPackage(filePath string) string {
	return ""
}
