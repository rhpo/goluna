package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type ValueType string

const (
	NULL_TYPE      ValueType = "null"
	UNDEF_TYPE     ValueType = "undef"
	VOID_TYPE      ValueType = "void"
	NUMBER_TYPE    ValueType = "number"
	BOOLEAN_TYPE   ValueType = "boolean"
	STRING_TYPE    ValueType = "string"
	FUNCTION_TYPE  ValueType = "function"
	NATIVE_FN_TYPE ValueType = "native-fn"
	ARRAY_TYPE     ValueType = "array"
	OBJECT_TYPE    ValueType = "object"
	RETURN_TYPE    ValueType = "return"
)

type RuntimeValue interface {
	Type() ValueType
	String() string
	IsTruthy() bool
	Prototypes() *[]RuntimeValue
}

// Null Value
type NullValue struct{}

func (n *NullValue) Type() ValueType { return NULL_TYPE }
func (n *NullValue) String() string  { return "null" }
func (n *NullValue) IsTruthy() bool  { return false }
func (n *NullValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue
	return &prototypes
}

// Undefined Value
type UndefinedValue struct{}

func (u *UndefinedValue) Type() ValueType { return UNDEF_TYPE }
func (u *UndefinedValue) String() string  { return "undef" }
func (u *UndefinedValue) IsTruthy() bool  { return false }
func (u *UndefinedValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue
	return &prototypes
}

// Void Value
type VoidValue struct{}

func (v *VoidValue) Type() ValueType { return VOID_TYPE }
func (v *VoidValue) String() string  { return "" }
func (v *VoidValue) IsTruthy() bool  { return false }
func (v *VoidValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue
	return &prototypes
}

// Number Value
type NumberValue struct {
	Value float64
}

func (n *NumberValue) Type() ValueType { return NUMBER_TYPE }
func (n *NumberValue) String() string {
	if n.Value == float64(int64(n.Value)) {
		return strconv.FormatInt(int64(n.Value), 10)
	}
	return strconv.FormatFloat(n.Value, 'g', -1, 64)
}
func (n *NumberValue) IsTruthy() bool { return n.Value != 0 && !math.IsNaN(n.Value) }

// Prototypes returns an empty slice for NumberValue
func (n *NumberValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue

	prototypes = append(prototypes, MakeNativeFunction("string", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("string() requires 1 argument")
		}
		return MakeString(args[0].String()), nil
	})) // NaN prototype

	return &prototypes
}

// Boolean Value
type BooleanValue struct {
	Value bool
}

func (b *BooleanValue) Type() ValueType { return BOOLEAN_TYPE }
func (b *BooleanValue) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}
func (b *BooleanValue) IsTruthy() bool { return b.Value }
func (b *BooleanValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue

	prototypes = append(prototypes, MakeNativeFunction("string", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if b.Value {
			return MakeString("true"), nil
		}
		return MakeString("false"), nil
	})) // Boolean prototype

	return &prototypes
}

// String Value
type StringValue struct {
	Value string
}

func (s *StringValue) Type() ValueType { return STRING_TYPE }
func (s *StringValue) String() string  { return fmt.Sprintf("'%s'", s.Value) }
func (s *StringValue) IsTruthy() bool  { return s.Value != "" }
func (s *StringValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue
	for name, f := range StringPrototype {
		prototypes = append(prototypes, MakeNativeFunction(name, func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
			val, err := f(s, args, env)
			if err != nil {
				return nil, err
			}
			return val, nil
		}))
	}

	return &prototypes
}

// Array Value
type ArrayValue struct {
	Elements []RuntimeValue
}

func (a *ArrayValue) Type() ValueType { return ARRAY_TYPE }
func (a *ArrayValue) String() string {
	var elements []string
	for _, elem := range a.Elements {
		elements = append(elements, elem.String())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}
func (a *ArrayValue) IsTruthy() bool { return len(a.Elements) > 0 }
func (a *ArrayValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue

	// arrayPrototype contains methods for ArrayValue
	for name, fn := range ArrayPrototype {
		prototypes = append(prototypes, MakeNativeFunction(name, func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
			val, err := fn(a, args, env)
			if err != nil {
				return nil, err
			}
			return val, nil
		}))
	}

	return &prototypes
}

// Object Value
type ObjectValue struct {
	Properties map[string]RuntimeValue
}

func (o *ObjectValue) Type() ValueType { return OBJECT_TYPE }
func (o *ObjectValue) String() string {
	var props []string
	for key, value := range o.Properties {
		props = append(props, fmt.Sprintf("%s: %s", key, value.String()))
	}
	return "{" + strings.Join(props, ", ") + "}"
}
func (o *ObjectValue) IsTruthy() bool { return len(o.Properties) > 0 }
func (o *ObjectValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue

	prototypes = append(prototypes, MakeNativeFunction("keys", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		keys := make([]RuntimeValue, 0, len(o.Properties))
		for key := range o.Properties {
			keys = append(keys, MakeString(key))
		}
		return MakeArray(keys), nil
	}))

	prototypes = append(prototypes, MakeNativeFunction("values", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		values := make([]RuntimeValue, 0, len(o.Properties))
		for _, value := range o.Properties {
			values = append(values, value)
		}
		return MakeArray(values), nil
	}))

	return &prototypes
}

