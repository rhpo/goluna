package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var startTime = time.Now()

func setupNativeFunctions(env *Environment) {

	// I/O functions

	// String functions
	env.DeclareVar("length", MakeNativeFunction("length", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("length expects 1 argument, got %d", len(args))
		}

		switch args[0].Type() {
		case STRING_TYPE:
			return MakeNumber(float64(len(args[0].(*StringValue).Value))), nil
		case ARRAY_TYPE:
			return MakeNumber(float64(len(args[0].(*ArrayValue).Elements))), nil
		case OBJECT_TYPE:
			return MakeNumber(float64(len(args[0].(*ObjectValue).Properties))), nil
		default:
			return nil, fmt.Errorf("length not supported for type %s", args[0].Type())
		}
	}), true)

	// Type conversion functions
	env.DeclareVar("int", MakeNativeFunction("int", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("int expects 1 argument, got %d", len(args))
		}

		switch args[0].Type() {
		case NUMBER_TYPE:
			value := args[0].(*NumberValue).Value
			return MakeNumber(float64(int64(value))), nil
		case STRING_TYPE:
			value := args[0].(*StringValue).Value
			if parsed, err := strconv.ParseFloat(value, 64); err == nil {
				return MakeNumber(float64(int64(parsed))), nil
			}
			return MakeNumber(0), nil
		default:
			return MakeNumber(0), nil
		}
	}), true)

	env.DeclareVar("float", MakeNativeFunction("float", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("float expects 1 argument, got %d", len(args))
		}

		switch args[0].Type() {
		case NUMBER_TYPE:
			return args[0], nil
		case STRING_TYPE:
			value := args[0].(*StringValue).Value
			if parsed, err := strconv.ParseFloat(value, 64); err == nil {
				return MakeNumber(parsed), nil
			}
			return MakeNumber(0), nil
		default:
			return MakeNumber(0), nil
		}
	}), true)

	env.DeclareVar("string", MakeNativeFunction("string", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("string expects 1 argument, got %d", len(args))
		}

		switch args[0].Type() {
		case STRING_TYPE:
			return args[0], nil
		case NUMBER_TYPE:
			value := args[0].(*NumberValue).Value
			return MakeString(strconv.FormatFloat(value, 'g', -1, 64)), nil
		case BOOLEAN_TYPE:
			value := args[0].(*BooleanValue).Value
			return MakeString(strconv.FormatBool(value)), nil
		default:
			return MakeString(args[0].String()), nil
		}
	}), true)

	// Type checking function
	env.DeclareVar("typeof", MakeNativeFunction("typeof", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("typeget expects 1 argument, got %d", len(args))
		}
		return MakeString(string(args[0].Type())), nil
	}), true)

	// Constants
	env.DeclareVar("true", MakeBool(true), true)
	env.DeclareVar("false", MakeBool(false), true)
	env.DeclareVar("null", MakeNull(), true)
	env.DeclareVar("undef", MakeUndefined(), true)

	// Exit function
	env.DeclareVar("exit", MakeNativeFunction("exit", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		fmt.Println(gray("Exiting..."))
		os.Exit(0)
		return MakeVoid(), nil
	}), true)

	// OBJECTS ---
	// Create IO object with all math functions
	IOObject := createIOObject()
	env.DeclareVar("io", IOObject, true)

	// Create math object with all math functions
	mathObject := createMathObject()
	env.DeclareVar("math", mathObject, true)
}

func createIOObject() RuntimeValue {
	ioProps := make(map[string]RuntimeValue)

	// Math functions
	ioProps["print"] = MakeNativeFunction("print", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		var output []string
		for _, arg := range args {
			if arg.Type() == STRING_TYPE {
				output = append(output, arg.(*StringValue).Value)
			} else {
				// Use colorized output for non-string values
				output = append(output, colorizeValue(arg, false, true))
			}
		}
		fmt.Println(strings.Join(output, " "))
		return MakeVoid(), nil
	})

	ioProps["input"] = MakeNativeFunction("input", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) > 0 && args[0].Type() == STRING_TYPE {
			fmt.Print(args[0].(*StringValue).Value)
		}

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			return MakeString(scanner.Text()), nil
		}
		return MakeString(""), nil
	})

	ioProps["time"] = MakeNativeFunction("time", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		elapsed := time.Since(startTime).Seconds() * 1000 // milliseconds
		return MakeNumber(elapsed), nil
	})

	return MakeObject(ioProps)
}

