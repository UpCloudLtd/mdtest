package utils

import (
	"fmt"

	"github.com/google/cel-go/cel"
)

func EvaluateExpression(env *cel.Env, expr string) (any, error) {
	ast, issues := env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("failed to compile expression: %w", issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("failed to construct CEL program: %w", err)
	}

	out, _, err := prg.Eval(map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate expression: %w", err)
	}

	return out.Value(), nil
}
