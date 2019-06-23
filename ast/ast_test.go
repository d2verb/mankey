package ast

import (
	"testing"

	"github.com/d2verb/monkey/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{Type: token.IDENT, Literal: "myVar"},
				Expression: &InfixExpression{
					Token: token.Token{Type: token.EQ, Literal: "="},
					Left: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "myVar"},
						Value: "myVar",
					},
					Operator: "=",
					Right: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
						Value: "anotherVar",
					},
				},
			},
		},
	}

	if program.String() != "(myVar = anotherVar)" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
