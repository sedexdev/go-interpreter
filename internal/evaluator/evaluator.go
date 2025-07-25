package evaluator

import (
	"fmt"

	"github.com/sedexdev/go-interpreter/internal/ast"
	"github.com/sedexdev/go-interpreter/internal/symbol"
)

/*
========================
Main Evaluation Function
========================
*/

// Evaluate walks the AST and evaluates each node. This function recursively
// go through each node and evaluates each statement until it arrives at a
// concrete value
func Evaluate(program ast.Node, symbolTable *symbol.SymbolTable) symbol.Symbol {
	switch node := program.(type) {
	case *ast.Program:
		return evaluateStatements(node.Statements, symbolTable)
	case *ast.VariableStatement:
		variableValue := Evaluate(node.Value, symbolTable)
		symbolTable.Set(node.Name.Value, variableValue)
		return variableValue
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression, symbolTable)
	case *ast.BlockStatement:
		return evaluateStatements(node.Statements, symbolTable)
	case *ast.IfStatement:
		return evaluateIfStatement(node, symbolTable)
	case *ast.WhileStatement:
		return evaluateWhileStatement(node, symbolTable)
	case *ast.PrintStatement:
		return evaluatePrintStatement(node, symbolTable)
	case *ast.InfixExpression:
		left := Evaluate(node.Left, symbolTable)
		right := Evaluate(node.Right, symbolTable)
		return evaluateInfix(node.Operator, left, right)
	case *ast.Identifier:
		return evaluateIdentifier(node, symbolTable)
	case *ast.Integer:
		return &symbol.Integer{Value: node.Value}
	default:
		return nil
	}
}

/*
================
Error generation
================
*/

func raiseError(format string, a ...interface{}) *symbol.Error {
	return &symbol.Error{Message: fmt.Sprintf(format, a...)}
}

/*
====================================================
Helper functions for evaluating different statements
====================================================
*/

func evaluateStatements(statements []ast.Statement, symbolTable *symbol.SymbolTable) symbol.Symbol {
	var result symbol.Symbol
	// Loop over each statement in the statements array
	for _, statement := range statements {
		result = Evaluate(statement, symbolTable)
	}
	return result
}

func evaluateIdentifier(node *ast.Identifier, symbolTable *symbol.SymbolTable) symbol.Symbol {
	// Check the symbol table to see if the identifier exists
	variableValue, ok := symbolTable.Get(node.Value)
	if !ok {
		return raiseError("Couldn't find identifier: " + node.Value)
	}
	return variableValue
}

func evaluateIfStatement(ifStatement *ast.IfStatement, symbolTable *symbol.SymbolTable) symbol.Symbol {
	// Evaluate the condition to digit (1 or 0) and then progress down the
	// appropriate branch based on the result
	condition := Evaluate(ifStatement.Condition, symbolTable)
	if condition.GetValue() == "1" {
		return Evaluate(ifStatement.FirstBranch, symbolTable)
	} else if condition.GetValue() == "0" && ifStatement.SecondBranch != nil {
		return Evaluate(ifStatement.SecondBranch, symbolTable)
	}
	// Dummy return means nothing is printed to the console when the function
	// has nothing to return
	return &symbol.Dummy{Value: ""}
}

func evaluateWhileStatement(whileStatement *ast.WhileStatement, symbolTable *symbol.SymbolTable) symbol.Symbol {
	// Keep checking the condition to make sure it is still true
	for Evaluate(whileStatement.Condition, symbolTable).GetValue() == "1" {
		Evaluate(whileStatement.Loop, symbolTable)
	}
	// Return the Dummy type when the loop has completed
	return &symbol.Dummy{Value: ""}
}

func evaluatePrintStatement(printStatement *ast.PrintStatement, symbolTable *symbol.SymbolTable) symbol.Symbol {
	for _, expression := range printStatement.Values {
		fmt.Print(Evaluate(expression, symbolTable).GetValue() + " ")
	}
	return &symbol.Dummy{Value: ""}
}

func evaluateInfix(operator string, left, right symbol.Symbol) symbol.Symbol {
	leftValue := left.(*symbol.Integer).Value
	rightValue := right.(*symbol.Integer).Value

	switch operator {
	case "+":
		return &symbol.Integer{Value: leftValue + rightValue}
	case "-":
		return &symbol.Integer{Value: leftValue - rightValue}
	case "*":
		return &symbol.Integer{Value: leftValue * rightValue}
	case "/":
		return &symbol.Integer{Value: leftValue / rightValue}
	case "%":
		return &symbol.Integer{Value: leftValue % rightValue}
	case "<":
		result := evaluateToBooleanInteger(leftValue < rightValue)
		return &symbol.Integer{Value: result}
	case ">":
		result := evaluateToBooleanInteger(leftValue > rightValue)
		return &symbol.Integer{Value: result}
	case "<=":
		result := evaluateToBooleanInteger(leftValue <= rightValue)
		return &symbol.Integer{Value: result}
	case ">=":
		result := evaluateToBooleanInteger(leftValue >= rightValue)
		return &symbol.Integer{Value: result}
	case "==":
		result := evaluateToBooleanInteger(leftValue == rightValue)
		return &symbol.Integer{Value: result}
	case "!=":
		result := evaluateToBooleanInteger(leftValue != rightValue)
		return &symbol.Integer{Value: result}
	case "&&":
		return evaluateBooleanInfix(operator, leftValue, rightValue)
	case "||":
		return evaluateBooleanInfix(operator, leftValue, rightValue)
	default:
		return nil
	}
}

/*
=========================
Evaluating boolean values
=========================
*/

func evaluateBooleanInfix(operator string, left, right int64) symbol.Symbol {

	// Set the left and right values to actual boolean values first
	// This avoids the error that arises from trying to compare int64
	// using the && and || operators
	leftBoolValue := setBooleanValue(left)
	rightBoolValue := setBooleanValue(right)

	switch operator {
	case "&&":
		result := leftBoolValue && rightBoolValue
		if result {
			return &symbol.Integer{Value: 1}
		}
		return &symbol.Integer{Value: 0}
	case "||":
		result := leftBoolValue || rightBoolValue
		if result {
			return &symbol.Integer{Value: 1}
		}
		return &symbol.Integer{Value: 0}
	default:
		return nil
	}
}

func evaluateToBooleanInteger(expression bool) int64 {
	if expression {
		return 1
	}
	return 0
}

func setBooleanValue(val int64) bool {
	if val == 1 {
		return true
	}
	return false
}
