package parse

import (
	"os/exec"

	"github.com/BurntSushi/toml"
)

type Command struct {
	*exec.Cmd
	raw any
}

func (c *Command) Raw() any {
	return c.raw
}

func (c Command) MarshalTOML() ([]byte, error) {
	return toml.Marshal(c.raw)
}

func (c *Command) UnmarshalTOML(data any) error {
	if data == nil {
		return nilErr("Command")
	}

	var cmds []string
	var err error

	switch d := data.(type) {
	case string:
		cmds, err = parseCommandFromString(d)
	case []any:
		cmds, err = parseStringArrayFromAnyArray(d)
	case []string:
		cmds = d
	default:
		return typeErr("Command", d)
	}

	if err != nil {
		return err
	}

	c.raw = data
	c.Cmd = exec.Command(cmds[0], cmds[1:]...)

	return nil
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
