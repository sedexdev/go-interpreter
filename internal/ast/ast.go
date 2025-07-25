package ast

import "github.com/sedexdev/go-interpreter/internal/token"

/*
==========
Interfaces
==========
*/

// Node interface for all other nodes to implement
type Node interface {
}

// Statement interface for producing statements
type Statement interface {
	Node
	statementNode()
}

// Expression interface for producing expressions
type Expression interface {
	Node
	expressionNode()
}

/*
=======================================
Structs to represent program components
=======================================
*/

// Program struct holds an array of program statements
type Program struct {
	// Slice of nodes that implement Statement
	Statements []Statement
}

// VariableStatement defines a variable declaration statement
type VariableStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (varStat *VariableStatement) statementNode() {}

// ExpressionStatement defines an expression to be evaluated
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expStat *ExpressionStatement) statementNode() {}

// BlockStatement for if/else statement and while loops
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (BlockStat *BlockStatement) statementNode() {}

// IfStatement struct to represent if/else statements
type IfStatement struct {
	Token        token.Token
	Condition    Expression
	FirstBranch  *BlockStatement
	SecondBranch *BlockStatement
}

func (ifStat *IfStatement) statementNode() {}

// WhileStatement struct to represent while loops
type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Loop      *BlockStatement
}

func (whileStat *WhileStatement) statementNode() {}

// PrintStatement struct for representing print statements
type PrintStatement struct {
	Token  token.Token
	Values []Expression
}

func (printStat *PrintStatement) statementNode() {}

// Identifier struct representing a C-- identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) expressionNode() {}

// Integer struct representing a C-- integer
type Integer struct {
	Token token.Token
	Value int64
}

func (integer *Integer) expressionNode() {}

// InfixExpression defines an infix expression to be evaluated
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (infix *InfixExpression) expressionNode() {}
