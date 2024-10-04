package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	PrintBlock(PrintBlockOptions{Title: "Welcome to gointerview!!!"})

	var activeCommand Command
	var lastActiveCommand Command

	activeCommand = StartCommand{}

	for {
		lastActiveCommand = activeCommand
		activeCommand.Prompt()

		var err error
		activeCommand, err = activeCommand.HandleInput(reader)
		if err != nil {
			fmt.Println("")
			fmt.Println("!!!")
			fmt.Println("Error:", err)
			fmt.Println("!!!")

			activeCommand = lastActiveCommand
			// Pause a moment so that the user can see the error before going back to the last command
			time.Sleep(time.Second)
		}
	}
}
