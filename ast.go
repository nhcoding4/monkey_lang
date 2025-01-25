package main

import (
	"bytes"
	"fmt"
	"strings"
)

// --------------------------------------------------------------------------------------------------------------------
// Ast Node types
// --------------------------------------------------------------------------------------------------------------------

type Node interface {
	tokenLiteral() string
	toString() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// --------------------------------------------------------------------------------------------------------------------
// Root
// --------------------------------------------------------------------------------------------------------------------

type Program struct {
	statements []Statement
}

func (p *Program) tokenLiteral() string {
	if len(p.statements) == 0 {
		return ""
	}
	return p.statements[0].tokenLiteral()
}

func (p *Program) toString() string {
	var buffer bytes.Buffer

	for _, stmt := range p.statements {
		if stmt != nil {
			buffer.WriteString(stmt.toString() + " ")
		}
	}

	return buffer.String()
}

// --------------------------------------------------------------------------------------------------------------------
// Stmts
// --------------------------------------------------------------------------------------------------------------------

type BlockStatement struct {
	token      Token
	statements []Statement
}

func (b *BlockStatement) statementNode() {}

func (b *BlockStatement) tokenLiteral() string { return b.token.literal }

func (b *BlockStatement) toString() string {
	var buffer bytes.Buffer

	for _, stmt := range b.statements {
		buffer.WriteString(stmt.toString())
	}

	return buffer.String()
}

// --------------------------------------------------------------------------------------------------------------------

type ExpressionStatement struct {
	token      Token
	expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) tokenLiteral() string {
	return e.token.literal
}

func (e *ExpressionStatement) toString() string {
	if e.expression != nil {
		return e.expression.toString()
	}
	return ""
}

// --------------------------------------------------------------------------------------------------------------------

type LetStatement struct {
	token Token
	name  *Identifier
	value Expression
}

func (l *LetStatement) statementNode() {}

func (l *LetStatement) tokenLiteral() string { return l.token.literal }

func (l *LetStatement) toString() string {
	if l.value != nil {
		return fmt.Sprintf("%v %v = %v", l.tokenLiteral(), l.name.toString(), l.value.toString())
	}

	return fmt.Sprintf("%v %v = null", l.tokenLiteral(), l.name.toString())
}

// --------------------------------------------------------------------------------------------------------------------

type ReturnStatement struct {
	token Token
	value Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) tokenLiteral() string {
	return r.token.literal
}

func (r *ReturnStatement) toString() string {
	if r.value != nil {
		return fmt.Sprintf("return %v", r.value.toString())
	}

	return "return null"
}

// --------------------------------------------------------------------------------------------------------------------
// Expressions
// --------------------------------------------------------------------------------------------------------------------

type ArrayLiteral struct {
	token    Token
	elements []Expression
}

func (a *ArrayLiteral) expressionNode() {}

func (a *ArrayLiteral) tokenLiteral() string { return a.token.literal }

func (a *ArrayLiteral) toString() string {
	var buffer bytes.Buffer
	elements := make([]string, 0)

	for _, element := range a.elements {
		elements = append(elements, element.toString())
	}

	buffer.WriteString("[")
	buffer.WriteString(strings.Join(elements, ", "))
	buffer.WriteString("]")

	return buffer.String()
}

// --------------------------------------------------------------------------------------------------------------------

type BooleanLiteral struct {
	token Token
	value bool
}

func (b *BooleanLiteral) expressionNode() {}

func (b *BooleanLiteral) tokenLiteral() string { return b.token.literal }

func (b *BooleanLiteral) toString() string {
	return b.tokenLiteral()
}

// --------------------------------------------------------------------------------------------------------------------

type CallExpression struct {
	token     Token
	function  Expression
	arguments []Expression
}

func (c *CallExpression) expressionNode() {}

func (c *CallExpression) tokenLiteral() string { return c.token.literal }

func (c *CallExpression) toString() string {
	var buffer bytes.Buffer
	args := make([]string, 0)

	for _, arg := range c.arguments {
		args = append(args, arg.toString())
	}

	buffer.WriteString(c.function.toString())
	buffer.WriteString("(")
	buffer.WriteString(strings.Join(args, ", "))
	buffer.WriteString(")")

	return buffer.String()
}

// --------------------------------------------------------------------------------------------------------------------

type FloatLiteral struct {
	token Token
	value float64
}

func (f *FloatLiteral) expressionNode() {}

func (f *FloatLiteral) tokenLiteral() string { return f.token.literal }

func (f *FloatLiteral) toString() string {
	return fmt.Sprintf("%v", f.value)
}

// --------------------------------------------------------------------------------------------------------------------

type FunctionLiteral struct {
	token      Token
	parameters []*Identifier
	body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode() {}

func (f *FunctionLiteral) tokenLiteral() string { return f.token.literal }

func (f *FunctionLiteral) toString() string {
	var buffer bytes.Buffer
	params := make([]string, 0)

	for _, param := range f.parameters {
		params = append(params, param.toString())
	}

	buffer.WriteString(f.tokenLiteral())
	buffer.WriteString("(")
	buffer.WriteString(strings.Join(params, ", "))
	buffer.WriteString(")")
	buffer.WriteString(f.body.toString())

	return buffer.String()
}

// --------------------------------------------------------------------------------------------------------------------

type Identifier struct {
	token Token
	value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) tokenLiteral() string { return i.token.literal }

func (i *Identifier) toString() string {
	return i.value
}

// --------------------------------------------------------------------------------------------------------------------

type IndexExpression struct {
	token Token
	left  Expression
	index Expression
}

func (i *IndexExpression) expressionNode() {}

func (i *IndexExpression) tokenLiteral() string { return i.token.literal }

func (i *IndexExpression) toString() string {
	return fmt.Sprintf("(%v[%v])", i.left.toString(), i.index.toString())
}

// --------------------------------------------------------------------------------------------------------------------

type IfExpression struct {
	token       Token
	condition   Expression
	consequence *BlockStatement
	alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {}

func (i *IfExpression) tokenLiteral() string { return i.token.literal }

func (i *IfExpression) toString() string {
	if i.alternative == nil {
		return fmt.Sprintf("if %v %v", i.condition.toString(), i.consequence.toString())
	}

	return fmt.Sprintf("if %v %v else %v", i.condition.toString(), i.consequence.toString(), i.alternative.toString())
}

// --------------------------------------------------------------------------------------------------------------------

type InfixExpression struct {
	token    Token
	left     Expression
	operator string
	right    Expression
}

func (i *InfixExpression) expressionNode() {}

func (i *InfixExpression) tokenLiteral() string { return i.token.literal }

func (i *InfixExpression) toString() string {
	if i.left == nil || i.right == nil {
		return "{ INVALID INFIX EXPR }"
	}

	return fmt.Sprintf("(%v %v %v)", i.left.toString(), i.operator, i.right.toString())
}

// --------------------------------------------------------------------------------------------------------------------

type IntegerLiteral struct {
	token Token
	value int64
}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) tokenLiteral() string { return i.token.literal }

func (i *IntegerLiteral) toString() string {
	return fmt.Sprintf("%v", i.value)
}

// --------------------------------------------------------------------------------------------------------------------

type StringLiteral struct {
	token Token
	value string
}

func (s *StringLiteral) expressionNode() {}

func (s *StringLiteral) tokenLiteral() string { return s.token.literal }

func (s *StringLiteral) toString() string { return s.token.literal }

// --------------------------------------------------------------------------------------------------------------------

type PrefixExpression struct {
	token    Token
	operator string
	right    Expression
}

func (p *PrefixExpression) expressionNode() {}

func (p *PrefixExpression) tokenLiteral() string { return p.token.literal }

func (p *PrefixExpression) toString() string {
	if p.right != nil {
		return fmt.Sprintf("(%v %v)", p.operator, p.right.toString())
	}

	return fmt.Sprintf("(%v)", p.operator)
}

// --------------------------------------------------------------------------------------------------------------------
