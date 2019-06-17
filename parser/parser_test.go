package parser

import (
	"fmt"
	"magot/ast"
	"magot/lexer"
	"strconv"
	"testing"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		value interface{}
	}{
		{"return 5;", 5},
		{"return false;", false},
		{"return foo;", "foo"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.input, 1)
		retStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("expected ReturntStatement, got=%T", program.Statements[0])
		}

		if retStmt.TokenLiteral() != "return" {
			t.Fatalf("ReturnStatement.TokenLiteral() not return, got=%s", retStmt.TokenLiteral())
		}
		if !testLiteralExpression(t, retStmt.ReturnValue, tt.value) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	program := getProgram(t, input, 1)
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

	program := getProgram(t, input, 1)
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

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		program := getProgram(t, tt.input, 1)
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
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}

	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "foo[1 + 2]"

	program := getProgram(t, input, 1)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression, got=%T", stmt.Expression)
	}
	testIdentifier(t, indexExp.Left, "foo")
	testInfixExpression(t, indexExp.Index, 1, "+", 2)
}

func TestBooleanExpression(t *testing.T) {
	boolTests := []struct {
		input string
		value bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range boolTests {
		program := getProgram(t, tt.input, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected *ast.ExpressionStatement, got=%T", program.Statements[0])
		}
		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("expected *ast.Boolean, got=%T", stmt.Expression)
		}
		if boolean.Value != tt.value {
			t.Fatalf("expected=%t, got=%t", tt.value, boolean.Value)
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true;", true, "==", true},
		{"false == false;", false, "==", false},
		{"true != false;", true, "!=", false},
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
		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
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
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := getProgram(t, input, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression not ast.IfExpression, got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements, got=%d", len(exp.Consequence.Statements))
	}
	consq, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] not an ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consq.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative not nil, got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

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

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression not ast.IfExpression, got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements, got=%d", len(exp.Consequence.Statements))
	}
	consq, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] not an ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consq.Expression, "x") {
		return
	}
	if exp.Alternative == nil {
		t.Fatalf("exp.Alternative nil")
	}
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statements, got=%d", len(exp.Alternative.Statements))
	}
	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] not an ast.ExpressionStatement, got=%T", exp.Alternative.Statements[0])
	}
	if !testIdentifier(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	program := getProgram(t, input, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression not ast.FunctionLiteral, got=%T", stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("parameters not 2, got=%d", len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	if len(function.Body.Statements) != 1 {
		t.Fatalf("Body.Statements not 1, got=%d", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body statement not an ast.ExpressionStatement, got=%T", function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParametersParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		parse := New(lex)
		program := parse.ParseProgram()
		checkParseErrors(t, parse)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expected) {
			t.Errorf("Parameters length wrong want=%d, got=%d", len(tt.expected), len(function.Parameters))
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	program := getProgram(t, input, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression not ast.CallExpression, got=%T", stmt.Expression)
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("parameters not 3, got=%d", len(exp.Arguments))
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
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

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input      string
		identifier string
		value      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foo = bar;", "foo", "bar"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.input, 1)
		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.identifier) {
			return
		}
		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.value) {
			return
		}
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world";`

	program := getProgram(t, input, 1)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	str, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral, got=%T", stmt.Expression)
	}
	expected := "hello world"
	if str.Value != expected {
		t.Errorf("string not %s, got=%s", expected, str.Value)
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 4]"

	program := getProgram(t, input, 1)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not *ast.ArrayLiteral, got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 4)
}

func checkParseErrors(t *testing.T, p *Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	for _, msg := range errors {
		t.Errorf("parse error: %s", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	t.Helper()
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

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	t.Helper()
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier, got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("expected value=%s, got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("expected TokenLiteral()=%s, got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	t.Helper()
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("exp type %T not handled", expected)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	t.Helper()
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("expected *ast.Boolean, got=%T", exp)
		return false
	}
	if boolean.Value != value {
		t.Errorf("Erroneous value, expected=%t, got=%t", value, boolean.Value)
		return false
	}
	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("Erroneous TokenLiteral() expected=%t, got=%s", value, boolean.TokenLiteral())
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	t.Helper()
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("Operator not '%s', got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func getProgram(t *testing.T, input string, statements int) *ast.Program {
	t.Helper()
	lex := lexer.New(input)
	parse := New(lex)
	program := parse.ParseProgram()
	checkParseErrors(t, parse)
	if len(program.Statements) != statements {
		t.Fatalf("expected program.Statements of size %d, got=%d", statements, len(program.Statements))
	}
	return program
}
