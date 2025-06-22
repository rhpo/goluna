package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	tokens   []Token
	position int
	code     string
}

func NewParser(tokens []Token, code string) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
		code:     code,
	}
}

func (p *Parser) ProduceAST() (Statement, error) {
	program := &Program{Body: []Statement{}}

	for !p.isEOF() {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			program.Body = append(program.Body, stmt)
		}
	}

	return program, nil
}

func (p *Parser) parseStatement() (Statement, error) {
	token := p.at()
	var returned Statement
	var err error

	switch token.Type {
	case OUT:
		returned, err = p.parseFunctionDeclaration()
	case FN:
		returned, err = p.parseFunctionDeclaration()
	case IF:
		returned, err = p.parseIfStatement()
	case WHILE:
		returned, err = p.parseWhileStatement()
	case FOR:
		returned, err = p.parseForStatement()
	case RETURN:
		returned, err = p.parseReturnStatement()
	case DEBUG:
		returned, err = p.parseDebugStatement()
	case USE:
		returned, err = p.parseUseStatement()
	case NEWLINE:
		p.eat() // Skip newlines
		returned, err = nil, nil
	default:
		returned, err = p.parseExpression()
	}

	// if ; then eat ;
	if p.at().Type == SEMICOLON {
		p.eat()
	}

	return returned, err
}

// Add error reporting helper
func (p *Parser) formatError(message string, token Token) error {
	lines := strings.Split(p.code, "\n")
	if token.Position.Line < len(lines) {
		line := lines[token.Position.Line]
		pointer := strings.Repeat(" ", token.Position.Column) + strings.Repeat("^", len(token.Value))
		if len(token.Value) == 0 {
			pointer = strings.Repeat(" ", token.Position.Column) + "^"
		}
		return fmt.Errorf("%s at line %d, column %d:\n%s\n%s",
			message, token.Position.Line+1, token.Position.Column+1, line, pointer)
	}
	return fmt.Errorf("%s at line %d, column %d", message, token.Position.Line+1, token.Position.Column+1)
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseAssignmentExpression()
}

func (p *Parser) parseAssignmentExpression() (Expression, error) {
	left, err := p.parseTernaryExpression()
	if err != nil {
		return nil, err
	}

	if p.at().Type == EQUALS {
		p.eat() // consume =
		// Fix: Use parseExpression to parse the right-hand side
		value, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return &AssignmentExpr{Assigne: left, Value: value}, nil
	}

	if p.at().Type == COLON {
		// Action assignment (const, var, out, etc.)
		p.eat() // consume :
		action := p.eat().Value

		if p.at().Type == EQUALS {
			p.eat() // consume =
			value, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			return &ActionAssignmentExpr{
				Assigne: left,
				Value:   value,
				Action:  ActionExpr{Name: action, Args: []Expression{}},
			}, nil
		}
	}

	return left, nil
}

func (p *Parser) parseTernaryExpression() (Expression, error) {
	expr, err := p.parseLogicalExpression()
	if err != nil {
		return nil, err
	}

	if p.at().Type == TERNARY {
		p.eat()                                       // consume ?
		consequent, err := p.parseLogicalExpression() // Parse up to logical level to avoid consuming the colon
		if err != nil {
			return nil, err
		}

		if p.at().Type != COLON {
			return nil, p.formatError("expected ':' in ternary expression", p.at())
		}
		p.eat() // consume :

		alternate, err := p.parseTernaryExpression() // Allow nested ternary
		if err != nil {
			return nil, err
		}

		return &TernaryExpr{
			Condition:  expr,
			Consequent: consequent,
			Alternate:  alternate,
		}, nil
	}

	return expr, nil
}

func (p *Parser) parseLogicalExpression() (Expression, error) {
	left, err := p.parseEqualityExpression()
	if err != nil {
		return nil, err
	}

	for p.at().Type == AND || p.at().Type == OR {
		operator := p.eat().Value
		right, err := p.parseEqualityExpression()
		if err != nil {
			return nil, err
		}
		left = &LogicalExpr{Left: left, Right: right, Operator: operator}
	}

	return left, nil
}

func (p *Parser) parseEqualityExpression() (Expression, error) {
	left, err := p.parseInequalityExpression()
	if err != nil {
		return nil, err
	}

	for p.at().Type == EQUALITY_OP || p.at().Type == INEQUALITY_OP {
		operator := p.eat().Value
		right, err := p.parseInequalityExpression()
		if err != nil {
			return nil, err
		}
		left = &EqualityExpr{Left: left, Right: right, Operator: operator}
	}

	return left, nil
}

