package utils

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

		case 502:
			resp.AddError("Bad Gateway Error", "There is problem in connecting to platform ,please try again after some time ")
			return

		default:
			resp.AddError("Internal provider error", "An unexpected error occurred. Please contact support or upgrade to the latest version.")
			return

		}
	}

	if invalidFieldErr, ok := err.(*InvalidFieldError); ok {
		resp.AddError("Invalid input", fmt.Sprintf("Invalid input for field '%s': %s", invalidFieldErr.Field, invalidFieldErr.Reason))
		return
	}

	// resp.AddError("Internal provider error", "An unexpected error occurred. Please contact support or upgrade to the latest version.")
	// return

	resp.AddError("Error aa gaya", err.Error())
	return

}
