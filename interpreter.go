package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Evaluate(node Statement, env *Environment) (RuntimeValue, error) {
	switch n := node.(type) {
	case *Program:
		return evaluateProgram(n, env)
	case *NumericLiteral:
		return MakeNumber(n.Value), nil
	case *StringLiteral:
		return evaluateStringLiteral(n, env)
	case *BooleanLiteral:
		return MakeBool(n.Value), nil
	case *UndefinedLiteral:
		return MakeUndefined(), nil
	case *NullLiteral:
		return MakeNull(), nil
	case *Identifier:
		return evaluateIdentifier(n, env)
	case *ArrayLiteral:
		return evaluateArrayLiteral(n, env)
	case *ObjectLiteral:
		return evaluateObjectLiteral(n, env)
	case *BinaryExpr:
		return evaluateBinaryExpression(n, env)
	case *UnaryExpr:
		return evaluateUnaryExpression(n, env)
	case *AssignmentExpr:
		return evaluateAssignmentExpression(n, env)
	case *ActionAssignmentExpr:
		return evaluateActionAssignmentExpression(n, env)
	case *CallExpr:
		return evaluateCallExpression(n, env)
	case *MemberExpr:
		return evaluateMemberExpression(n, env)
	case *TernaryExpr:
		return evaluateTernaryExpression(n, env)
	case *TypeofExpr:
		return evaluateTypeofExpression(n, env)
	case *EqualityExpr:
		return evaluateEqualityExpression(n, env)
	case *InequalityExpr:
		return evaluateInequalityExpression(n, env)
	case *LogicalExpr:
		return evaluateLogicalExpression(n, env)
	case *FunctionDeclaration:
		return evaluateFunctionDeclaration(n, env)
	case *IfStatement:
		return evaluateIfStatement(n, env)
	case *WhileStatement:
		return evaluateWhileStatement(n, env)
	case *ForStatement:
		return evaluateForStatement(n, env)
	case *ReturnExpr:
		value, err := Evaluate(n.Value, env)
		if err != nil {
			return nil, err
		}
		return MakeReturn(value), nil
	case *DebugStatement:
		return evaluateDebugStatement(n, env)
	default:
		return nil, fmt.Errorf("unsupported AST node: %T", node)
	}
}

func evaluateProgram(program *Program, env *Environment) (RuntimeValue, error) {
	var lastEvaluated RuntimeValue = MakeVoid()

	for _, statement := range program.Body {
		result, err := Evaluate(statement, env)
		if err != nil {
			return nil, err
		}
		if result != nil {
			lastEvaluated = result
		}
	}

	return lastEvaluated, nil
}

func evaluateStringLiteral(node *StringLiteral, env *Environment) (RuntimeValue, error) {
	// Handle string interpolation
	value := node.Value
	if strings.Contains(value, "{") {
		// Simple string interpolation - replace {variable} with variable value
		result := value
		// This is a simplified version - in a full implementation you'd parse expressions
		for name, val := range env.variables {
			placeholder := "{" + name + "}"
			if strings.Contains(result, placeholder) {
				if val.Type() == STRING_TYPE {
					result = strings.ReplaceAll(result, placeholder, val.(*StringValue).Value)
				} else {
					result = strings.ReplaceAll(result, placeholder, val.String())
				}
			}
		}
		return MakeString(result), nil
	}
	return MakeString(value), nil
}

func evaluateIdentifier(node *Identifier, env *Environment) (RuntimeValue, error) {
	myVar := env.LookupVar(node.Value)
	if myVar == nil {
		return nil, fmt.Errorf("undefined variable: %s", node.Value)
	}

	return myVar, nil
}

func evaluateArrayLiteral(node *ArrayLiteral, env *Environment) (RuntimeValue, error) {
	elements := make([]RuntimeValue, len(node.Elements))
	for i, elem := range node.Elements {
		value, err := Evaluate(elem, env)
		if err != nil {
			return nil, err
		}
		elements[i] = value
	}
	return MakeArray(elements), nil
}

func evaluateObjectLiteral(node *ObjectLiteral, env *Environment) (RuntimeValue, error) {
	properties := make(map[string]RuntimeValue)
	for _, prop := range node.Properties {
		value, err := Evaluate(prop.Value, env)
		if err != nil {
			return nil, err
		}
		properties[prop.Key] = value
	}
	return MakeObject(properties), nil
}