func (p *Parser) parseInequalityExpression() (Expression, error) {
	left, err := p.parseAdditiveExpression()
	if err != nil {
		return nil, err
	}

	for p.at().Type == SMALLER_THAN || p.at().Type == GREATER_THAN ||
		p.at().Type == SMALLER_OR_EQUAL || p.at().Type == GREATER_OR_EQUAL {
		operator := p.eat().Value
		right, err := p.parseAdditiveExpression()
		if err != nil {
			return nil, err
		}
		left = &InequalityExpr{Left: left, Right: right, Operator: operator}
	}

	return left, nil
}

func (p *Parser) parseAdditiveExpression() (Expression, error) {
	left, err := p.parseMultiplicativeExpression()
	if err != nil {
		return nil, err
	}

	for p.at().Value == "+" || p.at().Value == "-" {
		operator := p.eat().Value
		right, err := p.parseMultiplicativeExpression()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Right: right, Operator: operator}
	}

	return left, nil
}

func (p *Parser) parseMultiplicativeExpression() (Expression, error) {
	left, err := p.parseUnaryExpression()
	if err != nil {
		return nil, err
	}

	for p.at().Value == "*" || p.at().Value == "/" || p.at().Value == "%" || p.at().Value == "**" {
		operator := p.eat().Value
		right, err := p.parseUnaryExpression()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Right: right, Operator: operator}
	}

	return left, nil
}

// Add support for postfix increment/decrement (x++, x--)
func (p *Parser) parseUnaryExpression() (Expression, error) {
	// Prefix unary
	if p.at().Type == NEGATION_OP || p.at().Value == "+" || p.at().Value == "-" ||
		p.at().Type == INCREMENT || p.at().Type == DECREMENT {
		operator := p.eat().Value
		value, err := p.parseUnaryExpression()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{Value: value, Operator: operator}, nil
	}

	// Parse primary/call/member first
	expr, err := p.parseCallMemberExpression()
	if err != nil {
		return nil, err
	}

	// Postfix unary (x++ or x--)
	if p.at().Type == INCREMENT || p.at().Type == DECREMENT {
		operator := p.eat().Value
		return &UnaryExpr{Value: expr, Operator: operator + "_post"}, nil
	}

	return expr, nil
}

func (p *Parser) parseCallMemberExpression() (Expression, error) {
	member, err := p.parseMemberExpression()
	if err != nil {
		return nil, err
	}

	if p.at().Type == OPEN_PAREN {
		return p.parseCallExpression(member)
	}

	return member, nil
}

func (p *Parser) parseCallExpression(caller Expression) (Expression, error) {
	callExpr := &CallExpr{Caller: caller, Args: []Expression{}}

	p.eat() // consume (
	if p.at().Type != CLOSE_PAREN {
		for {
			arg, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			callExpr.Args = append(callExpr.Args, arg)

			if p.at().Type == COMMA {
				p.eat()
			} else {
				break
			}
		}
	}

	if p.at().Type != CLOSE_PAREN {
		return nil, fmt.Errorf("expected ')' after function arguments")
	}
	p.eat() // consume )

	// Handle chained calls
	if p.at().Type == OPEN_PAREN {
		return p.parseCallExpression(callExpr)
	}

	return callExpr, nil
}

func (p *Parser) parseMemberExpression() (Expression, error) {
	object, err := p.parsePrimaryExpression()
	if err != nil {
		return nil, err
	}

	for p.at().Type == DOT || p.at().Type == OPEN_BRACKET {
		if p.at().Type == DOT {
			p.eat() // consume .
			property, err := p.parsePrimaryExpression()
			if err != nil {
				return nil, err
			}
			object = &MemberExpr{Object: object, Property: property, Computed: false}
		} else {
			p.eat() // consume [
			property, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if p.at().Type != CLOSE_BRACKET {
				return nil, fmt.Errorf("expected ']' after computed member access")
			}
			p.eat() // consume ]
			object = &MemberExpr{Object: object, Property: property, Computed: true}
		}
	}

	return object, nil
}

