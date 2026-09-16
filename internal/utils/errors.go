package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type InvalidFieldError struct {
	Field  string
	Reason string
}

// Error implements the error interface for InvalidFieldError.
func (e *InvalidFieldError) Error() string {
	return fmt.Sprintf("invalid input for field '%s': %s", e.Field, e.Reason)
}

func NewInvalidError(field string, reason string) error {

	return &InvalidFieldError{
		Field:  field,
		Reason: reason,
	}

}

func AddError(ctx context.Context, resp *diag.Diagnostics, err error) {
	if err == nil {
		return
	}

	if gqlErr, ok := err.(*graphql.HTTPError); ok {
		switch gqlErr.StatusCode {
		case 401:
			resp.AddError("Unauthorized", "Please check your token and try again")
			return
		case 403:
			resp.AddError("Unauthorized", "Please check your token and try again")
			return

		case 502:
			resp.AddError("Bad Gateway Error", "There is problem in connecting to platform ,please try again after some time ")
			return

		default:
			resp.AddError("API error", gqlErr.Error())
			return

		}
	}

	if invalidFieldErr, ok := err.(*InvalidFieldError); ok {
		resp.AddError("Invalid input", fmt.Sprintf("Invalid input for field '%s': %s", invalidFieldErr.Field, invalidFieldErr.Reason))
		return
	}

	var gqlErrs gqlerror.List
	if errors.As(err, &gqlErrs) {
		resp.AddError("API error", formatGqlErrors(gqlErrs))
		return
	}

	resp.AddError("API error", err.Error())

	return

}

// formatGqlErrors renders each GraphQL error's message plus any extensions
// (e.g. classification, code) the API attached, since gqlerror.List.Error()
// discards Extensions entirely.
func formatGqlErrors(errs gqlerror.List) string {
	var lines []string
	for _, e := range errs {
		line := e.Message
		if len(e.Extensions) > 0 {
			var extParts []string
			for k, v := range e.Extensions {
				extParts = append(extParts, fmt.Sprintf("%s=%v", k, v))
			}
			line = fmt.Sprintf("%s (%s)", line, strings.Join(extParts, ", "))
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
