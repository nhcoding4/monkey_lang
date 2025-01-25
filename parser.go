package main

import (
	"fmt"
	"strconv"
)

// --------------------------------------------------------------------------------------------------------------------
// Parser
// --------------------------------------------------------------------------------------------------------------------

type Parser struct {
	lexer  *Lexer
	cur    Token
	peek   Token
	errors []string
}

func newParser(lexer *Lexer) *Parser {
	parser := &Parser{lexer: lexer}
	parser.nextToken()
	parser.nextToken()

	return parser
}

// --------------------------------------------------------------------------------------------------------------------
// Precedences
// --------------------------------------------------------------------------------------------------------------------

const (
	LOWEST = iota
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

// --------------------------------------------------------------------------------------------------------------------
// Parsing fns
// --------------------------------------------------------------------------------------------------------------------

type prefixParsingFn func() Expression
type infixParsingFn func(Expression) Expression

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseProgram() *Program {
	program := &Program{}
	stmts := make([]Statement, 0)

	for p.cur.tokenType != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
		p.nextToken()
	}

	program.statements = stmts

	return program
}

// --------------------------------------------------------------------------------------------------------------------
// Parsing Stmts
// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseStatement() Statement {
	switch p.cur.tokenType {
	case LET:
		return p.parseLetStatement()
	case RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{token: p.cur}
	stmts := make([]Statement, 0)
	p.nextToken()

	for p.cur.tokenType != RBRACE && p.cur.tokenType != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
		p.nextToken()
	}

	block.statements = stmts

	return block
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{token: p.cur}
	stmt.expression = p.parseExpression(LOWEST)
	if p.peek.tokenType == SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{token: p.cur}
	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.name = &Identifier{token: p.cur, value: p.cur.literal}
	if !p.expectPeek(ASSIGN) {
		return nil
	}
	p.nextToken()

	stmt.value = p.parseExpression(LOWEST)
	if p.peek.tokenType == SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{token: p.cur}
	p.nextToken()

	stmt.value = p.parseExpression(LOWEST)
	if p.peek.tokenType == SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// --------------------------------------------------------------------------------------------------------------------
