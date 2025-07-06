package parser

func ParseCommand(data any) ([]string, error) {
	if data == nil {
		return nil, nilErr([]string{"string", "[]string"})
	}

	switch d := data.(type) {
	case string:
		return parseCommandFromString(d)
	case []any:
		return parseStringArrayFromAnyArray(d)
	default:
		return nil, typeErr([]string{"string", "[]string"}, d)
	}
}

func parseCommandFromString(s string) ([]string, error) {
	var (
		res     []string
		acc     string
		quoted  string
		escaped bool
		tpl     bool
	)

	stream := newTokenStream(s)

	for stream.Next() {
		c := stream.Take()

		// If we find the escape char and we're not currently escaped then escape the
		// next character
		if c == "\\" && !escaped {
			escaped = true
			continue
		}

		// Escape the next character and reset escaped state
		if escaped {
			acc += c
			escaped = false
			continue
		}

		// If we find a quote char and we're in a quoted state with the matching
		// character then push the acc and reset state
		if (c == "'" && quoted == "'") || (c == "\"" && quoted == "\"") {
			res = append(res, acc)
			quoted = ""
			acc = ""

			continue
		}

		// If we find a quote char and we're not in a quoted state then set state
		// to be quoted and push acc if it's not empty
		if (c == "'" && quoted == "") || (c == "\"" && quoted == "") {
			quoted = c

			if acc != "" {
				res = append(res, acc)
				acc = ""
			}

			continue
		}

		next := stream.Peek()

		// If we have {{ and are not in a template context, set template context
		if c == "{" && next == "{" && !tpl {
			stream.Take()
			tpl = true

			acc += "{{"

			continue
		}

		// If we have }} and are in a template context, unset template context
		if c == "}" && next == "}" && tpl {
			stream.Take()
			tpl = false

			acc += "}}"

			continue
		}

		// If we find a space and we're not quoted and not in template context then
		// push acc and reset
		if c == " " && quoted == "" && !tpl {
			if acc != "" {
				res = append(res, acc)
				acc = ""
			}

			continue
		}

		acc += c
	}

	if acc != "" {
		res = append(res, acc)
	}

	return res, nil
}
