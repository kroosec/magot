package parser

import (
	"fmt"
	"magot/ast"
	"magot/lexer"
	"strconv"
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
	checkParseErrors(t, parse)
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

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	for _, msg := range errors {
		t.Errorf("parse error: %s", msg)
	}
	t.FailNow()
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 1234;
`
	lex := lexer.New(input)
	parse := New(lex)

	program := parse.ParseProgram()
	checkParseErrors(t, parse)
	if len(program.Statements) != 3 {
		t.Fatalf("expected program.Statements of size 3, got=%d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("expected *ast.ReturnStatement, got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("expected 'return', got=%T", returnStmt.TokenLiteral())
		}
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lex := lexer.New(input)
	parse := New(lex)

	program := parse.ParseProgram()
	checkParseErrors(t, parse)

	if len(program.Statements) != 1 {
		t.Fatalf("expected program.Statements of size 1, got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier, got=%T", ident)
	}
	expected := "foobar"
	if ident.Value != expected {
		t.Errorf("Erroneous value. expected=%s, got=%s", expected, ident.Value)
	}
	if ident.TokenLiteral() != expected {
		t.Errorf("Erroneous value. expected=%s, got=%s", expected, ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lex := lexer.New(input)
	parse := New(lex)

	program := parse.ParseProgram()
	checkParseErrors(t, parse)

	if len(program.Statements) != 1 {
		t.Fatalf("expected program.Statements of size 1, got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got=%T", stmt)
	}

	intLit, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got=%T", intLit)
	}
	expected := 5
	if intLit.Value != int64(expected) {
		t.Errorf("Erroneous value. expected=%d, got=%d", expected, intLit.Value)
	}
	if intLit.TokenLiteral() != strconv.Itoa(expected) {
		t.Errorf("Erroneous value. expected=%s, got=%s", string(expected), intLit.TokenLiteral())
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expected *ast.IntegerLiteral, got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("Erroneous integ.Value expected=%d, got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("Erroneous integ.TokenLiteral() expected=%d, got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		lex := lexer.New(tt.input)
		parse := New(lex)

		program := parse.ParseProgram()
		checkParseErrors(t, parse)

		if len(program.Statements) != 1 {
			t.Fatalf("expected program.Statements of size 1, got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt not a *ast.PrefixExpression, got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("Erroneous operator expcted=%s, got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}

	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		lex := lexer.New(tt.input)
		parse := New(lex)

		program := parse.ParseProgram()
		checkParseErrors(t, parse)

		if len(program.Statements) != 1 {
			t.Fatalf("expected program.Statements of size 1, got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt not a *ast.InfixExpression, got=%T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("Erroneous operator expcted=%s, got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}

	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	infixTests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for i, tt := range infixTests {
		lex := lexer.New(tt.input)
		parse := New(lex)

		program := parse.ParseProgram()
		checkParseErrors(t, parse)

		actual := program.String()
		if actual != tt.expected {
			t.Fatalf("Case #%d: expected=%q, got=%q", i, tt.expected, actual)
		}
	}
}
