package parser

import (
	"magot/ast"
	"magot/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 12345;
`
	lex := lexer.New(input)
	parse := New(lex)

	program := parse.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("expected program.Statements of size 3, got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}

}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	t.Helper()
	if stmt.TokenLiteral() != "let" {
		t.Errorf("Erroneous TokenLiteral, expected=let, got=%q", stmt.TokenLiteral())
		return false
	}
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected=*LetStatement, got=%T", stmt)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("Erroneous Name.Value, expected=%s, got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("Erroneous Name.TokenLiteral(), expected=%s, got=%s", name, letStmt.Name)
		return false
	}
	return true
}
