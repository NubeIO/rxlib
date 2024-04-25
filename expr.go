package rxlib

import (
	"github.com/expr-lang/expr"
)

func (inst *RuntimeImpl) ExprWithError(query string) (any, error) {
	exprEnv := map[string]interface{}{
		"objects": inst.Get(),
		"runtime": inst,
	}
	// Compile and run the expression
	program, err := expr.Compile(query, expr.Env(exprEnv))
	if err != nil {
		return nil, err
	}

	output, err := expr.Run(program, exprEnv)
	return output, err
}

func (inst *RuntimeImpl) Expr(query string) any {
	exprEnv := map[string]interface{}{
		"objects": inst.Get(),
		"runtime": inst,
	}
	// Compile and run the expression
	program, err := expr.Compile(query, expr.Env(exprEnv))
	if err != nil {
		return nil
	}

	output, err := expr.Run(program, exprEnv)
	return output
}