func (p *Parser) parsePrimaryExpression() (Expression, error) {
	token := p.at()

	switch token.Type {
	case IDENTIFIER:
		return &Identifier{Value: p.eat().Value}, nil

	case INT:
		value, err := strconv.ParseFloat(p.eat().Value, 64)
		if err != nil {
			return nil, err
		}
		return &NumericLiteral{Value: value}, nil

	case FLOAT:
		value, err := strconv.ParseFloat(p.eat().Value, 64)
		if err != nil {
			return nil, err
		}
		return &NumericLiteral{Value: value}, nil

	case STRING:
		return &StringLiteral{Value: p.eat().Value}, nil

	case BOOLEAN:
		value := p.eat().Value == "true"
		return &BooleanLiteral{Value: value}, nil

	case UNDEFINED:
		p.eat()
		return &UndefinedLiteral{}, nil

	case TYPEOF:
		p.eat()
		value, err := p.parseUnaryExpression()
		if err != nil {
			return nil, err
		}
		return &TypeofExpr{Value: value}, nil

	case OPEN_PAREN:
		p.eat() // consume (
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.at().Type != CLOSE_PAREN {
			return nil, fmt.Errorf("expected ')' after expression")
		}
		p.eat() // consume )
		return expr, nil

	case OPEN_BRACKET:
		return p.parseArrayLiteral()

	case OPEN_BRACE:
		return p.parseObjectLiteral()

	case FN, LAMBDA:
		return p.parseFunctionExpression()

	default:
		return nil, fmt.Errorf("unexpected token: %v", token.Value)
	}
}

func (p *Parser) parseArrayLiteral() (Expression, error) {
	p.eat() // consume [
	elements := []Expression{}

	if p.at().Type != CLOSE_BRACKET {
		for {
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			elements = append(elements, expr)

			if p.at().Type == COMMA {
				p.eat()
			} else {
				break
			}
		}
	}

	if p.at().Type != CLOSE_BRACKET {
		return nil, fmt.Errorf("expected ']' after array elements")
	}
	p.eat() // consume ]

	return &ArrayLiteral{Elements: elements}, nil
}

func (p *Parser) parseObjectLiteral() (Expression, error) {
	p.eat() // consume {
	properties := []Property{}

	if p.at().Type != CLOSE_BRACE {
		for {
			if p.at().Type != IDENTIFIER && p.at().Type != STRING {
				return nil, fmt.Errorf("expected property name")
			}
			key := p.eat().Value

			// Support shorthand property syntax: { x, y } instead of { x: x, y: y }
			if p.at().Type == COMMA || p.at().Type == CLOSE_BRACE {
				// Shorthand property
				properties = append(properties, Property{Key: key, Value: &Identifier{Value: key}})
			} else {
				if p.at().Type != COLON {
					return nil, fmt.Errorf("expected ':' after property name")
				}
				p.eat() // consume :

				value, err := p.parseExpression()
				if err != nil {
					return nil, err
				}

				properties = append(properties, Property{Key: key, Value: value})
			}

			if p.at().Type == COMMA {
				p.eat()
			} else {
				break
			}
		}
	}

	if p.at().Type != CLOSE_BRACE {
		return nil, fmt.Errorf("expected '}' after object properties")
	}
	p.eat() // consume }

	return &ObjectLiteral{Properties: properties}, nil
}

// Update parseFunctionExpression to handle fn:: syntax
func (p *Parser) parseFunctionExpression() (Expression, error) {
	isLambda := p.at().Type == LAMBDA
	p.eat() // consume fn or lambda

	name := ""

	// Check for anonymous function syntax: fn: or fn::
	if p.at().Type == COLON {
		p.eat() // consume :

		// Check for direct call syntax: fn::
		if p.at().Type == COLON {
			p.eat() // consume second :
			// Parse the expression to call immediately
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			// Create anonymous function that returns the expression and call it immediately
			body := []Statement{&ReturnExpr{Value: expr}}
			fn := &FunctionDeclaration{
				Name:       "",
				Parameters: []Parameter{},
				Body:       body,
				Export:     false,
			}
			// Return a call expression
			return &CallExpr{Caller: fn, Args: []Expression{}}, nil
		}

		// Parse parameters for fn: syntax
		parameters, err := p.parseParameterList()
		if err != nil {
			return nil, err
		}

		if p.at().Type != COLON {
			return nil, p.formatError("expected ':' after function parameters in anonymous function", p.at())
		}
		p.eat() // consume :

		// Parse the expression body
		expr, err := p.parseTernaryExpression()
		if err != nil {
			return nil, err
		}
		body := []Statement{&ReturnExpr{Value: expr}}

		return &FunctionDeclaration{
			Name:       "",
			Parameters: parameters,
			Body:       body,
			Export:     false,
		}, nil
	}

	// Regular function syntax
	if !isLambda && p.at().Type == IDENTIFIER {
		name = p.eat().Value
	}

	parameters, err := p.parseParameterList()
	if err != nil {
		return nil, err
	}

	if p.at().Type != OPEN_BRACE && p.at().Type != COLON {
		return nil, p.formatError("expected '{' or ':' after function parameters", p.at())
	}

	var body []Statement
	if p.at().Type == OPEN_BRACE {
		p.eat() // consume {
		for p.at().Type != CLOSE_BRACE && !p.isEOF() {
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				body = append(body, stmt)
			}
		}
		if p.at().Type != CLOSE_BRACE {
			return nil, p.formatError("expected '}' after function body", p.at())
		}
		p.eat() // consume }
	} else {
		p.eat() // consume :
		// Parse the full expression including ternary
		expr, err := p.parseTernaryExpression()
		if err != nil {
			return nil, err
		}
		body = []Statement{&ReturnExpr{Value: expr}}
	}

	return &FunctionDeclaration{
		Name:       name,
		Parameters: parameters,
		Body:       body,
		Export:     false,
	}, nil
}

