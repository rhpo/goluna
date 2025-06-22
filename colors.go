package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	Italic = "\033[3m"
	Under  = "\033[4m"

	// Foreground colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Gray    = "\033[90m"

	// Background colors
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// Color functions
func colorize(text, color string) string {
	return color + text + Reset
}

func red(text string) string     { return colorize(text, Red) }
func green(text string) string   { return colorize(text, Green) }
func yellow(text string) string  { return colorize(text, Yellow) }
func blue(text string) string    { return colorize(text, Blue) }
func magenta(text string) string { return colorize(text, Magenta) }
func cyan(text string) string    { return colorize(text, Cyan) }
func white(text string) string   { return colorize(text, White) }
func gray(text string) string    { return colorize(text, Gray) }
func bold(text string) string    { return colorize(text, Bold) }
func dim(text string) string     { return colorize(text, Dim) }
func italic(text string) string  { return colorize(text, Italic) }
func under(text string) string   { return colorize(text, Under) }

// Colorize runtime values for output
func colorizeValue(result RuntimeValue, isInner bool, noString bool) string {
	if result == nil {
		return gray("null")
	}

	switch result.Type() {
	case STRING_TYPE:
		str := result.(*StringValue).Value
		if noString {
			return str
		}
		return green("'" + strings.ReplaceAll(str, "'", dim("'")) + "'")

	case ARRAY_TYPE:
		array := result.(*ArrayValue)
		maxElements := 16

		if len(array.Elements) <= maxElements {
			var elements []string
			for _, elem := range array.Elements {
				elements = append(elements, colorizeValue(elem, true, false))
			}
			return cyan("[") + strings.Join(elements, ", ") + cyan("]")
		} else {
			var elements []string
			for i := 0; i < maxElements; i++ {
				elements = append(elements, colorizeValue(array.Elements[i], true, false))
			}
			return cyan(fmt.Sprintf("(%d elements) ", len(array.Elements))) +
				yellow("[") + strings.Join(elements, ", ") + gray(", ...") + yellow("]")
		}

	case NUMBER_TYPE:
		num := result.(*NumberValue).Value
		if num != num { // NaN check
			return cyan("NaN")
		}
		if num == float64(int64(num)) {
			return yellow(strconv.FormatInt(int64(num), 10))
		}
		return yellow(strconv.FormatFloat(num, 'g', -1, 64))

	case UNDEF_TYPE:
		return gray("undef")

	case VOID_TYPE:
		return ""

	case FUNCTION_TYPE:
		fn := result.(*FunctionValue)
		var name string

		if fn.IsAnonymous() {
			var paramStrs []string
			for _, param := range fn.Parameters {
				if param.DefaultValue != nil {
					paramStrs = append(paramStrs, param.Name+"=(...)")
				} else {
					paramStrs = append(paramStrs, param.Name)
				}
			}
			name = magenta("lambda") + " " + strings.Join(paramStrs, " ")
		} else {
			exportPrefix := ""
			if fn.Export {
				exportPrefix = green("out") + " "
			}

			var paramStrs []string
			for _, param := range fn.Parameters {
				if param.DefaultValue != nil {
					paramStrs = append(paramStrs, green(param.Name)+yellow("=(...)"))
				} else {
					paramStrs = append(paramStrs, green(param.Name))
				}
			}

			name = exportPrefix + magenta("fn") + " " + bold(blue(fn.Name)) + " " +
				strings.Join(paramStrs, " ")

			if len(fn.Parameters) > 0 {
				name += " "
			}
		}

		bodyIndicator := ""
		if len(fn.Body) > 0 {
			bodyIndicator = " ... "
		}

		return name + gray(fmt.Sprintf("{%s}", bodyIndicator))

	case NATIVE_FN_TYPE:
		fn := result.(*NativeFunctionValue)
		if isInner {
			return magenta("fn") + " " + cyan(fn.Name)
		}
		return magenta("fn") + " " + cyan(fn.Name) + " {\n" +
			"  " + italic("(NAT-C)...") + "\n" +
			"}"

	case BOOLEAN_TYPE:
		return magenta(result.String())

	case NULL_TYPE:
		return magenta("null")

	case OBJECT_TYPE:
		obj := result.(*ObjectValue)
		if isInner {
			return gray("{ ... }")
		}

		var props []string
		for key, value := range obj.Properties {
			props = append(props, fmt.Sprintf("  %s: %s", blue(key), colorizeValue(value, true, false)))
		}

		if len(props) == 0 {
			return gray("{}")
		}

		return gray("{") + "\n" + strings.Join(props, ",\n") + "\n" + gray("}")

	default:
		return yellow(result.String())
	}
}

// Format error messages with colors
func formatError(errType, message string) string {
	return fmt.Sprintf("%s: %s", red(under(bold(errType))), gray(message))
}

// Format debug output
func formatDebug(props []string) string {
	debugStyle := BgYellow + Red
	return colorize(" DEBUG: ", debugStyle) + strings.Join(props, ", ")
}
