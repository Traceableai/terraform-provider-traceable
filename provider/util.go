
package provider

import(
	"strings"
	"fmt"
)


func listToString(stringArray []string) string{
	var formattedStrings []string
	for _, s := range stringArray {
		formattedStrings = append(formattedStrings, fmt.Sprintf(`"%s"`, s))
	}
	return strings.Join(formattedStrings, ", ")
}

func toStringSlice(interfaceSlice []interface{}) []string {
	stringSlice := make([]string, len(interfaceSlice))
	for i, v := range interfaceSlice {
		stringSlice[i] = v.(string)
	}
	return stringSlice
}

func convertToStringSlice(data []interface{}) []interface{} {
    var result []interface{}
    for _, v := range data {
        result = append(result, v.(interface{}))
    }
    return result
}