package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/zawlinnnaing/monkey-language-in-golang/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello, %v. Welcome to the monkey programming language! \n", user.Username)
	fmt.Println("Type a command")
	repl.Start(os.Stdin, os.Stdout)
}