func evaluateBinaryExpression(node *BinaryExpr, env *Environment) (RuntimeValue, error) {
	left, err := Evaluate(node.Left, env)
	if err != nil {
		return nil, err
	}

	right, err := Evaluate(node.Right, env)
	if err != nil {
		return nil, err
	}

	return evaluateBinaryOperation(left, right, node.Operator)
}

func evaluateBinaryOperation(left, right RuntimeValue, operator string) (RuntimeValue, error) {
	// Handle numeric operations
	if left.Type() == NUMBER_TYPE && right.Type() == NUMBER_TYPE {
		leftVal := left.(*NumberValue).Value
		rightVal := right.(*NumberValue).Value

		switch operator {
		case "+":
			return MakeNumber(leftVal + rightVal), nil
		case "-":
			return MakeNumber(leftVal - rightVal), nil
		case "*":
			return MakeNumber(leftVal * rightVal), nil
		case "/":
			if rightVal == 0 {
				return MakeNumber(math.Inf(1)), nil
			}
			return MakeNumber(leftVal / rightVal), nil
		case "%":
			return MakeNumber(math.Mod(leftVal, rightVal)), nil
		case "**":
			return MakeNumber(math.Pow(leftVal, rightVal)), nil
		}
	}

	// Handle string concatenation
	if operator == "+" && (left.Type() == STRING_TYPE || right.Type() == STRING_TYPE) {
		leftStr := left.String()
		rightStr := right.String()
		if left.Type() == STRING_TYPE {
			leftStr = left.(*StringValue).Value
		}
		if right.Type() == STRING_TYPE {
			rightStr = right.(*StringValue).Value
		}
		return MakeString(leftStr + rightStr), nil
	}

	return nil, fmt.Errorf("unsupported binary operation: %s %s %s", left.Type(), operator, right.Type())
}

func evaluateUnaryExpression(node *UnaryExpr, env *Environment) (RuntimeValue, error) {
	// Handle postfix increment/decrement
	if node.Operator == "++_post" || node.Operator == "--_post" {
		// Only valid on identifiers
		ident, ok := node.Value.(*Identifier)
		if !ok {
			return nil, fmt.Errorf("postfix operator only valid on identifiers")
		}
		val := env.LookupVar(ident.Value)
		if val == nil || val.Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("cannot apply %s to non-number variable", node.Operator[:2])
		}
		oldVal := val.(*NumberValue).Value
		var newVal float64
		if node.Operator == "++_post" {
			newVal = oldVal + 1
		} else {
			newVal = oldVal - 1
		}
		env.AssignVar(ident.Value, MakeNumber(newVal))
		return MakeNumber(oldVal), nil // Return old value (postfix)
	}

	// Prefix unary
	switch node.Operator {
	case "!":
		value, err := Evaluate(node.Value, env)
		if err != nil {
			return nil, err
		}
		return MakeBool(!value.IsTruthy()), nil
	case "-":
		value, err := Evaluate(node.Value, env)
		if err != nil {
			return nil, err
		}
		if value.Type() == NUMBER_TYPE {
			return MakeNumber(-value.(*NumberValue).Value), nil
		}
		return nil, fmt.Errorf("cannot negate non-number value")
	case "+":
		value, err := Evaluate(node.Value, env)
		if err != nil {
			return nil, err
		}
		if value.Type() == NUMBER_TYPE {
			return value, nil
		}
		return nil, fmt.Errorf("cannot apply unary plus to non-number value")
	case "++":
		ident, ok := node.Value.(*Identifier)
		if !ok {
			return nil, fmt.Errorf("prefix ++ only valid on identifiers")
		}
		val := env.LookupVar(ident.Value)
		if val == nil || val.Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("cannot increment non-number variable")
		}
		newVal := val.(*NumberValue).Value + 1
		env.AssignVar(ident.Value, MakeNumber(newVal))
		return MakeNumber(newVal), nil // Return new value (prefix)
	case "--":
		ident, ok := node.Value.(*Identifier)
		if !ok {
			return nil, fmt.Errorf("prefix -- only valid on identifiers")
		}
		val := env.LookupVar(ident.Value)
		if val == nil || val.Type() != NUMBER_TYPE {
			return nil, fmt.Errorf("cannot decrement non-number variable")
		}
		newVal := val.(*NumberValue).Value - 1
		env.AssignVar(ident.Value, MakeNumber(newVal))
		return MakeNumber(newVal), nil // Return new value (prefix)
	}

	return nil, fmt.Errorf("unsupported unary operator: %s", node.Operator)
}

