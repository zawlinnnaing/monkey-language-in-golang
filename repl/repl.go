package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/zawlinnnaing/monkey-language-in-golang/evaluator"
	"github.com/zawlinnnaing/monkey-language-in-golang/lexer"
	"github.com/zawlinnnaing/monkey-language-in-golang/object"
	"github.com/zawlinnnaing/monkey-language-in-golang/parser"
)

const PROMPT = ">>"

const MONKEY_FACE = `                                                                
                            ▓▓▓▓▓▓▓▓▓▓                          
                          ▓▓▓▓▓▓▓▓▓▓▓▓▓▓                        
                        ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓                      
                      ▓▓▓▓░░░░░░▓▓░░░░░░▓▓▓▓                    
                  ░░░░▓▓░░░░░░░░░░░░░░░░░░▓▓░░░░                
                  ░░░░▓▓░░  ██░░░░░░  ██░░▓▓░░░░                
                    ░░▓▓░░████░░░░░░████░░▓▓░░                  
                      ▓▓░░▒▒▒▒░░░░░░▒▒▒▒░░▓▓                    
                        ▓▓░░░░░░░░░░░░░░▓▓                      
                          ▓▓▓▓░░░░░░▓▓▓▓                        
                              ▓▓▓▓▓▓        ░░                  
                            ▓▓▓▓▓▓▓▓▓▓      ▓▓                  
                            ▓▓▓▓▓▓▓▓▓▓    ▓▓▓▓                  
                          ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓                    
                          ▓▓▓▓░░▓▓░░▓▓▓▓                                                                            
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "🐒 Whoops!, we ran into some error 🙈.\n")
	io.WriteString(out, "Parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()

		l := lexer.New(line)
		parser := parser.New(l)

		program := parser.ParseProgram()

		if len(parser.Errors()) > 0 {
			printParserErrors(out, parser.Errors())
			// Stop further evaluation if there are parser errors
			continue
		}
		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
