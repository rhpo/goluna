package main

type Luna struct {
	env *Environment
}

func NewLuna(env *Environment) *Luna {
	return &Luna{env: env}
}

func (l *Luna) Tokenize(code string) ([]Token, error) {
	tokenizer := NewTokenizer(code)
	return tokenizer.Tokenize()
}

func (l *Luna) Parse(tokens []Token) (Statement, error) {
	parser := NewParser(tokens, "")
	return parser.ProduceAST()
}

func (l *Luna) Evaluate(code string) (RuntimeValue, error) {
	tokens, err := l.Tokenize(code)
	if err != nil {
		return nil, err
	}

	ast, err := l.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return l.EvaluateAST(ast)
}

func (l *Luna) EvaluateAST(ast Statement) (RuntimeValue, error) {
	return Evaluate(ast, l.env)
}
