package models

type Mapper map[string]Value

func ParseMapper(data map[string]interface{}) Mapper {
	mapper := make(Mapper)
	for key, content := range data {
		value := ParseValue(content)
		mapper[key] = value
	}
	return mapper
}
