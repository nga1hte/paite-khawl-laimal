package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nga1hte/interpreter/evaluator"
	"github.com/nga1hte/interpreter/lexer"
	"github.com/nga1hte/interpreter/object"
	"github.com/nga1hte/interpreter/parser"
)

const PROMPT = ">> "

const LOGO = `

__________        .__  __            ____  __.__                  .__    .____           .__               .__   
\______   \_____  |__|/  |_  ____   |    |/ _|  |__ _____ __  _  _|  |   |    |   _____  |__| _____ _____  |  |  
 |     ___/\__  \ |  \   __\/ __ \  |      < |  |  \\__  \\ \/ \/ /  |   |    |   \__  \ |  |/     \\__  \ |  |  
 |    |     / __ \|  ||  | \  ___/  |    |  \|   Y  \/ __ \\     /|  |__ |    |___ / __ \|  |  Y Y  \/ __ \|  |__
 |____|    (____  /__||__|  \___  > |____|__ \___|  (____  /\/\_/ |____/ |_______ (____  /__|__|_|  (____  /____/
                \/              \/          \/    \/     \/                      \/    \/         \/     \/      
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, LOGO)
	io.WriteString(out, "Woops! A diklou khat awm e!\n")
	io.WriteString(out, " Enzui in!\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
