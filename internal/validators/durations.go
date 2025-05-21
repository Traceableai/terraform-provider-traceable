package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var iso8601DurationRegex = regexp.MustCompile(`^P(T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?$`)

func ValidDurationFormat() validator.String {
	return durationValidator{}
}

type durationValidator struct{}

func (v durationValidator) Description(_ context.Context) string {
	return "Validates that the string is an ISO-8601 duration (e.g. PT1H30M)"
}

func (v durationValidator) MarkdownDescription(_ context.Context) string {
	return v.Description(context.Background())
}

func (v durationValidator) ValidateString(
	_ context.Context,
	req validator.StringRequest,
	resp *validator.StringResponse,
) {
	val := req.ConfigValue

	if val.IsUnknown() || val.IsNull() {
		return
	}

	if !iso8601DurationRegex.MatchString(val.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid ISO-8601 duration format",
			fmt.Sprintf("Expected a value like PT1H30M20S, but got: %q", val.ValueString()),
		)
	}
}
