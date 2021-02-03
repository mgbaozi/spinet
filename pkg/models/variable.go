package models

func buildInVariables(override map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	vars := GetBuildInVariables()
	for name, item := range vars {
		res[name] = item.New(nil).Data()
	}
	if override != nil {
		for name, item := range override {
			res[name] = item
		}
	}
	return res
}
