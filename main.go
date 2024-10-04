package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	PrintBlock(PrintBlockOptions{Title: "Welcome to gointerview!!!"})

	var activeCommand Command
	activeCommand = StartCommand{}

	for {
		activeCommand.Prompt()

		var err error
		activeCommand, err = activeCommand.HandleInput(reader)
		if err != nil {
			fmt.Println("error handling input:", err)
			return
		}
	}
}