// Parse Expressions
// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseExpression(prec int) Expression {
	prefix := p.getPrefixFn(p.cur.tokenType)
	if prefix == nil {
		p.noPrefixParsingFnError(p.cur.tokenType)
		return nil
	}
	leftExpr := prefix()

	for p.peek.tokenType != SEMICOLON && prec < p.getPrec(p.peek.tokenType) {
		infix := p.getInfixFn(p.peek.tokenType)
		if infix == nil {
			return leftExpr
		}

		p.nextToken()
		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseArrayLiteral() Expression {
	array := &ArrayLiteral{token: p.cur}
	array.elements = p.parseExpressionList(RBRACKET)

	return array
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{token: p.cur, value: p.cur.tokenType == TRUE}
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseCallExpression(function Expression) Expression {
	expr := &CallExpression{token: p.cur, function: function}
	expr.arguments = p.parseExpressionList(RPAREN)
	return expr
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseCallArgs() []Expression {
	args := make([]Expression, 0)

	if p.peek.tokenType == RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peek.tokenType == COMMA {
		p.nextToken()
		p.nextToken()

		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return args
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := make([]Expression, 0)

	if p.peek.tokenType == end {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peek.tokenType == COMMA {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseFloatLiteral() Expression {
	value, err := strconv.ParseFloat(p.cur.literal, 64)
	if err != nil {
		p.numberParsingError()
		return nil
	}

	return &FloatLiteral{token: p.cur, value: value}
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseFunctionLiteral() Expression {
	funcLit := &FunctionLiteral{token: p.cur}
	if !p.expectPeek(LPAREN) {
		return nil
	}

	funcLit.parameters = p.parseFunctionParameters()
	if !p.expectPeek(LBRACE) {
		return nil
	}

	funcLit.body = p.parseBlockStatement()

	return funcLit
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseFunctionParameters() []*Identifier {
	idents := make([]*Identifier, 0)

	if p.peek.tokenType == RPAREN {
		p.nextToken()
		return idents
	}

	p.nextToken()
	ident := &Identifier{token: p.cur, value: p.cur.literal}
	idents = append(idents, ident)

	for p.peek.tokenType == COMMA {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{token: p.cur, value: p.cur.literal}
		idents = append(idents, ident)
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return idents
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseGroupedExpr() Expression {
	p.nextToken()
	expr := p.parseExpression(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}

	return expr
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{token: p.cur, value: p.cur.literal}
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseIfExpression() Expression {
	expr := &IfExpression{token: p.cur}
	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	expr.condition = p.parseExpression(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}
	if !p.expectPeek(LBRACE) {
		return nil
	}

	expr.consequence = p.parseBlockStatement()

	if p.peek.tokenType == ELSE {
		p.nextToken()
		if !p.expectPeek(LBRACE) {
			return nil
		}
		expr.alternative = p.parseBlockStatement()
	}

	return expr
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseIndexExpression(left Expression) Expression {
	expr := &IndexExpression{token: p.cur, left: left}
	p.nextToken()
	expr.index = p.parseExpression(LOWEST)

	if !p.expectPeek(RBRACKET) {
		return nil
	}

	return expr
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expr := &InfixExpression{token: p.cur, operator: p.cur.literal, left: left}
	prec := p.getPrec(p.cur.tokenType)
	p.nextToken()
	expr.right = p.parseExpression(prec)

	return expr
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseIntegerLiteral() Expression {
	value, err := strconv.ParseInt(p.cur.literal, 0, 64)
	if err != nil {
		p.numberParsingError()
		return nil
	}

	return &IntegerLiteral{token: p.cur, value: value}
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parsePrefixExpression() Expression {
	expr := &PrefixExpression{token: p.cur, operator: p.cur.literal}
	p.nextToken()
	expr.right = p.parseExpression(PREFIX)

	return expr
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{token: p.cur, value: p.cur.literal}
}

// --------------------------------------------------------------------------------------------------------------------
// Helpers
// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) getPrec(tokenType TokenType) int {
	switch tokenType {
	case EQ, NOTEQ:
		return EQUALS
	case LT, LTEQ, GT, GTEQ:
		return LESSGREATER
	case PLUS, MINUS:
		return SUM
	case SLASH, ASTERIX:
		return PRODUCT
	case LPAREN:
		return CALL
	case LBRACKET:
		return INDEX
	default:
		return LOWEST
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) getInfixFn(tokenType TokenType) infixParsingFn {
	switch tokenType {
	case PLUS, MINUS, SLASH, ASTERIX, EQ, NOTEQ, LT, LTEQ, GT, GTEQ:
		return p.parseInfixExpression
	case LPAREN:
		return p.parseCallExpression
	case LBRACKET:
		return p.parseIndexExpression
	default:
		return nil
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) getPrefixFn(tokenType TokenType) prefixParsingFn {
	switch tokenType {
	case BANG, MINUS:
		return p.parsePrefixExpression
	case FALSE, TRUE:
		return p.parseBooleanLiteral
	case FLOAT:
		return p.parseFloatLiteral
	case FUNCTION:
		return p.parseFunctionLiteral
	case IDENT:
		return p.parseIdentifier
	case IF:
		return p.parseIfExpression
	case INT:
		return p.parseIntegerLiteral
	case LBRACKET:
		return p.parseArrayLiteral
	case LPAREN:
		return p.parseGroupedExpr
	case STRING:
		return p.parseStringLiteral
	default:
		return nil
	}
}

// --------------------------------------------------------------------------------------------------------------------
// Advance tokens
// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) expectPeek(expected TokenType) bool {
	if p.peek.tokenType == expected {
		p.nextToken()
		return true
	}
	p.peekError(expected)
	return false
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.lexer.nextToken()
}

// --------------------------------------------------------------------------------------------------------------------
// Errors
// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) noPrefixParsingFnError(tokenType TokenType) {
	errMsg := fmt.Sprintf(
		"Error: no prefix parsing fn found for -> { %v }. On line %v, column %v",
		tokenType,
		p.cur.line,
		p.cur.column,
	)

	p.errors = append(p.errors, errMsg)
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) numberParsingError() {
	errMsg := fmt.Sprintf(
		"Error: could not parse -> { %v } into a number. On line %v , column %v.",
		p.cur.literal,
		p.cur.line,
		p.cur.column,
	)

	p.errors = append(p.errors, errMsg)
}

// --------------------------------------------------------------------------------------------------------------------

func (p *Parser) peekError(expected TokenType) {
	errMsg := fmt.Sprintf(
		"Error: got wrong expected type -> { %v }, wanted -> { %v }. On line: %v, column %v. ",
		p.peek.tokenType,
		expected,
		p.peek.line,
		p.peek.column,
	)

	p.errors = append(p.errors, errMsg)
}

// --------------------------------------------------------------------------------------------------------------------
