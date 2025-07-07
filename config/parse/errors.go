package parse

import (
	"fmt"
)

func nilErr(expectedType string) error {
	return fmt.Errorf("cannot unmarshal nil to type %s", expectedType)
}

func typeErr(expectedType string, found any) error {
	return fmt.Errorf("cannot unmarshal type %T to type %s", found, expectedType)
}
