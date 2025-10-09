package embed

import (
	"slices"
	"strings"
)

// embed start

func typeNameCollapseRecursive(input any) string {
	str, ok := input.(string)
	if ok {
		return str
	}
	array, ok := input.([]any)
	if !ok {
		panic("typeNameCollapseRecursive input is not an array")
	}
	if len(array) == 0 {
		panic("typeNameCollapseRecursive input is an empty array")
	}
	resultOr := []string{}
	for _, element := range array {
		resultOr = append(resultOr, typeNameCollapseRecursive(element))
	}
	slices.Sort(resultOr)
	resultOr = slices.Compact(resultOr)
	return strings.Join(resultOr, " | ")
}
