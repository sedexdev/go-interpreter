package parser

import (
	"fmt"
	"strconv"

	"github.com/sedexdev/go-interpreter/internal/ast"
	"github.com/sedexdev/go-interpreter/internal/lexer"
	"github.com/sedexdev/go-interpreter/internal/token"
)

// Constant denoting no precedence
const (
	NILPRECEDENCE = 0
)

// Map that binds operators to precedences
var opPrecedences = map[string]int{
	"*":  6,
	"/":  6,
	"%":  6,
	"+":  5,
	"-":  5,
	"<":  4,
	">":  4,
	"<=": 4,
	">=": 4,
	"==": 3,
	"!=": 3,
	"&&": 2,
	"||": 1,
}

var leftBraceCount int32

// Function types for token association
// Includes:
//   - Functions for prefix expressions
//   - Functions for infix expressions
//   - Functions for prefix statements
type (
	prefixExpFunc func() ast.Expression
	infixExpFunc  func(ast.Expression) ast.Expression
)

// Parser struct for defining the the parser elements
type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	nextToken    token.Token
	// These maps create associations between tokens and functions that
	// will be called when a token of a particular type is found
	prefixExpFuncs map[string]prefixExpFunc
	infixExpFuncs  map[string]infixExpFunc
	errors         []string
}

// CreateParser creates a new Parser object
func CreateParser(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}
	parser.setTokens()
	parser.setTokens()

	/*This program uses Vaughan Pratt parse function association to link parse functions to certain tokens */

	// Prefix tokens and their associated expression functions
	parser.prefixExpFuncs = make(map[string]prefixExpFunc)
	parser.registerPrefixExpFunc(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefixExpFunc(token.INTEGER, parser.parseInteger)
	parser.registerPrefixExpFunc("LEFTPARENTHESES", parser.parseBoundExpression)

	// Infix tokens and their associated expression functions
	parser.infixExpFuncs = make(map[string]infixExpFunc)
	parser.registerInfixExpFunc("*", parser.parseInfix)
	parser.registerInfixExpFunc("/", parser.parseInfix)
	parser.registerInfixExpFunc("%", parser.parseInfix)
	parser.registerInfixExpFunc("+", parser.parseInfix)
	parser.registerInfixExpFunc("-", parser.parseInfix)
	parser.registerInfixExpFunc("<", parser.parseInfix)
	parser.registerInfixExpFunc(">", parser.parseInfix)
	parser.registerInfixExpFunc("<=", parser.parseInfix)
	parser.registerInfixExpFunc(">=", parser.parseInfix)
	parser.registerInfixExpFunc("==", parser.parseInfix)
	parser.registerInfixExpFunc("!=", parser.parseInfix)
	parser.registerInfixExpFunc("&&", parser.parseInfix)
	parser.registerInfixExpFunc("||", parser.parseInfix)

	return parser
}

/*
=================================
This is the main parsing function
=================================
*/

// ParseProgram is the function that will parse the C-- code
func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.currentToken.Type != token.END {
		statement := parser.parseStatement(false)
		if statement != nil {
			program.Statements = append(program.Statements, statement)
			leftBraceCount = 0
		}
		parser.setTokens()
	}
	return program
}

/*
================================
Error handling for syntax errors
================================
*/

// GetErrors will return any errors that arise from invalid syntax
func (parser *Parser) GetErrors() []string {
	return parser.errors
}

func (parser *Parser) logError(tokenValue string) {
	errorMsg := fmt.Sprintf("Syntax error, didn't expect %s", tokenValue)
	parser.errors = append(parser.errors, errorMsg)
}

/*
=============================
Precedence handling functions
=============================
*/

// Check the precedence of the current token
func (parser *Parser) checkCurrentPrecedence() int {
	if result, ok := opPrecedences[parser.currentToken.Value]; ok {
		return result
	}
	return NILPRECEDENCE
}

// Check the precedence of the next token
func (parser *Parser) checkNextPrecedence() int {
	if result, ok := opPrecedences[parser.nextToken.Value]; ok {
		return result
	}
	return NILPRECEDENCE
}

/*
============================================================
Functions for registering infix/prefix functions with tokens
============================================================
*/

