package parser

import (
	"fmt"
	"strings"
)

func nilErr(expectedTypes []string) error {
	return fmt.Errorf("expected %s, found nil", expectedType(expectedTypes))
}

func typeErr(expectedTypes []string, found any) error {
	return fmt.Errorf("expected %s, found %T", expectedType(expectedTypes), found)
}

func expectedType(expectedTypes []string) string {
	if len(expectedTypes) == 1 {
		return fmt.Sprintf("type %s", expectedTypes[0])
	}

	return fmt.Sprintf("one of types [%s]", strings.Join(expectedTypes, ", "))
}
