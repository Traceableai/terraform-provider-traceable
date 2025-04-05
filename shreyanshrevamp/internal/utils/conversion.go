package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
)

func ConvertStringPtrSliceToTerraformList(values []*string) types.List {
	ans := []attr.Value{}
	for _, val := range values {
		ans = append(ans, types.StringValue(*val))
	}
	return types.ListValueMust(types.StringType, ans)
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

	return types.ListValueMust(types.StringType, values), nil
}
