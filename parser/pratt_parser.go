package parser

import "github.com/zawlinnnaing/monkey-language-in-golang/ast"

type (
	prefixParsingFn func() ast.Expression
	infixParsingFn  func(ast.Expression) ast.Expression
)