// Register a function to a prefix expression token
func (parser *Parser) registerPrefixExpFunc(tokenType string, prefixFn prefixExpFunc) {
	parser.prefixExpFuncs[tokenType] = prefixFn
}

// Register a function to an infix expression token
func (parser *Parser) registerInfixExpFunc(tokenValue string, infixFn infixExpFunc) {
	parser.infixExpFuncs[tokenValue] = infixFn
}

/*
==============
Parser methods
==============
*/

// Parse statements - boolean fromBlock controls flow to avoid identifiers in the
// middle of expressions being treated as variable declarations, also allows variables
// to be declared and reassigned within blocks
func (parser *Parser) parseStatement(fromBlock bool) ast.Statement {
	switch parser.currentToken.Type {
	case token.IDENTIFIER:
		if fromBlock {
			return parser.parseExpressionStatement()
		}
		return parser.parseVariableDeclaration()
	case token.IF:
		return parser.parseIfStatement()
	case token.WHILE:
		return parser.parseWhileStatement()
	case token.PRINT:
		return parser.parsePrintStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

// Parse expressions statements
func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: parser.currentToken}
	statement.Expression = parser.parseExpression(NILPRECEDENCE, false)
	return statement
}

// Parse expressions - the fromPrint boolean checks to see if the expression
// is being created as part of a Print statement. If it isn't, error checking
// is performed to check for missing operators
func (parser *Parser) parseExpression(precedence int, fromPrint bool) ast.Expression {
	prefix := parser.prefixExpFuncs[parser.currentToken.Type]

	if prefix == nil {
		// If this token has no function associated with it, log error and return nil
		parser.logError(parser.currentToken.Value)
		return nil
	}
	// Call the function returned from prefixExpFuncs
	leftExpression := prefix()

	if !fromPrint {
		// If the expression did not come from a print statement and is missing an
		// operator, log an error
		if parser.nextToken.Type == "INTEGER" || parser.nextToken.Type == "IDENTIFIER" {
			parser.logError(parser.nextToken.Value)
			return nil
		}
	}

	for precedence < parser.checkNextPrecedence() {
		infix := parser.infixExpFuncs[parser.nextToken.Value]
		if infix == nil {
			return leftExpression
		}

		parser.setTokens()
		leftExpression = infix(leftExpression)
	}
	return leftExpression
}

// Parse a block statement - a block of code after an if, else or while statement
func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	if parser.currentToken.Type == "LEFTCURLYBRACE" {
		parser.setTokens()
		// Setting a count for the number of left curly braces found allows the
		// program to handle nested block statements
		leftBraceCount++
	}

	block := &ast.BlockStatement{Token: parser.currentToken}
	block.Statements = []ast.Statement{}

	returnBlock := false

	for parser.currentToken.Type != "RIGHTCURLYBRACE" && parser.currentToken.Type != token.END {
		if returnBlock {
			return block
		}
		var statement ast.Statement
		// Check for ELSE tokens - return this block if ELSE is found
		if parser.currentToken.Type == "ELSE" {
			return block
		}
		// If a variable declaration is found inside a block, set fromBlock to false in order
		// to call parseVariableDeclaration() and create a new variable declaration
		if parser.currentToken.Type == "IDENTIFIER" && parser.nextToken.Type == "ASSIGNMENT" {
			statement = parser.parseStatement(false)
			parser.setTokens()
		} else {
			statement = parser.parseStatement(true)
		}
		if statement != nil {
			block.Statements = append(block.Statements, statement)
			// If this is the end of an inner block, shift tokens and reduce leftBraceCount,
			// then set returnBlock to true so that this block gets returned as a part of the
			// outer block
			if parser.currentToken.Type == "RIGHTCURLYBRACE" && leftBraceCount > 1 {
				parser.setTokens()
				leftBraceCount--
				returnBlock = true
			}
		} else {
			return block
		}
	}
	return block
}

// Parse infix expressions
func (parser *Parser) parseInfix(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Value,
		Left:     left,
	}

	precedence := parser.checkCurrentPrecedence()
	parser.setTokens()
	expression.Right = parser.parseExpression(precedence, false)
	return expression
}