func evaluateAssignmentExpression(node *AssignmentExpr, env *Environment) (RuntimeValue, error) {
	if identifier, ok := node.Assigne.(*Identifier); ok {
		value, err := Evaluate(node.Value, env)
		if err != nil {
			return nil, err
		}

		// Fix: Check if variable exists in current or parent environment
		// If it exists, assign to existing variable instead of creating new one
		if env.HasVar(identifier.Value) {
			return env.AssignVar(identifier.Value, value), nil
		} else {
			return env.DeclareVar(identifier.Value, value, false), nil
		}
	} else if memberExpr, ok := node.Assigne.(*MemberExpr); ok {
		object, err := Evaluate(memberExpr.Object, env)
		if err != nil {
			return nil, err
		}

		var property RuntimeValue
		if memberExpr.Property.Kind() == IDENTIFIER_NODE {
			ident := memberExpr.Property.(*Identifier)
			property = MakeString(ident.Value)
		} else {
			prop, err := Evaluate(memberExpr.Property, env)
			if err != nil {
				return nil, err
			}
			property = prop
		}

		value, err := Evaluate(node.Value, env)
		if err != nil {
			return nil, err
		}
		// is it object or array
		if object.Type() == OBJECT_TYPE {
			objectVal := object.(*ObjectValue)
			if str, ok := property.(*StringValue); ok {
				objectVal.Properties[str.Value] = value
			} else {
				numVal := fmt.Sprint(property.(*NumberValue).Value)
				objectVal.Properties[numVal] = value
			}
			return value, nil
		} else if object.Type() == ARRAY_TYPE {
			arrayVal := object.(*ArrayValue)
			index := int(property.(*NumberValue).Value)
			arrayVal.Elements[index] = value
			return value, nil
		} else {
			return nil, fmt.Errorf("cannot assign to non-object (%s)", object.Type())
		}
	}

	return nil, fmt.Errorf("invalid assignment target")
}

func evaluateActionAssignmentExpression(node *ActionAssignmentExpr, env *Environment) (RuntimeValue, error) {
	if identifier, ok := node.Assigne.(*Identifier); ok {
		value, err := Evaluate(node.Value, env)
		if err != nil {
			return nil, err
		}

		switch node.Action.Name {
		case "const":
			return env.DeclareVar(identifier.Value, value, true), nil
		case "var":
			return env.DeclareVar(identifier.Value, value, false), nil
		case "out":
			// Mark as exported (simplified - just declare normally for now)
			return env.DeclareVar(identifier.Value, value, false), nil
		default:
			return nil, fmt.Errorf("unsupported action: %s", node.Action.Name)
		}
	}

	return nil, fmt.Errorf("invalid assignment target")
}

func evaluateCallExpression(node *CallExpr, env *Environment) (RuntimeValue, error) {
	fn, err := Evaluate(node.Caller, env)
	if err != nil {
		return nil, err
	}

	args := make([]RuntimeValue, len(node.Args))
	for i, arg := range node.Args {
		value, err := Evaluate(arg, env)
		if err != nil {
			return nil, err
		}
		args[i] = value
	}

	switch f := fn.(type) {
	case *FunctionValue:
		return callFunction(f, args, env)
	case *NativeFunctionValue:
		return f.Call(args, env)
	default:
		return nil, fmt.Errorf("cannot call non-function value")
	}
}

func callFunction(fn *FunctionValue, args []RuntimeValue, env *Environment) (RuntimeValue, error) {
	// Create new scope for function execution
	fnEnv := NewEnvironment(fn.DeclarationEnv)

	// Bind parameters with default value support
	for i, param := range fn.Parameters {
		var value RuntimeValue = MakeUndefined()

		if i < len(args) {
			// Use provided argument
			value = args[i]
		} else if param.DefaultValue != nil {
			// Use default value
			defaultVal, err := Evaluate(param.DefaultValue, fn.DeclarationEnv)
			if err != nil {
				return nil, fmt.Errorf("error evaluating default parameter %s: %v", param.Name, err)
			}
			value = defaultVal
		}
		// If no argument and no default, value remains undefined

		fnEnv.DeclareVar(param.Name, value, false)
	}

	// Execute function body
	var result RuntimeValue = MakeVoid()
	for _, stmt := range fn.Body {
		val, err := Evaluate(stmt, fnEnv)
		if err != nil {
			return nil, err
		}
		if val != nil {
			if val.Type() == RETURN_TYPE {
				return val.(*ReturnValue).Value, nil
			}
			result = val
		}
	}

	return result, nil
}

