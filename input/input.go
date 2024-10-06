package input

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

func SelectNumericOption(reader *bufio.Reader, options []int) (int, error) {
	input, err := ReadUserInput(reader)
	if err != nil {
		return 0, err
	}

	intVal, err := strconv.Atoi(input)

	if err != nil {
		return 0, errors.New("please enter number of command")
	}

	var validOption bool
	for _, option := range options {
		if option == intVal {
			validOption = true
			break
		}
	}

	if !validOption {
		return 0, errors.New("option not found")
	}

	return intVal, nil
}
