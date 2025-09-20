package parser

import (
	"fmt"
	"strconv"

	"github.com/zawlinnnaing/monkey-language-in-golang/ast"
	"github.com/zawlinnnaing/monkey-language-in-golang/lexer"
	"github.com/zawlinnnaing/monkey-language-in-golang/token"
)

// Important to keep the order, lower in the list means it takes greater precedence
const (
	_ int = iota
	LOWEST
	EQUALS       // == or !=
	LESS_GREATER // > OR <
	SUM          // + or -
	PRODUCT      // * or /
	PREFIX       // -X or !X
	CALL         // myFunction(x)
)

var precedencesMap = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LT:       LESS_GREATER,
	token.GT:       LESS_GREATER,
	token.LPAREN:   CALL,
}

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string

	prefixParsingFns map[token.TokenType]prefixParsingFn
	infixParsingFns  map[token.TokenType]infixParsingFn
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.currentTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Assign identifier
	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

/*
Check whether peek token is of expected type. If it is, move current token to next
*/
func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	p.peekErrors(tokenType)
	return false
}

func (p *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return p.currentToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	statement.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{
		Token: p.currentToken,
	}
	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()
	expression := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expression
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixParsingFn := p.prefixParsingFns[p.currentToken.Type]
	if prefixParsingFn == nil {
		errorMsg := fmt.Sprintf("No prefix parsing function for token: %s", p.currentToken.Type)
		p.errors = append(p.errors, errorMsg)
		return nil
	}
	left := prefixParsingFn()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.getTokenPrecedence(p.peekToken) {
		infixParsingFn := p.infixParsingFns[p.peekToken.Type]
		if infixParsingFn == nil {
			return left
		}
		p.nextToken()
		left = infixParsingFn(left)
	}
	return left
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.currentToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatements()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatements()
	}

	return expression
}

func (p *Parser) parseBlockStatements() *ast.BlockStatement {

	blockStatement := &ast.BlockStatement{
		Token:      p.currentToken,
		Statements: []ast.Statement{},
	}

	p.nextToken()
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			blockStatement.Statements = append(blockStatement.Statements, statement)
		}
		p.nextToken()
	}

	return blockStatement
}

func (p *Parser) peekErrors(expectTokenType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", expectTokenType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefixFn(token token.TokenType, fn prefixParsingFn) {
	p.prefixParsingFns[token] = fn
}

func (p *Parser) registerInfixFn(token token.TokenType, fn infixParsingFn) {
	p.infixParsingFns[token] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	return identifier
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currentToken}
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		errorMsg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, errorMsg)
		return nil
	}
	literal.Value = value
	return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Operator: p.currentToken.Literal,
		Token:    p.currentToken,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infixExpression := &ast.InfixExpression{
		Left:     left,
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	precedence := p.getTokenPrecedence(p.currentToken)
	p.nextToken()
	infixExpression.Right = p.parseExpression(precedence)
	return infixExpression
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	exp := &ast.BooleanLiteral{
		Token: p.currentToken,
		Value: p.currentTokenIs(token.TRUE),
	}
	return exp
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	functionLiteral := &ast.FunctionLiteral{
		Token: p.currentToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	functionLiteral.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	functionLiteral.Body = p.parseBlockStatements()
	return functionLiteral
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	identifiers = append(identifiers, p.parseIdentifier().(*ast.Identifier))
	for p.peekTokenIs(token.COMMA) {
		// Skips comma and reaches next parameter
		p.nextToken()
		p.nextToken()

		identifiers = append(identifiers, p.parseIdentifier().(*ast.Identifier))
	}

	if !p.expectPeek(token.RPAREN) {
		// After parsing all parameters, if right paren is not found, parse error
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currentToken, Function: function}
	exp.Arguments = p.parseArguments()
	return exp
}

func (p *Parser) parseArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) getTokenPrecedence(token token.Token) int {
	precedence, ok := precedencesMap[token.Type]
	if !ok {
		return LOWEST
	}
	return precedence
}

func (p *Parser) Errors() []string {
	return p.errors
}

func New(l *lexer.Lexer) *Parser {
	p := Parser{
		lexer:  l,
		errors: []string{},

		prefixParsingFns: make(map[token.TokenType]prefixParsingFn),
		infixParsingFns:  make(map[token.TokenType]infixParsingFn),
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefixFn(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefixFn(token.LPAREN, p.parseGroupExpression)
	p.registerPrefixFn(token.IF, p.parseIfExpression)
	p.registerPrefixFn(token.FUNCTION, p.parseFunctionLiteral)

	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.LPAREN, p.parseCallExpression)

	return &p
}
