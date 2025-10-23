package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/zawlinnnaing/monkey-language-in-golang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, statement := range p.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}
func (p *Program) TokenLiteral() string {
	output := ""
	for _, stmt := range p.Statements {
		output += fmt.Sprintf("%s\n", stmt.TokenLiteral())
	}
	return output
}

// Identifier implements Expression interface
type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) expressionNode() {}
func (id *Identifier) TokenLiteral() string {
	return id.Token.Literal
}
func (id *Identifier) String() string {
	return id.Value
}

// IntegerLiteral implements Expression interface
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) TokenLiteral() string {
	return b.Token.Literal
}
func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

// PrefixExpression implements Expression Interface
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (prefix *PrefixExpression) expressionNode() {}
func (prefix *PrefixExpression) TokenLiteral() string {
	return prefix.Token.Literal
}
func (prefix *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", prefix.Operator, prefix.Right.String())
}

// InfixExpression implements Expression interface
type InfixExpression struct {
	Operator string
	Left     Expression
	Right    Expression
	Token    token.Token
}

func (infix *InfixExpression) expressionNode() {}
func (infix *InfixExpression) TokenLiteral() string {
	return infix.Token.Literal
}
func (infix *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", infix.Left.String(), infix.Operator, infix.Right.String())
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	value := ""
	if ls.Value != nil {
		value = ls.Value.String()
	}
	return fmt.Sprintf("%v %v = %v;", ls.Token.Literal, ls.Name.String(), value)
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	val := ""
	if rs.ReturnValue != nil {
		val = rs.ReturnValue.String()
	}
	return fmt.Sprintf("return %v;", val)
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	val := ""
	if es.Expression != nil {
		val = es.Expression.String()
	}
	return val
}

type IfExpression struct {
	Token       token.Token // 'If' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, statement := range bs.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // "(" token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// Compile time checks
var _ Expression = (*Identifier)(nil)
var _ Expression = (*IntegerLiteral)(nil)
var _ Expression = (*PrefixExpression)(nil)
var _ Expression = (*InfixExpression)(nil)
var _ Expression = (*IfExpression)(nil)
var _ Statement = (*BlockStatement)(nil)
var _ Expression = (*FunctionLiteral)(nil)
var _ Expression = (*CallExpression)(nil)
