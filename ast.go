package main

type NodeType string

const (
	// Statements
	PROGRAM_NODE         NodeType = "Program"
	FUNCTION_DECLARATION NodeType = "FunctionDeclaration"
	IF_STATEMENT         NodeType = "IfStatement"
	WHILE_STATEMENT      NodeType = "WhileStatement"
	FOR_STATEMENT        NodeType = "ForStatement"
	RETURN_EXPR          NodeType = "ReturnExpr"
	DEBUG_STATEMENT      NodeType = "DebugStatement"
	USE_STATEMENT        NodeType = "UseStatement"

	// Expressions
	IDENTIFIER_NODE   NodeType = "Identifier"
	NUMERIC_LITERAL   NodeType = "NumericLiteral"
	STRING_LITERAL    NodeType = "StringLiteral"
	BOOLEAN_LITERAL   NodeType = "BooleanLiteral"
	UNDEFINED_LITERAL NodeType = "UndefinedLiteral"
	NULL_LITERAL      NodeType = "NullLiteral"
	ARRAY_LITERAL     NodeType = "ArrayLiteral"
	OBJECT_LITERAL    NodeType = "ObjectLiteral"

	BINARY_EXPR            NodeType = "BinaryExpr"
	UNARY_EXPR             NodeType = "UnaryExpr"
	ASSIGNMENT_EXPR        NodeType = "AssignmentExpr"
	ACTION_ASSIGNMENT_EXPR NodeType = "ActionAssignmentExpr"
	CALL_EXPR              NodeType = "CallExpr"
	MEMBER_EXPR            NodeType = "MemberExpr"
	TERNARY_EXPR           NodeType = "TernaryExpr"
	TYPEOF_EXPR            NodeType = "TypeofExpr"

	EQUALITY_EXPR   NodeType = "EqualityExpr"
	INEQUALITY_EXPR NodeType = "InequalityExpr"
	LOGICAL_EXPR    NodeType = "LogicalExpr"
)

type Statement interface {
	Kind() NodeType
}

type Expression interface {
	Statement
}

// Program
type Program struct {
	Body []Statement
}

func (p *Program) Kind() NodeType { return PROGRAM_NODE }

// Literals
type Identifier struct {
	Value string
}

func (i *Identifier) Kind() NodeType { return IDENTIFIER_NODE }

type NumericLiteral struct {
	Value float64
}

func (n *NumericLiteral) Kind() NodeType { return NUMERIC_LITERAL }

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Kind() NodeType { return STRING_LITERAL }

type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) Kind() NodeType { return BOOLEAN_LITERAL }

type UndefinedLiteral struct{}

func (u *UndefinedLiteral) Kind() NodeType { return UNDEFINED_LITERAL }

type NullLiteral struct{}

func (n *NullLiteral) Kind() NodeType { return NULL_LITERAL }

// Complex Literals
type ArrayLiteral struct {
	Elements []Expression
}

func (a *ArrayLiteral) Kind() NodeType { return ARRAY_LITERAL }

type Property struct {
	Key   string
	Value Expression
}

type ObjectLiteral struct {
	Properties []Property
}

func (o *ObjectLiteral) Kind() NodeType { return OBJECT_LITERAL }

// Expressions
type BinaryExpr struct {
	Left     Expression
	Right    Expression
	Operator string
}

func (b *BinaryExpr) Kind() NodeType { return BINARY_EXPR }

type UnaryExpr struct {
	Value    Expression
	Operator string
}

func (u *UnaryExpr) Kind() NodeType { return UNARY_EXPR }

type AssignmentExpr struct {
	Assigne Expression
	Value   Expression
}

func (a *AssignmentExpr) Kind() NodeType { return ASSIGNMENT_EXPR }

type ActionExpr struct {
	Name string
	Args []Expression
}

type ActionAssignmentExpr struct {
	Assigne Expression
	Value   Expression
	Action  ActionExpr
}

func (a *ActionAssignmentExpr) Kind() NodeType { return ACTION_ASSIGNMENT_EXPR }

type CallExpr struct {
	Caller Expression
	Args   []Expression
}

func (c *CallExpr) Kind() NodeType { return CALL_EXPR }

type MemberExpr struct {
	Object   Expression
	Property Expression
	Computed bool
}

func (m *MemberExpr) Kind() NodeType { return MEMBER_EXPR }

type TernaryExpr struct {
	Condition  Expression
	Consequent Expression
	Alternate  Expression
}

func (t *TernaryExpr) Kind() NodeType { return TERNARY_EXPR }

type TypeofExpr struct {
	Value Expression
}

func (t *TypeofExpr) Kind() NodeType { return TYPEOF_EXPR }

type EqualityExpr struct {
	Left     Expression
	Right    Expression
	Operator string
}

func (e *EqualityExpr) Kind() NodeType { return EQUALITY_EXPR }

type InequalityExpr struct {
	Left     Expression
	Right    Expression
	Operator string
}

func (i *InequalityExpr) Kind() NodeType { return INEQUALITY_EXPR }

type LogicalExpr struct {
	Left     Expression
	Right    Expression
	Operator string
}

func (l *LogicalExpr) Kind() NodeType { return LOGICAL_EXPR }

// Add a new struct for function parameters with defaults
type Parameter struct {
	Name         string
	DefaultValue Expression
}

// Statements
// Update FunctionDeclaration to use Parameter struct
type FunctionDeclaration struct {
	Name       string
	Parameters []Parameter
	Body       []Statement
	Export     bool
}

func (f *FunctionDeclaration) Kind() NodeType { return FUNCTION_DECLARATION }

type IfStatement struct {
	Test       Expression
	Consequent []Statement
	Alternate  []Statement
}

func (i *IfStatement) Kind() NodeType { return IF_STATEMENT }

type WhileStatement struct {
	Test       Expression
	Consequent []Statement
}

func (w *WhileStatement) Kind() NodeType { return WHILE_STATEMENT }

type ForStatement struct {
	Declaration Expression
	Test        Expression
	Increaser   Expression
	Body        []Statement
}

func (f *ForStatement) Kind() NodeType { return FOR_STATEMENT }

type ReturnExpr struct {
	Value Expression
}

func (r *ReturnExpr) Kind() NodeType { return RETURN_EXPR }

type DebugStatement struct {
	Props []Expression
}

func (d *DebugStatement) Kind() NodeType { return DEBUG_STATEMENT }

type UseStatement struct {
	Path string
}

func (u *UseStatement) Kind() NodeType { return USE_STATEMENT }
