package parser

func ParseMap(data any) (map[string]any, error) {
	if data == nil {
		return nil, nilErr([]string{"map[string]any"})
	}

	m, ok := data.(map[string]any)
	if !ok {
		return nil, typeErr([]string{"map[string]any"}, data)
	}

	return m, nil
}

func ParseMapArray(data any) ([]map[string]any, error) {
	if data == nil {
		return nil, nilErr([]string{"[]map[string]any"})
	}

	res, ok := data.([]map[string]any)
	if !ok {
		return nil, typeErr([]string{"[]map[string]any"}, data)
	}

	return res, nil
}