func evaluateMemberExpression(node *MemberExpr, env *Environment) (RuntimeValue, error) {
	object, err := Evaluate(node.Object, env)
	if err != nil {
		return nil, err
	}

	var key string
	if node.Computed {
		prop, err := Evaluate(node.Property, env)
		if err != nil {
			return nil, err
		}
		if prop.Type() == STRING_TYPE {
			key = prop.(*StringValue).Value
		} else if prop.Type() == NUMBER_TYPE {
			key = strconv.FormatFloat(prop.(*NumberValue).Value, 'g', -1, 64)
		} else {
			return nil, fmt.Errorf("invalid property key type")
		}
	} else {
		if identifier, ok := node.Property.(*Identifier); ok {
			key = identifier.Value
		} else {
			return nil, fmt.Errorf("invalid property access")
		}
	}

	switch obj := object.(type) {
	case *ArrayValue:
		if index, err := strconv.Atoi(key); err == nil {
			if index >= 0 && index < len(obj.Elements) {
				return obj.Elements[index], nil
			}
		}

		// Check prototypes for native functions
		for _, protoFn := range *obj.Prototypes() {
			if protoFn.(*NativeFunctionValue).Name == key {
				return protoFn, nil
			}
		}
		return MakeUndefined(), nil

	case *ObjectValue:
		if value, exists := obj.Properties[key]; exists {
			return value, nil
		}
		// Check prototypes for native functions
		for _, protoFn := range *obj.Prototypes() {
			if protoFn.(*NativeFunctionValue).Name == key {
				return protoFn, nil
			}
		}
		return MakeUndefined(), nil
	default:
		// Check prototypes for native functions
		for _, protoFn := range *obj.Prototypes() {
			if protoFn.(*NativeFunctionValue).Name == key {
				return protoFn, nil
			}
		}
		return MakeUndefined(), nil
	}
}

func evaluateTernaryExpression(node *TernaryExpr, env *Environment) (RuntimeValue, error) {
	condition, err := Evaluate(node.Condition, env)
	if err != nil {
		return nil, err
	}

	if condition.IsTruthy() {
		return Evaluate(node.Consequent, env)
	} else {
		return Evaluate(node.Alternate, env)
	}
}

func evaluateTypeofExpression(node *TypeofExpr, env *Environment) (RuntimeValue, error) {
	value, err := Evaluate(node.Value, env)
	if err != nil {
		return nil, err
	}

	return MakeString(string(value.Type())), nil
}

func evaluateEqualityExpression(node *EqualityExpr, env *Environment) (RuntimeValue, error) {
	left, err := Evaluate(node.Left, env)
	if err != nil {
		return nil, err
	}

	right, err := Evaluate(node.Right, env)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case "==":
		return MakeBool(isEqual(left, right)), nil
	case "!=":
		return MakeBool(!isEqual(left, right)), nil
	default:
		return nil, fmt.Errorf("unsupported equality operator: %s", node.Operator)
	}
}

func evaluateInequalityExpression(node *InequalityExpr, env *Environment) (RuntimeValue, error) {
	left, err := Evaluate(node.Left, env)
	if err != nil {
		return nil, err
	}

	right, err := Evaluate(node.Right, env)
	if err != nil {
		return nil, err
	}

	if left.Type() != NUMBER_TYPE || right.Type() != NUMBER_TYPE {
		return nil, fmt.Errorf("cannot compare non-numeric values")
	}

	leftVal := left.(*NumberValue).Value
	rightVal := right.(*NumberValue).Value

	switch node.Operator {
	case "<":
		return MakeBool(leftVal < rightVal), nil
	case ">":
		return MakeBool(leftVal > rightVal), nil
	case "<=":
		return MakeBool(leftVal <= rightVal), nil
	case ">=":
		return MakeBool(leftVal >= rightVal), nil
	default:
		return nil, fmt.Errorf("unsupported inequality operator: %s", node.Operator)
	}
}

