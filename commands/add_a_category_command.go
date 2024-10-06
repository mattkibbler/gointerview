package commands

import (
	"bufio"
	"database/sql"
	"fmt"

	"github.com/mattkibbler/gointerview/data"
	"github.com/mattkibbler/gointerview/input"
	"github.com/mattkibbler/gointerview/output"
)

type AddACategoryCommand struct{}

func (cmd AddACategoryCommand) Name() string {
	return "Add a category"
}

func (cmd AddACategoryCommand) Prompt(db *sql.DB) error {
	output.PrintBlock(output.PrintBlockOptions{Title: "Add a category"})
	return nil
}

func (cmd AddACategoryCommand) HandleInput(db *sql.DB, r *bufio.Reader) (Command, error) {
	categoryName, err := input.ReadUserInput(r)
	if err != nil {
		return nil, err
	}
	err = data.CreateQuestionCategory(db, categoryName)
	if err != nil {
		return nil, err
	}

	fmt.Println("")
	output.TypewriterPrint(fmt.Sprintf("Category '%v' added\n", categoryName))

	return StartCommand{}, nil
}