// Add new method to parse parameter list with defaults
func (p *Parser) parseParameterList() ([]Parameter, error) {
	var parameters []Parameter

	for p.at().Type == IDENTIFIER {
		paramName := p.eat().Value
		var defaultValue Expression

		// Check for default parameter syntax: param=(defaultValue)
		if p.at().Type == EQUALS {
			p.eat() // consume =
			if p.at().Type != OPEN_PAREN {
				return nil, p.formatError("expected '(' after '=' in default parameter", p.at())
			}
			p.eat() // consume (

			defaultExpr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			defaultValue = defaultExpr

			if p.at().Type != CLOSE_PAREN {
				return nil, p.formatError("expected ')' after default parameter value", p.at())
			}
			p.eat() // consume )
		}

		parameters = append(parameters, Parameter{
			Name:         paramName,
			DefaultValue: defaultValue,
		})
	}

	return parameters, nil
}

// Update parseFunctionDeclaration to use new parameter parsing
func (p *Parser) parseFunctionDeclaration() (Statement, error) {
	var t Token = p.eat() // consume fn/out

	var out bool = false
	if t.Type == OUT {
		out = true

		// expect fn keyword
		if p.at().Type != FN {
			return nil, p.formatError("expected 'fn' after 'out'", p.at())
		}
		p.eat() // consume fn
	}

	// Check for fn:: syntax at statement level
	if p.at().Type == COLON {
		p.eat() // consume :
		if p.at().Type == COLON {
			p.eat() // consume second :
			// Parse the expression to call immediately
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			// Create anonymous function that returns the expression and call it immediately
			body := []Statement{&ReturnExpr{Value: expr}}
			fn := &FunctionDeclaration{
				Name:       "",
				Parameters: []Parameter{},
				Body:       body,
				Export:     out,
			}
			// Return a call expression as a statement
			return fn, nil
		}
		// If only one colon, this is an error for function declaration
		return nil, p.formatError("unexpected ':' after 'fn' in function declaration", p.at())
	}

	if p.at().Type != IDENTIFIER {
		return nil, p.formatError("expected function name", p.at())
	}
	name := p.eat().Value

	parameters, err := p.parseParameterList()
	if err != nil {
		return nil, err
	}

	if p.at().Type != OPEN_BRACE && p.at().Type != COLON {
		return nil, p.formatError("expected '{' or ':' after function parameters", p.at())
	}

	var body []Statement
	if p.at().Type == OPEN_BRACE {
		p.eat() // consume {
		for p.at().Type != CLOSE_BRACE && !p.isEOF() {
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				body = append(body, stmt)
			}
		}
		if p.at().Type != CLOSE_BRACE {
			return nil, p.formatError("expected '}' after function body", p.at())
		}
		p.eat() // consume }
	} else {
		// Colon syntax - single expression
		p.eat() // consume :
		// Parse the full expression including ternary
		expr, err := p.parseTernaryExpression()
		if err != nil {
			return nil, err
		}
		// Wrap the expression in a return statement
		body = []Statement{&ReturnExpr{Value: expr}}
	}

	return &FunctionDeclaration{
		Name:       name,
		Parameters: parameters,
		Body:       body,
		Export:     out,
	}, nil
}

