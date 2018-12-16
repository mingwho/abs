package repl

import (
	"abs/evaluator"
	"abs/lexer"
	"abs/object"
	"abs/parser"
	"fmt"
	"io"
	"os"

	prompt "github.com/c-bata/go-prompt"
)

var env *object.Environment

func init() {
	env = object.NewEnvironment()
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}

	for _, key := range env.GetKeys() {
		s = append(s, prompt.Suggest{Text: key})
	}

	if len(d.GetWordBeforeCursor()) == 0 {
		return nil
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
}

func Start(in io.Reader, out io.Writer) {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("⧐  "),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("abs-repl"),
	)
	p.Run()
}

func executor(line string) {
	if line == "quit" {
		fmt.Printf("%s", "Adios!")
		fmt.Printf("%s", "\n")
		os.Exit(0)
	}

	if line == "help" {
		fmt.Printf("Try typing something along the lines of:")
		fmt.Printf("%s", "\n")
		fmt.Printf("%s", "\n")
		fmt.Print("  ⧐  current_date = $(date)")
		fmt.Printf("%s", "\n")
		fmt.Printf("%s", "\n")
		fmt.Print("A command should be triggered in your system. Then try printing the result of that command with:")
		fmt.Printf("%s", "\n")
		fmt.Printf("%s", "\n")
		fmt.Printf("  ⧐  current_date")
		fmt.Printf("%s", "\n")
		return
	}

	Run(line, true)
}

func Run(code string, addNewline bool) {
	l := lexer.New(code)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Printf("%s", evaluated.Inspect())

		if addNewline {
			fmt.Printf("%s", "\n")
		}
	}
}

func printParserErrors(errors []string) {
	fmt.Printf("%s", " parser errors:\n")
	for _, msg := range errors {
		fmt.Printf("%s", "\t"+msg+"\n")
	}
}
