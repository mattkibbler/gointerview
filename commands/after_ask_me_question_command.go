package commands

import (
	"bufio"
	"database/sql"
	"fmt"

	"github.com/mattkibbler/gointerview/output"
)

type AfterAskMeQuestionCommand struct{}

func (cmd AfterAskMeQuestionCommand) Name() string {
	return "What would you like to do next?"
}

func (cmd AfterAskMeQuestionCommand) Prompt(db *sql.DB) error {
	output.PrintBlock(output.PrintBlockOptions{Title: "What would you like to do next?", Message: "Select an option"})
	return nil
}

func (cmd AfterAskMeQuestionCommand) HandleInput(db *sql.DB, r *bufio.Reader) (Command, error) {
	commands := []Command{
		&AskMeQuestionCommand{},
		&StartCommand{},
	}
	fmt.Println("")
	PresentCommands(commands)
	return SelectCommand(r, commands)
}
