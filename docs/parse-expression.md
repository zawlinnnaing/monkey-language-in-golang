# Parse Expression flow explanation
---
description: Explanation of how parse expression work
mode: ask
model: Claude Sonnet 4
---
<user>
can you explain this logic to me?
</user>
<assistant>
This is the core expression parsing logic using **Pratt parsing** (also known as precedence climbing). Let me break down how it works:

## The `parseExpression` Function Logic

```go
func (p *Parser) parseExpression(precedence int) ast.Expression {
```

This function parses expressions with operator precedence handling. The `precedence` parameter represents the minimum precedence level we're willing to parse at this level.

### Step 1: Parse the Prefix
```go
prefixParsingFn := p.prefixParsingFns[p.currentToken.Type]
if prefixParsingFn == nil {
    errorMsg := fmt.Sprintf("No prefix parsing function for token: %s", p.currentToken.Type)
    p.errors = append(p.errors, errorMsg)
    return nil
}
left := prefixParsingFn()
```

- Gets the prefix parsing function for the current token (e.g., identifier, number, `-`, `!`)
- If no prefix function exists, it's a syntax error
- Calls the function to parse the left side of the expression

### Step 2: Handle Infix Operations
```go
for !p.peekTokenIs(token.SEMICOLON) && precedence < p.getTokenPrecedence(p.peekToken) {
    infixParsingFn := p.infixParsingFns[p.peekToken.Type]
    if infixParsingFn == nil {
        return left
    }
    p.nextToken()
    left = infixParsingFn(left)
}
```

This loop continues while:
1. The next token isn't a semicolon (end of expression)
2. **AND** the current precedence is less than the next operator's precedence

## How Precedence Works

For an expression like `1 + 2 * 3`:

1. Start with `precedence = LOWEST`
2. Parse `1` (prefix)
3. See `+` with precedence `SUM`, which is > `LOWEST`, so enter loop
4. Parse infix `+`, which recursively calls `parseExpression(SUM)`
5. In the recursive call, parse `2`, then see `*` 
6. Since `PRODUCT > SUM`, it continues and parses `2 * 3` first
7. Returns `(2 * 3)` to be the right side of `1 + (2 * 3)`

This ensures `*` binds tighter than `+` automatically through the precedence comparison!
</assistant>
<user>
can you explain how parseGroupExpress work?
</user>
<assistant>
The `parseGroupExpression` function handles parsing expressions wrapped in parentheses `()`. Let me explain how it works:

```go
func (p *Parser) parseGroupExpression() ast.Expression {
	expression := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expression
}
```

## How it Works

### 1. Entry Point
This function is called when the parser encounters a `(` token (registered as a prefix function for `token.LPAREN`).

### 2. Parse Inner Expression
```go
expression := p.parseExpression(LOWEST)
```
- **Key insight**: It doesn't advance the token first - it immediately parses whatever comes after the `(`
- Uses `LOWEST` precedence, meaning it will parse the entire expression inside the parentheses
- For example, in `(1 + 2 * 3)`, it will parse the full `1 + 2 * 3` expression

### 3. Expect Closing Parenthesis
```go
if !p.expectPeek(token.RPAREN) {
    return nil
}
```
- After parsing the inner expression, it expects to find a `)`
- If not found, returns `nil` (syntax error)

### 4. Return the Expression
```go
return expression
```
- Returns the parsed inner expression directly
- The parentheses themselves don't create a new AST node - they just control parsing precedence

## Example Walkthrough

For input `(1 + 2)`:

1. Parser sees `(` token, calls `parseGroupExpression`
2. Current token is now `1`, calls `parseExpression(LOWEST)`
3. Parses `1 + 2` as a complete infix expression
4. Expects and finds `)` 
5. Returns the `1 + 2` infix expression

## Purpose

Parentheses override operator precedence. For example:
- `1 + 2 * 3` → `1 + (2 * 3)` (multiplication first)
- `(1 + 2) * 3` → `(1 + 2) * 3` (addition first due to parentheses)

The function essentially "unwraps" the parentheses while ensuring the inner expression is parsed completely before anything outside can interfere.
</assistant>