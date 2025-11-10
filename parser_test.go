package main

import (
	"fmt"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	testCases := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tC := range testCases {
		t.Run(tC.expectedIdentifier, func(t *testing.T) {
			stmt := program.Statements[i]
			if !testLetStatement(t, stmt, tC.expectedIdentifier) {
				return
			}
		})
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, s Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*LetStatement)
	if !ok {
		t.Errorf("s not *LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*Identifier)
	if !ok {
		t.Fatalf("exp not *Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	testCases := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{
			"!5;", "!", 5,
		},
		{
			"-15;", "-", 15,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			l := NewLexer(tC.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*PrefixExpression)
			if !ok {
				t.Fatalf("stmt.Expression is not PrefixExpression. got=%T", stmt.Expression)
			}

			if exp.Operator != tC.operator {
				t.Fatalf("exp.Operator is not '%s'. got=%s", tC.operator, exp.Operator)
			}

			if !testIntegerLiteral(t, exp.Right, tC.integerValue) {
				return
			}
		})
	}
}

func testIntegerLiteral(t *testing.T, il Expression, value int64) bool {
	integ, ok := il.(*IntegerLiteral)
	if !ok {
		t.Errorf("il not *IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	testCases := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5/ 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			l := NewLexer(tC.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*InfixExpression)
			if !ok {
				t.Fatalf("stmt.Expression is not InfixExpression. got=%T", stmt.Expression)
			}

			if !testIntegerLiteral(t, exp.Left, tC.leftValue) {
				return
			}

			if exp.Operator != tC.operator {
				t.Fatalf("exp.Operator is not '%s'. got=%s", tC.operator, exp.Operator)
			}

			if !testIntegerLiteral(t, exp.Right, tC.rightValue) {
				return
			}
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b/ c",
			"((a * b) / c)",
		},
		{
			"a + b/ c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d/ e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			l := NewLexer(tC.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			actual := program.String()

			if actual != tC.want {
				t.Errorf("want=%q, got=%q", tC.want, actual)
			}
		})
	}
}
