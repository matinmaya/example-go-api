package rediservice

import (
	"fmt"
)

// private func
func validateNotEmpty(params map[string]string) error {
	for name, value := range params {
		if value == "" {
			return fmt.Errorf("argument [%s] cannot be empty", name)
		}
	}
	return nil
}
