package testutils

import "fmt"

func GetErrorString(identifier, expected, got string) string {
	return fmt.Sprintf("Incorrect %s: expected %s, got %s", identifier, expected, got)
}