func createMathObject() RuntimeValue {
	mathProps := make(map[string]RuntimeValue)

	// Math functions
	mathProps["abs"] = MakeNativeFunction("abs", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("abs expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("abs expects a number")
		}
		value := args[0].(*NumberValue).Value
		return MakeNumber(math.Abs(value)), nil
	})

	mathProps["sqrt"] = MakeNativeFunction("sqrt", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("sqrt expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("sqrt expects a number")
		}
		value := args[0].(*NumberValue).Value
		return MakeNumber(math.Sqrt(value)), nil
	})

	mathProps["pow"] = MakeNativeFunction("pow", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("pow expects 2 arguments, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE || args[1].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("pow expects numbers")
		}
		base := args[0].(*NumberValue).Value
		exp := args[1].(*NumberValue).Value
		return MakeNumber(math.Pow(base, exp)), nil
	})

	mathProps["sin"] = MakeNativeFunction("sin", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("sin expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("sin expects a number")
		}
		value := args[0].(*NumberValue).Value
		return MakeNumber(math.Sin(value)), nil
	})

	mathProps["cos"] = MakeNativeFunction("cos", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("cos expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("cos expects a number")
		}
		value := args[0].(*NumberValue).Value
		return MakeNumber(math.Cos(value)), nil
	})

	mathProps["tan"] = MakeNativeFunction("tan", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("tan expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("tan expects a number")
		}
		value := args[0].(*NumberValue).Value
		return MakeNumber(math.Tan(value)), nil
	})

	mathProps["floor"] = MakeNativeFunction("floor", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("floor expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("floor expects a number")
		}
		value := args[0].(*NumberValue).Value
		return MakeNumber(math.Floor(value)), nil
	})

	mathProps["ceil"] = MakeNativeFunction("ceil", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("ceil expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("ceil expects a number")
		}
		value := args[0].(*NumberValue).Value
		return MakeNumber(math.Ceil(value)), nil
	})

	mathProps["round"] = MakeNativeFunction("round", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("round expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("round expects a number")
		}
		value := args[0].(*NumberValue).Value
		return MakeNumber(math.Round(value)), nil
	})

	mathProps["log"] = MakeNativeFunction("log", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("log expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("log expects a number")
		}
		value := args[0].(*NumberValue).Value
		return MakeNumber(math.Log(value)), nil
	})

	mathProps["exp"] = MakeNativeFunction("exp", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("exp expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("exp expects a number")
		}
		value := args[0].(*NumberValue).Value
		return MakeNumber(math.Exp(value)), nil
	})

	mathProps["min"] = MakeNativeFunction("min", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) == 0 {
			return MakeNumber(math.Inf(1)), nil
		}

		min := math.Inf(1)
		for _, arg := range args {
			if arg.Type() != NUMBER_TYPE {
				return nil, fmt.Errorf("min expects numbers")
			}
			value := arg.(*NumberValue).Value
			if value < min {
				min = value
			}
		}
		return MakeNumber(min), nil
	})

	mathProps["max"] = MakeNativeFunction("max", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		if len(args) == 0 {
			return MakeNumber(math.Inf(-1)), nil
		}

		max := math.Inf(-1)
		for _, arg := range args {
			if arg.Type() != NUMBER_TYPE {
				return nil, fmt.Errorf("max expects numbers")
			}
			value := arg.(*NumberValue).Value
			if value > max {
				max = value
			}
		}
		return MakeNumber(max), nil
	})

	mathProps["random"] = MakeNativeFunction("random", func(args []RuntimeValue, env *Environment) (RuntimeValue, error) {
		return MakeNumber(rand.Float64()), nil
	})

	// Math constants
	mathProps["PI"] = MakeNumber(math.Pi)
	mathProps["E"] = MakeNumber(math.E)
	mathProps["LN2"] = MakeNumber(math.Ln2)
	mathProps["LN10"] = MakeNumber(math.Ln10)
	mathProps["LOG2E"] = MakeNumber(math.Log2E)
	mathProps["LOG10E"] = MakeNumber(math.Log10E)
	mathProps["SQRT1_2"] = MakeNumber(math.Sqrt2 / 2)
	mathProps["SQRT2"] = MakeNumber(math.Sqrt2)

	return MakeObject(mathProps)
}
