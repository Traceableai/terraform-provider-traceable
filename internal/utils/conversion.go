package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertElementsSet[T any](elements types.Set, target *[]T) error {
	diags := elements.ElementsAs(context.Background(), target, false)
	if diags.HasError() {
		return fmt.Errorf("failed to convert elements: %v", diags)
	}
	return nil
}

func ConvertCustomStringPtrsToTerraformSet[T any](input []*T) (types.Set, error) {
	values := []attr.Value{}

	if input == nil {
		return types.SetNull(types.StringType), nil
	}

	for i, val := range input {
		if val == nil {
			values = append(values, types.StringNull())
			continue
		}

		// Use reflect to ensure *T is convertible to string
		valValue := reflect.ValueOf(val).Elem()
		if valValue.Kind() != reflect.String {
			return types.Set{}, fmt.Errorf("element at index %d is not a string-compatible type", i)
		}

		str := valValue.String()
		values = append(values, types.StringValue(str))
	}

	listValue, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {

		return types.SetNull(types.StringType), fmt.Errorf("conversion to type set failed ")
	}

	return listValue, nil

}

func ConvertStringPtrToTerraformSet(input []*string) (types.Set, error) {
	fmt.Println("entering the utils functions")

	if input == nil {
		return types.SetNull(types.StringType), nil
	}
	values := []attr.Value{}
	for _, val := range input {
		if val != nil {
			values = append(values, types.StringValue(*val))
		}
	}

	fmt.Println(values)

	setValue, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("converting String Pointer to Terraform Set fails ")
	}
	fmt.Println("existing the utils functions")
	return setValue, nil

}

func ConvertStringPtrSliceToTerraformList(input []*string) (types.List, error) {
	values := []attr.Value{}

	if input == nil {
		return types.ListNull(types.StringType), nil
	}

	for _, val := range input {
		if val != nil {
			values = append(values, types.StringValue(*val))
		}
	}

	listValue, diags := types.ListValueFrom(
		context.Background(),
		types.StringType,
		values,
	)

	if diags.HasError() {
		return types.ListNull(types.StringType), fmt.Errorf("converting String Pointer to Terraform List fails ")
	}

	return listValue, nil
}

func ConvertCustomStringPtrsToTerraformList[T any](input []*T) (types.List, error) {
	values := []attr.Value{}

	if input == nil {
		return types.ListNull(types.StringType), nil
	}

	for i, val := range input {
		if val == nil {
			values = append(values, types.StringNull())
			continue
		}

		// Use reflect to ensure *T is convertible to string
		valValue := reflect.ValueOf(val).Elem()
		if valValue.Kind() != reflect.String {
			return types.List{}, fmt.Errorf("element at index %d is not a string-compatible type", i)
		}

		str := valValue.String()
		values = append(values, types.StringValue(str))
	}

	listValue, diags := types.ListValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		// Handle or log error

		return types.ListNull(types.StringType), fmt.Errorf("conversion to type list failed ")
	}

	return listValue, nil

}

func ConvertSetToStrPointer(data types.Set) ([]*string, error) {
	values := []*string{}
	for _, elem := range data.Elements() {

		if elem, ok := elem.(types.String); ok {
			str := elem.ValueString()
			values = append(values, &str)
		} else {
			return nil, fmt.Errorf("Failed to convert %s to string pointer", elem)
		}
	}
	return values, nil
}
func ConvertDurationToSeconds(duration string) (string, error) {
	re := regexp.MustCompile(`P(T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?`)
	matches := re.FindStringSubmatch(duration)
	if len(matches) == 0 {
		return "", fmt.Errorf("Invalid duration format")
	}

	totalSeconds := 0

	// Parse hours, minutes, and seconds if present
	if matches[2] != "" { // Hours
		hours, _ := strconv.Atoi(matches[2])
		totalSeconds += hours * 3600
	}
	if matches[3] != "" { // Minutes
		minutes, _ := strconv.Atoi(matches[3])
		totalSeconds += minutes * 60
	}
	if matches[4] != "" { // Seconds
		seconds, _ := strconv.Atoi(matches[4])
		totalSeconds += seconds
	}

	return fmt.Sprintf("PT%dS", totalSeconds), nil
}

func GenerateRandomString(n int) string {
	bytes := make([]byte, n/2)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func EscapeString(input string) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		line = strings.ReplaceAll(line, `\`, `\\`)
		line = strings.ReplaceAll(line, `"`, `\"`)
		lines[i] = line
	}
	return strings.Join(lines, ` `)
}
