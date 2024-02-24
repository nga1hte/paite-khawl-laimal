package ast

import (
	"testing"

	"github.com/nga1hte/paite-khawl-laimal/token"
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
