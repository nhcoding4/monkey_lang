package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	env := newEnvironment()
	loadNativeBuiltins(env)

	for {
		input := takeInput()
		if input == "quit" {
			return
		}

		lexer := newLexer(input)
		parser := newParser(lexer)
		program := parser.parseProgram()

		if len(parser.errors) != 0 {
			for _, err := range parser.errors {
				fmt.Println(err)
			}
		} else {
			evaluated := eval(program, env)
			if evaluated != nil {
				fmt.Println(evaluated.inspect())
			}
		}
		fmt.Println("")
	}
}

func takeInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(input)
}
