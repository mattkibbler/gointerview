package commands

import (
	"bufio"
	"database/sql"
	"fmt"

	"github.com/mattkibbler/gointerview/output"
)

type StartCommand struct {
}

func (cmd StartCommand) Name() string {
	return "Start screen"
}

func (cmd StartCommand) Prompt(db *sql.DB) error {
	output.PrintBlock(output.PrintBlockOptions{Title: "What would you like to do?", Message: "Select an option"})
	return nil
}

func (cmd StartCommand) HandleInput(db *sql.DB, r *bufio.Reader) (Command, error) {
	commands := []Command{
		&AskMeQuestionCommand{},
		&AddAQuestionCommand{},
		AddACategoryCommand{},
	}

	fmt.Println("")

	PresentCommands(commands)
	return SelectCommand(r, commands)
}
