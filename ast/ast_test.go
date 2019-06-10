package ast

import (
	"magot/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "foo"},
					Value: "foo",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "var"},
					Value: "bar",
				},
			},
		},
	}

	expected := "let foo = bar;"
	if program.String() != expected {
		t.Errorf("Erroneous program.String(), got='%s', expected='%s'", program.String(), expected)
	}
}
