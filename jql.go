package rxlib

import (
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/expr-lang/expr"
)

func (inst *RuntimeImpl) QueryObjectsConfig(query string) ([]*runtime.ObjectConfig, error) {
	matches, err := inst.matchValue(query)
	return SerializeCurrentFlowArray(matches), err
}

func (inst *RuntimeImpl) QueryObjects(query string) ([]Object, error) {
	matches, err := inst.matchValue(query)
	return matches, err
}

func (inst *RuntimeImpl) matchValue(query string) ([]Object, error) {
	var matches []Object
	for _, obj := range inst.Get() {
		// Prepare environment with each object wrapped
		exprEnv := map[string]interface{}{
			"jql": obj,
		}

		// Compile and run the expression
		program, err := expr.Compile(query, expr.Env(exprEnv))
		if err != nil {
			return nil, err
		}

		output, err := expr.Run(program, exprEnv)
		if err != nil {
			continue
		}

		if result, ok := output.(bool); ok && result {
			matches = append(matches, obj)
		}
	}

	return matches, nil
}
