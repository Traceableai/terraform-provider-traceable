package utils

import (
	"context"
	"fmt"
	"reflect"

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
	if input == nil {
		return types.SetNull(types.StringType), nil
	}
	values := []attr.Value{}
	for _, val := range input {
		if val != nil {
			values = append(values, types.StringValue(*val))
		}
	}

	setValue, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("converting String Pointer to Terraform Set fails ")
	}
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
