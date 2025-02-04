package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PROMPT = ">> "

const LLAMA = string('\U0001F999')

const EYE rune = '◕'

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

// NOTE: to exit the REPL use <CTRL + D> or <CTRL + Z>

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
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
	//io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "\n"+string(LLAMA))
	io.WriteString(out, string(LLAMA))
	io.WriteString(out, string(LLAMA))
	io.WriteString(out, string(LLAMA))
	io.WriteString(out, string(LLAMA))
	io.WriteString(out, "Woops! We ran into some monkey business here!")
	io.WriteString(out, string(LLAMA))
	io.WriteString(out, string(LLAMA))
	io.WriteString(out, string(LLAMA))
	io.WriteString(out, string(LLAMA))
	io.WriteString(out, string(LLAMA)+"\n")
	io.WriteString(out, "\nParser Errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