func evaluateLogicalExpression(node *LogicalExpr, env *Environment) (RuntimeValue, error) {
	left, err := Evaluate(node.Left, env)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case "&&":
		if !left.IsTruthy() {
			return left, nil
		}
		return Evaluate(node.Right, env)
	case "||":
		if left.IsTruthy() {
			return left, nil
		}
		return Evaluate(node.Right, env)
	default:
		return nil, fmt.Errorf("unsupported logical operator: %s", node.Operator)
	}
}

func evaluateFunctionDeclaration(node *FunctionDeclaration, env *Environment) (RuntimeValue, error) {
	anonymous := node.Name == ""
	fn := MakeFunction(node.Name, node.Parameters, node.Body, env, node.Export, anonymous)
	if !anonymous {
		env.DeclareVar(node.Name, fn, true)
	}
	return fn, nil
}

func evaluateIfStatement(node *IfStatement, env *Environment) (RuntimeValue, error) {
	condition, err := Evaluate(node.Test, env)
	if err != nil {
		return nil, err
	}

	// Don't create new environment for if statements - use parent environment
	var result RuntimeValue = MakeVoid()

	if condition.IsTruthy() {
		for _, stmt := range node.Consequent {
			val, err := Evaluate(stmt, env) // Use parent env instead of new env
			if err != nil {
				return nil, err
			}
			if val != nil {
				if val.Type() == RETURN_TYPE {
					return val, nil
				}
				result = val
			}
		}
	} else if len(node.Alternate) > 0 {
		for _, stmt := range node.Alternate {
			val, err := Evaluate(stmt, env) // Use parent env instead of new env
			if err != nil {
				return nil, err
			}
			if val != nil {
				if val.Type() == RETURN_TYPE {
					return val, nil
				}
				result = val
			}
		}
	}

	return result, nil
}

func evaluateWhileStatement(node *WhileStatement, env *Environment) (RuntimeValue, error) {
	var result RuntimeValue = MakeVoid()

	for {
		condition, err := Evaluate(node.Test, env)
		if err != nil {
			return nil, err
		}

		if !condition.IsTruthy() {
			break
		}

		for _, stmt := range node.Consequent {
			val, err := Evaluate(stmt, env)
			if err != nil {
				return nil, err
			}
			if val != nil {
				if val.Type() == RETURN_TYPE {
					return val, nil
				}
				result = val
			}
		}
	}

	return result, nil
}

func evaluateForStatement(node *ForStatement, env *Environment) (RuntimeValue, error) {
	forEnv := NewEnvironment(env)
	var result RuntimeValue = MakeVoid()

	// Execute declaration
	_, err := Evaluate(node.Declaration, forEnv)
	if err != nil {
		return nil, err
	}

	for {
		// Test condition
		condition, err := Evaluate(node.Test, forEnv)
		if err != nil {
			return nil, err
		}

		if !condition.IsTruthy() {
			break
		}

		// Execute body
		for _, stmt := range node.Body {
			val, err := Evaluate(stmt, forEnv)
			if err != nil {
				return nil, err
			}
			if val != nil {
				if val.Type() == RETURN_TYPE {
					return val, nil
				}
				result = val
			}
		}

		// Execute increaser
		_, err = Evaluate(node.Increaser, forEnv)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func evaluateDebugStatement(node *DebugStatement, env *Environment) (RuntimeValue, error) {
	var props []string
	for _, prop := range node.Props {
		value, err := Evaluate(prop, env)
		if err != nil {
			return nil, err
		}
		props = append(props, colorizeValue(value, false, false))
	}

	fmt.Println(formatDebug(props))
	return MakeVoid(), nil
}

func isEqual(left, right RuntimeValue) bool {
	if left.Type() != right.Type() {
		return false
	}

	switch left.Type() {
	case NUMBER_TYPE:
		return left.(*NumberValue).Value == right.(*NumberValue).Value
	case BOOLEAN_TYPE:
		return left.(*BooleanValue).Value == right.(*BooleanValue).Value
	case STRING_TYPE:
		return left.(*StringValue).Value == right.(*StringValue).Value
	case NULL_TYPE, UNDEF_TYPE, VOID_TYPE:
		return true
	default:
		return false // Objects and arrays need deep comparison
	}
}