// Parse an expression surrounded by ordering parentheses
func (parser *Parser) parseBoundExpression() ast.Expression {
	parser.setTokens()
	expression := parser.parseExpression(NILPRECEDENCE, false)

	if !parser.expectNext("RIGHTPARENTHESES") {
		parser.logError(parser.nextToken.Value)
		return nil
	}
	return expression
}

// Parse an if statement
func (parser *Parser) parseIfStatement() ast.Statement {
	statement := &ast.IfStatement{Token: parser.currentToken}

	if !parser.expectNext("LEFTPARENTHESES") {
		parser.logError(parser.nextToken.Value)
		return nil
	}

	parser.setTokens()
	statement.Condition = parser.parseExpression(NILPRECEDENCE, false)

	if statement.Condition == nil {
		parser.logError("condition of statement to be nil")
		return nil
	}

	if !parser.expectNext("RIGHTPARENTHESES") {
		parser.logError(parser.nextToken.Value)
		return nil
	}

	parser.setTokens()
	statement.FirstBranch = parser.parseBlockStatement()

	if parser.currentToken.Type == "ELSE" {
		parser.setTokens()
		statement.SecondBranch = parser.parseBlockStatement()
	}
	return statement
}

// Parse while statement
func (parser *Parser) parseWhileStatement() ast.Statement {
	statement := &ast.WhileStatement{Token: parser.currentToken}

	if !parser.expectNext("LEFTPARENTHESES") {
		parser.logError(parser.nextToken.Value)
		return nil
	}

	parser.setTokens()
	statement.Condition = parser.parseExpression(NILPRECEDENCE, false)

	if statement.Condition == nil {
		parser.logError("condition of statement to be nil")
		return nil
	}

	if !parser.expectNext("RIGHTPARENTHESES") {
		parser.logError(parser.nextToken.Value)
		return nil
	}

	parser.setTokens()
	statement.Loop = parser.parseBlockStatement()
	return statement
}

// Parse print statement
func (parser *Parser) parsePrintStatement() ast.Statement {
	statement := &ast.PrintStatement{Token: parser.currentToken}

	// In order to recursively call a function declared inside another function,
	// the function must be defined as a closure
	var setValues func([]ast.Expression) []ast.Expression

	setValues = func(values []ast.Expression) []ast.Expression {
		parser.setTokens()
		values = append(values, parser.parseExpression(NILPRECEDENCE, true))

		if parser.nextToken.Type == "COMMA" {
			parser.setTokens()
			return setValues(values)
		}
		return values
	}
	// Call setValues to start recursively adding values to statement.Values
	statement.Values = setValues(statement.Values)
	parser.setTokens()

	return statement
}

// Parse variable declarations
func (parser *Parser) parseVariableDeclaration() *ast.VariableStatement {
	varStatement := &ast.VariableStatement{Token: parser.currentToken}
	varStatement.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Value}

	if !parser.expectNext("ASSIGNMENT") {
		parser.logError(parser.nextToken.Value)
		return nil
	}

	parser.setTokens()
	varStatement.Value = parser.parseExpression(NILPRECEDENCE, false)

	return varStatement
}

// Parse identifiers
func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Value}
}

// Parse integers
func (parser *Parser) parseInteger() ast.Expression {
	integer := &ast.Integer{Token: parser.currentToken}
	// Convert the string value into an integer
	value, ok := strconv.ParseInt(parser.currentToken.Value, 0, 64)

	if ok != nil {
		errorMsg := fmt.Sprintf("Unable to parse %q as an integer", parser.currentToken.Value)
		parser.errors = append(parser.errors, errorMsg)
		return nil
	}
	integer.Value = value
	return integer
}

/*
=========================================================================
Helper methods on the Parser for parsing tokens and creating AST nodes
=========================================================================
*/

func (parser *Parser) setTokens() {
	parser.currentToken = parser.nextToken
	parser.nextToken = parser.lexer.ReadNextToken()
}

func (parser *Parser) expectNext(tokenType string) bool {
	if parser.nextToken.Type == tokenType {
		parser.setTokens()
		return true
	}
	parser.logError(parser.nextToken.Value)
	return false
}
