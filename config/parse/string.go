package parse

func parseStringArrayFromAnyArray(data []any) ([]string, error) {
	res := make([]string, len(data))

	for idx, d := range data {
		s, ok := d.(string)
		if !ok {
			return nil, typeErr("string", d)
		}

		res[idx] = s
	}

	return res, nil
}
