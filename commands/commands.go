package commands

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/mattkibbler/gointerview/input"
)

type Command interface {
	Name() string
	Prompt(*sql.DB) error
	HandleInput(*sql.DB, *bufio.Reader) (Command, error)
}

func PresentCommands(commands []Command) {
	for i, cmd := range commands {
		fmt.Printf("%v. %v\n", i+1, cmd.Name())
	}
}

func SelectCommand(reader *bufio.Reader, commands []Command) (Command, error) {
	input, err := input.ReadUserInput(reader)
	if err != nil {
		return nil, err
	}

	intVal, err := strconv.Atoi(input)

	if err != nil {
		return nil, errors.New("please enter number of command")
	}

	if intVal < 0 || intVal-1 > len(commands) {
		return nil, errors.New("command not found")
	}

	selectedCommand := commands[intVal-1]

	return selectedCommand, nil
}
