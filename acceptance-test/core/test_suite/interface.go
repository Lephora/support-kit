package test_suite

type Executable interface {
	Execute() error
}

type VariableContainer interface {
	SetVar(key string, value any)
	Var() map[string]any
}

var _ Executable = (*Stage)(nil)
var _ Executable = (*Case)(nil)
var _ Executable = (*Suite)(nil)

var _ VariableContainer = (*Stage)(nil)
var _ VariableContainer = (*Case)(nil)

type StageType string

const (
	API = StageType("api")
)
