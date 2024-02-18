package ast

import (
	"github.com/nga1hte/interpreter/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "huchin"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "ruth"},
					Value: "ruth",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "mas"},
					Value: "mas",
				},

			},
		},
	}

	if program.String() != "huchin ruth = mas;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
