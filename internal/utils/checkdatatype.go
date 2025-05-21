package utils

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func HasValue(field interface{}) bool {

	if field == nil {
		return false
	}

	val := reflect.ValueOf(field)
	if (val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface) && val.IsNil() {
		return false
	}

	switch v := field.(type) {
	case attr.Value:
		if v.IsNull() || v.IsUnknown() {
			return false
		}

		// Handle specific framework types
		switch concrete := v.(type) {
		case types.String:
			return concrete.ValueString() != ""
		case types.List, types.Set:
			elements := reflect.ValueOf(concrete).MethodByName("Elements").Call(nil)[0]
			for i := 0; i < elements.Len(); i++ {
				element := elements.Index(i).Interface().(attr.Value)
				if HasValue(element) {
					return true
				}
			}
			return false
		case types.Map:
			for _, v := range concrete.Elements() {
				if HasValue(v) {
					return true
				}
			}
			return false

		case types.Object:
			attrs := concrete.Attributes()
			if len(attrs) == 0 {
				return false
			}
			for _, attrVal := range attrs {
				if HasValue(attrVal) {
					return true
				}
			}
			return false
		default:
			return true
		}
	default:
		return true
	}

}
