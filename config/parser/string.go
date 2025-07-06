package parser

func ParseString(data any) (string, error) {
	if data == nil {
		return "", nilErr([]string{"string"})
	}
	s, ok := data.(string)
	if !ok {
		return "", typeErr([]string{"string"}, data)
	}

	return s, nil
}

func ParseStringArray(data any) ([]string, error) {
	if data == nil {
		return nil, nilErr([]string{"string", "[]string"})
	}

	switch d := data.(type) {
	case string:
		return []string{d}, nil
	case []any:
		return parseStringArrayFromAnyArray(d)
	default:
		return nil, typeErr([]string{"string", "[]string"}, data)
	}
}

func parseStringArrayFromAnyArray(a []any) ([]string, error) {
	res := make([]string, len(a))

	for idx, d := range a {
		s, ok := d.(string)
		if !ok {
			return nil, typeErr([]string{"string"}, d)
		}

		res[idx] = s
	}

	return res, nil
}
