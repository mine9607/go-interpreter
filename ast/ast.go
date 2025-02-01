package ast

import (
	"bytes"
	"monkey/token"
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

// NOTE: Every NODE must implement the NODE interface (i.e. have a TokenLiteral method which returns the literal value of the token)

// 1 - Define a Program struct which implements the Node interface
// This will be the root node for every AST the parser produces
// Each Monkey program is a series of statements
// The statements are contained in Program.Statements (a slice of AST nodes)
// Each AST node implements the Statement interface
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// method creates a buffer and writes the return value of each statement's String() method to it
// returns the buffer as a string
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// define the LetStatment struct which implements the Node and Statement Node interfaces
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier // "x"
	Value Expression
}

func (ls *LetStatement) statementNode()       {}                          // implements Statement interface
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal } // implements Node interface
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// define the Identifier struct which implements the Node and Expression interfaces
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (ident *Identifier) expressionNode()      {}                             // implements Expression interface
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal } // implements Node interface
func (ident *Identifier) String() string       { return ident.Value }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
