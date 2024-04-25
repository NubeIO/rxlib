package docgen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

type Helper struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Args        []string `json:"args"`
	Return      []string `json:"return"`
	Help        string   `json:"help"`
}

// ParseInterfaceMethods reads a Go source file and extracts methods from the specified interface.
func ParseInterfaceMethods(filename, interfaceName string) ([]Helper, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var helpers []Helper

	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != interfaceName {
				continue
			}

			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return false
			}

			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue
				}
				methodName := method.Names[0].Name
				var args, returns []string
				if funcType, ok := method.Type.(*ast.FuncType); ok {
					if funcType.Params != nil {
						for _, param := range funcType.Params.List {
							for _, name := range param.Names {
								args = append(args, name.Name+" "+astPrint(param.Type))
							}
						}
					}
					if funcType.Results != nil {
						for _, res := range funcType.Results.List {
							returns = append(returns, astPrint(res.Type))
						}
					}
				}
				description := ""
				if method.Doc != nil {
					description = strings.TrimSpace(method.Doc.Text())
				}

				helpers = append(helpers, Helper{
					Name:        methodName,
					Description: description,
					Args:        args,
					Return:      returns,
					Help:        "", // Placeholder for manual entry
				})
			}
		}
		return false
	})
	return helpers, nil
}

// astPrint returns a string representation of an AST node, used for parameter and return types.
func astPrint(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + astPrint(e.X)
	case *ast.SelectorExpr:
		return astPrint(e.X) + "." + e.Sel.Name
	case *ast.ArrayType:
		return "[]" + astPrint(e.Elt)
	case *ast.Ellipsis:
		return "..." + astPrint(e.Elt)
	default:
		return ""
	}
}
