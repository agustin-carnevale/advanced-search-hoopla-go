package cli

import (
	"fmt"
	"strings"
)

func ValidateFlagEnum(value, flagName string, allowed ...string) error {
	if value == "" {
		return nil
	}

	for _, v := range allowed {
		if value == v {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid value for --%s: %q (allowed: %s)",
		flagName,
		value,
		strings.Join(allowed, ", "),
	)
}
