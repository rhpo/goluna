package main

type Environment struct {
	parent    *Environment
	variables map[string]RuntimeValue
	constants map[string]bool
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent:    parent,
		variables: make(map[string]RuntimeValue),
		constants: make(map[string]bool),
	}
}

func (env *Environment) DeclareVar(name string, value RuntimeValue, isConstant bool) RuntimeValue {
	env.variables[name] = value
	if isConstant {
		env.constants[name] = true
	}
	return value
}

func (env *Environment) AssignVar(name string, value RuntimeValue) RuntimeValue {
	// Check if it's a constant
	if env.constants[name] {
		// For now, just return the value without error - could add error handling later
		return value
	}

	// Find the environment that contains this variable
	current := env
	for current != nil {
		if _, exists := current.variables[name]; exists {
			current.variables[name] = value
			return value
		}
		current = current.parent
	}

	// If not found, declare it in current environment
	env.variables[name] = value
	return value
}

func (env *Environment) LookupVar(name string) RuntimeValue {
	current := env
	for current != nil {
		if value, exists := current.variables[name]; exists {
			return value
		}
		current = current.parent
	}
	// Return undefined instead of panicking
	return MakeUndefined()
}

func (env *Environment) HasVar(name string) bool {
	current := env
	for current != nil {
		if _, exists := current.variables[name]; exists {
			return true
		}
		current = current.parent
	}
	return false
}
