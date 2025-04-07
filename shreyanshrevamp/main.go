package main

import "fmt"

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/traceableai/terraform-provider-traceable/shreyanshrevamp/internal/provider"
)

var (
	//comes from goreleaser
	version string = "dev"
)

func main() {

	var debug bool
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/traceableai/traceable",
		Debug:   debug,
	}
	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(version)

}

// package utils

// import (
// 	"context"
// 	"fmt"
// 	"reflect"

// 	"github.com/hashicorp/terraform-plugin-framework/attr"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// )

// func ConvertStringPtrSliceToTerraformList(values []*string) types.List {
// 	ans := []attr.Value{}
// 	for _, val := range values {
// 		if val == nil {
// 			continue
// 		}
// 		ans = append(ans, types.StringValue(*val))

// 	}
// 	listValue, diags := types.ListValueFrom(context.Background(), types.StringType, ans)
// 	if diags.HasError() {
// 		// Handle or log error
// 		return types.ListNull(types.StringType)
// 	}
// 	if listValue.ElementType(context.Background()).Equal(types.StringType) {
// 		fmt.Printf("hello I am inside")
// 		return types.ListNull(types.StringType)
// 	}

// 	return listValue
// }

// // func ConvertCustomStringPtrsToTerraformList[T ~string](input []*T) types.List {
// // 	result := make([]*string, len(input))
// // 	for i, val := range input {
// // 		if val != nil {
// // 			str := string(*val)
// // 			result[i] = &str
// // 		}
// // 	}
// // 	ans := []attr.Value{}
// // 	for _, val := range result {
// // 		ans = append(ans, types.StringValue(*val))
// // 	}
// // 	return types.ListValueMust(types.StringType, ans)

// // }
// //

// func ConvertCustomStringPtrsToTerraformList[T any](input []*T) (any, error) {
// 	values := []attr.Value{}

// 	for i, val := range input {
// 		if val == nil {
// 			values = append(values, types.StringNull())
// 			continue
// 		}

// 		// Use reflect to ensure *T is convertible to string
// 		valValue := reflect.ValueOf(val).Elem()
// 		if valValue.Kind() != reflect.String {
// 			return types.List{}, fmt.Errorf("element at index %d is not a string-compatible type", i)
// 		}

// 		str := valValue.String()
// 		values = append(values, types.StringValue(str))
// 	}

// 	listValue, diags := types.ListValueFrom(context.Background(), types.StringType, values)
// 	if diags.HasError() {
// 		// Handle or log error

// 		return types.ListNull(types.StringType), fmt.Errorf("conversion to type list failed ")
// 	}

// 	return listValue, nil

// }

// // Generic function to convert []*T to a Terraform list of string values
// // func ConvertCustomStringToTerraformList[T any](input []*T) (types.List, error) {
// // 	values := make([]attr.Value, 0, len(input))

// // 	for _, val := range input {
// // 		if val == nil {
// // 			values = append(values, types.StringNull())
// // 			continue
// // 		}

// // 		// Convert to interface{} to use type assertion
// // 		switch v := any(*val).(type) {
// // 		case string:
// // 			values = append(values, types.StringValue(v))
// // 		case fmt.Stringer:
// // 			values = append(values, types.StringValue(v.String()))
// // 		default:
// // 			return types.List{}, fmt.Errorf("unsupported type: %T", v)
// // 		}
// // 	}

// //		return types.ListValueMust(types.StringType, values), nil
// //	}
// func isStringList(listVal types.List) bool {
// 	ctx := context.Background()

// 	return listVal.ElementType(ctx).Equal(types.StringType)
// }
// func main() {

// 	dynamicList, _ := types.ListValue(types.DynamicType, []attr.Value{
// 		types.StringValue("a"),
// 	})
// 	a := "shreyansh"
// 	b := "gupta"
// 	arr := []*string{&a, &b}
// 	list := ConvertStringPtrSliceToTerraformList(arr)
// 	isList := isStringList(list)
// 	isDynamicList := isStringList(dynamicList)
// 	fmt.Println(isList)
// 	fmt.Println(isDynamicList)
// }
