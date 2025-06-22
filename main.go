package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func main() {

	// get args
	args := make([]string, 0)
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--") {
			// skip flags
			continue
		}
		if strings.HasPrefix(arg, "-") {
			// skip short flags
			continue
		}
		args = append(args, arg)
	}

	// If there are arguments, treat them as a file to execute
	if len(args) > 0 {
		filename := args[0]
		if len(args) > 1 {
			fmt.Println("Error: Too many arguments. Only one file can be executed at a time.")
			return
		}

		// try to read the relative file (using fs library)
		data, err := fs.ReadFile(os.DirFS("."), filename)
		if err != nil {
			fmt.Printf("Error: Could not read file '%s': %v\n", filename, err)
			return
		}

		// Create a new Luna instance and< evaluate the file content
		env := NewEnvironment(nil)
		setupNativeFunctions(env)

		luna := NewLuna(env)
		result, err := luna.Evaluate(string(data))

		if err != nil {
			fmt.Println(formatError("Error", err.Error()))
			return
		}

		if result != nil && result.Type() != VOID_TYPE {
			// Colorize the output
			output := colorizeValue(result, false, false)
			if output != "" {
				fmt.Println(output)
			}
		}

		return

	}

	// Welcome message with colors
	fmt.Println(green("Welcome to the Luna REPL!"))
	fmt.Println(gray("Type ") + green(under("exit()")) + gray(" to leave..."))

	env := NewEnvironment(nil)
	setupNativeFunctions(env)

	readline := NewReadline(white(">> "))

	for {
		input, err := readline.ReadLine()
		if err != nil {
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if input == "exit()" {
			fmt.Println(gray("Exiting..."))
			break
		}

		// Check for balanced brackets
		if !isBalanced(input) {
			for {
				nesting := countNesting(input)
				fmt.Print(strings.Repeat("  ", nesting) + gray("... "))
				line, err := readline.ReadLine(false)
				if err != nil {
					break
				}
				input += " " + line
				if isBalanced(input) {
					break
				}
			}
		}

		luna := NewLuna(env)
		result, err := luna.Evaluate(input)
		if err != nil {
			// Format error with colors
			fmt.Println(formatError("Error", err.Error()))
		} else if result != nil && result.Type() != VOID_TYPE {
			// Colorize the output
			output := colorizeValue(result, false, false)
			if output != "" {
				fmt.Println(output)
			}
		}
	}
}

func isBalanced(input string) bool {
	stack := 0
	inString := false
	escaped := false

	for _, char := range input {
		if escaped {
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if char == '"' || char == '\'' {
			inString = !inString
			continue
		}

		if !inString {
			switch char {
			case '{', '(', '[':
				stack++
			case '}', ')', ']':
				stack--
			}
		}
	}

	return stack == 0
}

// countNesting returns the current nesting level of brackets in the input string.
func countNesting(input string) int {
	stack := 0
	inString := false
	escaped := false

	for _, char := range input {
		if escaped {
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if char == '"' || char == '\'' {
			inString = !inString
			continue
		}

		if !inString {
			switch char {
			case '{', '(', '[':
				stack++
			case '}', ')', ']':
				if stack > 0 {
					stack--
				}
			}
		}
	}

	return stack
}
