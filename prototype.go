package main

import (
	"fmt"
	"strings"
)

// ARRAY PROTOTYPE FUNCTIONS ---
func arrayLength(a *ArrayValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	result := MakeNumber(float64(len(a.Elements)))
	return result, nil
}

func arrayPush(a *ArrayValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("array.push requires at least one argument")
	}
	a.Elements = append(a.Elements, args...)
	result := MakeNumber(float64(len(a.Elements)))
	return result, nil
}

func arrayPop(a *ArrayValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	if len(a.Elements) == 0 {
		return nil, fmt.Errorf("array.pop called on an empty array")
	}
	lastIndex := len(a.Elements) - 1
	poppedElement := a.Elements[lastIndex]
	a.Elements = a.Elements[:lastIndex]
	return poppedElement, nil
}

func arrayJoin(a *ArrayValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("array.join requires exactly one argument")
	}
	separator, ok := args[0].(*StringValue)
	if !ok {
		return nil, fmt.Errorf("array.join argument must be a string")
	}
	var parts []string
	for _, elem := range a.Elements {
		if strElem, ok := elem.(*StringValue); ok {
			parts = append(parts, strElem.Value)
		} else {
			return nil, fmt.Errorf("array.join elements must be strings")
		}
	}
	result := MakeString(strings.Join(parts, separator.Value))
	return result, nil
}

//
// func arrayFilter(a *ArrayValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
// 	if len(args) != 1 {
// 		return nil, fmt.Errorf("array.filter requires exactly one argument")
// 	}
// 	filterFunc, ok := args[0].(*FunctionValue)
// 	if !ok {
// 		return nil, fmt.Errorf("array.filter argument must be a function")
// 	}
//
// 	filteredElements := []RuntimeValue{}
// 	for _, elem := range a.Elements {
// 		result, err := evaluateCallExpression(&CallExpr{
// 			Caller: &Identifier{Value: filterFunc.Name},
// 			Args:   []Expression{elem.(Expression)},
// 		}, env)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if boolResult, ok := result.(*BooleanValue); ok && boolResult.Value {
// 			filteredElements = append(filteredElements, elem)
// 		}
// 	}
//
// 	result := MakeArray(filteredElements)
// 	return result, nil
// }
//
// func arrayMap(a *ArrayValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
// 	if len(args) != 1 {
// 		return nil, fmt.Errorf("array.map requires exactly one argument")
// 	}
// 	mapFunc, ok := args[0].(*FunctionValue)
// 	if !ok {
// 		return nil, fmt.Errorf("array.map argument must be a function")
// 	}
//
// 	mappedElements := []RuntimeValue{}
// 	for _, elem := range a.Elements {
// 		result, err := evaluateCallExpression(&CallExpr{
// 			Caller: &Identifier{Value: mapFunc.Name},
// 			Args:   []Expression{elem.(Expression)},
// 		}, env)
// 		if err != nil {
// 			return nil, err
// 		}
// 		mappedElements = append(mappedElements, result)
// 	}
//
// 	result := MakeArray(mappedElements)
// 	return result, nil
// }
//
// func arrayFind(a *ArrayValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
// 	if len(args) != 1 {
// 		return nil, fmt.Errorf("array.find requires exactly one argument")
// 	}
// 	findFunc, ok := args[0].(*FunctionValue)
// 	if !ok {
// 		return nil, fmt.Errorf("array.find argument must be a function")
// 	}
//
// 	for _, elem := range a.Elements {
// 		result, err := evaluateCallExpression(&CallExpr{
// 			Caller: &Identifier{Value: findFunc.Name},
// 			Args:   []Expression{elem.(Expression)},
// 		}, env)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if boolResult, ok := result.(*BooleanValue); ok && boolResult.Value {
// 			return elem, nil
// 		}
// 	}
//
// 	return MakeNull(), nil // Return null if no element matches
// }

func arrayIncludes(a *ArrayValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("array.includes requires exactly one argument")
	}
	element := args[0]
	found := false
	for _, elem := range a.Elements {
		if elem.Type() == element.Type() && elem.String() == element.String() {
			found = true
			break
		}
	}
	return MakeBool(found), nil
}

// STRING PROTOTYPE FUNCTIONS ---

func stringLength(s *StringValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	result := MakeNumber(float64(len(s.Value)))
	return result, nil
}

func stringToUpperCase(s *StringValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	result := MakeString(strings.ToUpper(s.Value))
	return result, nil
}

func stringToLowerCase(s *StringValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	result := MakeString(strings.ToLower(s.Value))
	return result, nil
}

func stringCharAt(s *StringValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string.charAt requires exactly one argument")
	}
	index, ok := args[0].(*NumberValue)
	if !ok {
		return nil, fmt.Errorf("string.charAt argument must be a number")
	}
	if index.Value < 0 || int(index.Value) >= len(s.Value) {
		return MakeString(""), nil // Return empty string for out of bounds
	}
	result := MakeString(string(s.Value[int(index.Value)]))
	return result, nil
}

func stringSubstring(s *StringValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("string.substring requires one or two arguments")
	}
	start, ok := args[0].(*NumberValue)
	if !ok {
		return nil, fmt.Errorf("string.substring first argument must be a number")
	}
	end := len(s.Value)
	if len(args) == 2 {
		endArg, ok := args[1].(*NumberValue)
		if !ok {
			return nil, fmt.Errorf("string.substring second argument must be a number")
		}
		end = int(endArg.Value)
	}
	if start.Value < 0 || start.Value > float64(len(s.Value)) || end < 0 || end > len(s.Value) {
		return nil, fmt.Errorf("string.substring indices out of bounds")
	}
	result := MakeString(s.Value[int(start.Value):end])
	return result, nil
}

func stringSplit(s *StringValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string.split requires exactly one argument")
	}
	separator, ok := args[0].(*StringValue)
	if !ok {
		return nil, fmt.Errorf("string.split argument must be a string")
	}
	parts := strings.Split(s.Value, separator.Value)
	result := MakeArray([]RuntimeValue{})
	for _, part := range parts {
		result.(*ArrayValue).Elements = append(result.(*ArrayValue).Elements, MakeString(part))
	}
	return result, nil
}

var ArrayPrototype = map[string]func(a *ArrayValue, args []RuntimeValue, env *Environment) (RuntimeValue, error){
	"length": arrayLength,
	"push":   arrayPush,
	"pop":    arrayPop,
	"join":   arrayJoin,
	// "filter":   arrayFilter,
	// "map":      arrayMap,
	// "find":     arrayFind,
	"includes": arrayIncludes,
}

// map to prototype functions
var StringPrototype = map[string]func(s *StringValue, args []RuntimeValue, env *Environment) (RuntimeValue, error){
	"length":      stringLength,
	"toUpperCase": stringToUpperCase,
	"toLowerCase": stringToLowerCase,
	"charAt":      stringCharAt,
	"substring":   stringSubstring,
	"split":       stringSplit,
}
