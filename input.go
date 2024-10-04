package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ReadUserInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("couldn't read input: %w", err)
	}
	input = strings.TrimSpace(input)
	return input, nil
}

func SelectCommand(reader *bufio.Reader, commands []Command) (Command, error) {
	input, err := ReadUserInput(reader)
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
