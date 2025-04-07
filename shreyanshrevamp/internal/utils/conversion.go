package utils

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertStringPtrSliceToTerraformList(input []*string) (types.List, error) {
	values := []attr.Value{}

	if input != nil {

		for _, val := range input {
			if val != nil {
				values = append(values, types.StringValue(*val))
			}
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

// func ConvertCustomStringPtrsToTerraformList[T ~string](input []*T) types.List {
// 	result := make([]*string, len(input))
// 	for i, val := range input {
// 		if val != nil {
// 			str := string(*val)
// 			result[i] = &str
// 		}
// 	}
// 	ans := []attr.Value{}
// 	for _, val := range result {
// 		ans = append(ans, types.StringValue(*val))
// 	}
// 	return types.ListValueMust(types.StringType, ans)

// }
//

func ConvertCustomStringPtrsToTerraformList[T any](input []*T) (types.List, error) {
	values := []attr.Value{}

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

// Generic function to convert []*T to a Terraform list of string values
// func ConvertCustomStringToTerraformList[T any](input []*T) (types.List, error) {
// 	values := make([]attr.Value, 0, len(input))

// 	for _, val := range input {
// 		if val == nil {
// 			values = append(values, types.StringNull())
// 			continue
// 		}

// 		// Convert to interface{} to use type assertion
// 		switch v := any(*val).(type) {
// 		case string:
// 			values = append(values, types.StringValue(v))
// 		case fmt.Stringer:
// 			values = append(values, types.StringValue(v.String()))
// 		default:
// 			return types.List{}, fmt.Errorf("unsupported type: %T", v)
// 		}
// 	}

// 	return types.ListValueMust(types.StringType, values), nil
// }
