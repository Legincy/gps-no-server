package validation

import (
	"context"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Error())
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	errMessages := make([]string, len(v))
	for i, err := range v {
		errMessages[i] = err.Error()
	}

	return strings.Join(errMessages, "; ")
}

type Validator interface {
	Validate(ctx context.Context) error
}