// Function Value
type FunctionValue struct {
	Name           string
	Parameters     []Parameter
	Body           []Statement
	DeclarationEnv *Environment
	Export         bool
	Anonymous      bool
}

func (f *FunctionValue) String() string {
	var paramStrs []string
	for _, param := range f.Parameters {
		if param.DefaultValue != nil {
			paramStrs = append(paramStrs, fmt.Sprintf("%s=(...)", param.Name))
		} else {
			paramStrs = append(paramStrs, param.Name)
		}
	}

	if f.IsAnonymous() {
		return fmt.Sprintf("lambda %s { ... }", strings.Join(paramStrs, " "))
	}
	return fmt.Sprintf("fn %s %s { ... }", f.Name, strings.Join(paramStrs, " "))
}
func (f *FunctionValue) IsTruthy() bool    { return true }
func (f *FunctionValue) IsAnonymous() bool { return f.Anonymous }
func (f *FunctionValue) Type() ValueType   { return FUNCTION_TYPE }
func (f *FunctionValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue

	// Add a prototype for calling the function
	prototypes = append(prototypes, MakeNativeFunction("call", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) < len(f.Parameters) {
			return nil, fmt.Errorf("not enough arguments to call function %s", f.Name)
		}
		if len(args) > len(f.Parameters) {
			return nil, fmt.Errorf("too many arguments to call function %s", f.Name)
		}

		// Create a new environment for the function call
		callEnv := NewEnvironment(f.DeclarationEnv)

		for i, param := range f.Parameters {
			callEnv.DeclareVar(param.Name, args[i], false)
		}

		// Execute the function body
		var returnValue RuntimeValue
		for _, stmt := range f.Body {
			result, err := Evaluate(stmt, callEnv)
			if err != nil {
				return nil, err
			}
			if result != nil && result.Type() == RETURN_TYPE {
				returnValue = result.(*ReturnValue).Value
				break
			}
		}

		if returnValue == nil {
			return MakeVoid(), nil
		}
		return returnValue, nil
	}))

	return &prototypes
}

// Native Function Value
type NativeFunctionCall func(args []RuntimeValue, env *Environment) (RuntimeValue, error)

type NativeFunctionValue struct {
	Name string
	Call NativeFunctionCall
}

func (n *NativeFunctionValue) Type() ValueType { return NATIVE_FN_TYPE }
func (n *NativeFunctionValue) String() string {
	return fmt.Sprintf("fn %s", n.Name)
}
func (n *NativeFunctionValue) IsTruthy() bool { return true }
func (n *NativeFunctionValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue
	return &prototypes
}

// Return Value (for control flow)
type ReturnValue struct {
	Value RuntimeValue
}

func (r *ReturnValue) Type() ValueType { return RETURN_TYPE }
func (r *ReturnValue) String() string  { return r.Value.String() }
func (r *ReturnValue) IsTruthy() bool  { return r.Value.IsTruthy() }
func (r *ReturnValue) Prototypes() *[]RuntimeValue {
	var prototypes []RuntimeValue
	return &prototypes
}

// Helper functions to create values
func MakeNull() RuntimeValue {
	return &NullValue{}
}

func MakeUndefined() RuntimeValue {
	return &UndefinedValue{}
}

func MakeVoid() RuntimeValue {
	return &VoidValue{}
}

func MakeNumber(value float64) RuntimeValue {
	return &NumberValue{Value: value}
}

func MakeBool(value bool) RuntimeValue {
	return &BooleanValue{Value: value}
}

func MakeString(value string) RuntimeValue {
	return &StringValue{Value: value}
}

func MakeArray(elements []RuntimeValue) RuntimeValue {
	return &ArrayValue{Elements: elements}
}

func MakeObject(properties map[string]RuntimeValue) RuntimeValue {
	return &ObjectValue{Properties: properties}
}

// Update MakeFunction to use Parameter struct
func MakeFunction(name string, parameters []Parameter, body []Statement, env *Environment, export bool, anonymous bool) RuntimeValue {
	return &FunctionValue{
		Name:           name,
		Parameters:     parameters,
		Body:           body,
		DeclarationEnv: env,
		Export:         export,
		Anonymous:      anonymous,
	}
}

func MakeNativeFunction(name string, call NativeFunctionCall) RuntimeValue {
	return &NativeFunctionValue{Name: name, Call: call}
}

func MakeReturn(value RuntimeValue) RuntimeValue {
	return &ReturnValue{Value: value}
}
