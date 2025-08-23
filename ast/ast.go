package ast

import (
	"bytes"
	"fmt"

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
	return fmt.Sprintf("(%s %s)", prefix.Operator, prefix.Right.String())
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

// Compile time checks
var _ Expression = (*Identifier)(nil)
var _ Expression = (*IntegerLiteral)(nil)