func (p *Parser) parseIfStatement() (Statement, error) {
	p.eat() // consume if

	test, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.at().Type != OPEN_BRACE && p.at().Type != COLON {
		return nil, fmt.Errorf("expected '{' or ':' after if condition")
	}

	var consequent []Statement
	if p.at().Type == OPEN_BRACE {
		p.eat() // consume {
		for p.at().Type != CLOSE_BRACE && !p.isEOF() {
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				consequent = append(consequent, stmt)
			}
		}
		if p.at().Type != CLOSE_BRACE {
			return nil, fmt.Errorf("expected '}' after if body")
		}
		p.eat() // consume }
	} else {
		p.eat() // consume :
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			consequent = []Statement{stmt}
		}
	}

	var alternate []Statement
	if p.at().Type == ELSE {
		p.eat() // consume else
		if p.at().Type == IF {
			// else if
			elseIf, err := p.parseIfStatement()
			if err != nil {
				return nil, err
			}
			alternate = []Statement{elseIf}
		} else {
			// else block
			if p.at().Type == OPEN_BRACE {
				p.eat() // consume {
				for p.at().Type != CLOSE_BRACE && !p.isEOF() {
					stmt, err := p.parseStatement()
					if err != nil {
						return nil, err
					}
					if stmt != nil {
						alternate = append(alternate, stmt)
					}
				}
				if p.at().Type != CLOSE_BRACE {
					return nil, fmt.Errorf("expected '}' after else body")
				}
				p.eat() // consume }
			} else {
				stmt, err := p.parseStatement()
				if err != nil {
					return nil, err
				}
				if stmt != nil {
					alternate = []Statement{stmt}
				}
			}
		}
	}

	return &IfStatement{
		Test:       test,
		Consequent: consequent,
		Alternate:  alternate,
	}, nil
}

func (p *Parser) parseWhileStatement() (Statement, error) {
	p.eat() // consume while

	test, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.at().Type != OPEN_BRACE && p.at().Type != COLON {
		return nil, fmt.Errorf("expected '{' or ':' after while condition")
	}

	var consequent []Statement
	if p.at().Type == OPEN_BRACE {
		p.eat() // consume {
		for p.at().Type != CLOSE_BRACE && !p.isEOF() {
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				consequent = append(consequent, stmt)
			}
		}
		if p.at().Type != CLOSE_BRACE {
			return nil, fmt.Errorf("expected '}' after while body")
		}
		p.eat() // consume }
	} else {
		p.eat() // consume :
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			consequent = []Statement{stmt}
		}
	}

	return &WhileStatement{
		Test:       test,
		Consequent: consequent,
	}, nil
}

func (p *Parser) parseForStatement() (Statement, error) {
	p.eat() // consume for

	declaration, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.at().Type != SEMICOLON {
		return nil, fmt.Errorf("expected ';' after for declaration")
	}
	p.eat() // consume ;

	test, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.at().Type != SEMICOLON {
		return nil, fmt.Errorf("expected ';' after for test")
	}
	p.eat() // consume ;

	increaser, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.at().Type != OPEN_BRACE {
		return nil, fmt.Errorf("expected '{' after for header")
	}

	p.eat() // consume {

	body := []Statement{}
	for p.at().Type != CLOSE_BRACE && !p.isEOF() {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			body = append(body, stmt)
		}
	}

	if p.at().Type != CLOSE_BRACE {
		return nil, fmt.Errorf("expected '}' after for body")
	}
	p.eat() // consume }

	return &ForStatement{
		Declaration: declaration,
		Test:        test,
		Increaser:   increaser,
		Body:        body,
	}, nil
}

func (p *Parser) parseReturnStatement() (Statement, error) {
	p.eat() // consume return

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ReturnExpr{Value: value}, nil
}

func (p *Parser) parseDebugStatement() (Statement, error) {
	p.eat() // consume debug

	props := []Expression{}
	if p.at().Type == OPEN_BRACE {
		p.eat() // consume {
		if p.at().Type != CLOSE_BRACE {
			for {
				expr, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				props = append(props, expr)

				if p.at().Type == COMMA {
					p.eat()
				} else {
					break
				}
			}
		}
		if p.at().Type != CLOSE_BRACE {
			return nil, fmt.Errorf("expected '}' after debug props")
		}
		p.eat() // consume }
	} else {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		props = []Expression{expr}
	}

	return &DebugStatement{Props: props}, nil
}

func (p *Parser) parseUseStatement() (Statement, error) {
	p.eat() // consume use

	if p.at().Type != STRING {
		return nil, fmt.Errorf("expected string after use")
	}
	path := p.eat().Value

	return &UseStatement{Path: path}, nil
}

func (p *Parser) at() Token {
	if p.position >= len(p.tokens) {
		return Token{Type: EOF, Value: "", Position: Position{}}
	}
	return p.tokens[p.position]
}

func (p *Parser) eat() Token {
	token := p.at()
	p.position++
	return token
}

func (p *Parser) isEOF() bool {
	return p.at().Type == EOF
}
