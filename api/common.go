package api

import (
	"fmt"
)

// wrap wrap an error
func wrap(msg string, err error) error {
	return fmt.Errorf("%s: %v", msg, err)
}
