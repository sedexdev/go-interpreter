package symbol

import "fmt"

// Symbol is for creating symbols that represent
// values when evaluating the AST
type Symbol interface {
	GetType() string
	GetValue() string
}

// Integer symbol
type Integer struct {
	Value int64
}

// GetType returns the INTEGER symbol type
func (integer *Integer) GetType() string {
	return "INTEGER"
}

// GetValue returns a string representation of the
// value of an integer
func (integer *Integer) GetValue() string {
	return fmt.Sprintf("%d", integer.Value)
}

// Dummy symbol for when a function has no value to return
type Dummy struct {
	Value string
}

// GetType returns the "" Dummy symbol type
func (dummy *Dummy) GetType() string {
	return ""
}

// GetValue returns an empty string
func (dummy *Dummy) GetValue() string {
	return fmt.Sprintf(dummy.Value)
}

// SymbolTable for storing variables at evaluation - uses
// a map of key value pairs that can be accessed and updated
type SymbolTable struct {
	Table map[string]Symbol
}

// CreateSymbolTable creates a new instance of SymbolTable
func CreateSymbolTable() *SymbolTable {
	table := make(map[string]Symbol)
	return &SymbolTable{Table: table}
}

// Get will get a value from a SymbolTable instance
func (symbolTable *SymbolTable) Get(identifier string) (Symbol, bool) {
	table, ok := symbolTable.Table[identifier]
	return table, ok
}

// Set will set a value in a SymbolTable instance
func (symbolTable *SymbolTable) Set(identifier string, value Symbol) Symbol {
	symbolTable.Table[identifier] = value
	return value
}

// Error symbol stores errors that occur in evaluation
type Error struct {
	Message string
}

// GetType returns the ERROR type
func (err *Error) GetType() string {
	return "ERROR"
}

// GetValue returns the error message
func (err *Error) GetValue() string {
	return "ERROR: " + err.Message
}
