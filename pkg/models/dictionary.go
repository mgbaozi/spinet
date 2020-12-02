package models

func ProcessAppDictionary(app string, dictionray map[string]interface{}, ctx *Context) {
	data := ctx.AppData[app]
	for key, content := range dictionray {
		value := ParseValue(content)
		if v, err := value.Extract(ctx.Dictionary, data); err == nil {
			ctx.Dictionary[key] = v
		}
	}
}
