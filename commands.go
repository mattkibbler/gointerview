package main

import (
	"bufio"
	"fmt"
)

type Command interface {
	Name() string
	Prompt()
	HandleInput(*bufio.Reader) (Command, error)
}

type StartCommand struct {
}

func (cmd StartCommand) Name() string {
	return "Start screen"
}

func (cmd StartCommand) Prompt() {
	PrintBlock(PrintBlockOptions{Title: "What would you like to do?", Message: "Select an option"})
}

func (cmd StartCommand) HandleInput(r *bufio.Reader) (Command, error) {
	commands := []Command{
		&AskMeQuestionCommand{},
		&AskMeQuestionCommand{},
		&AskMeQuestionCommand{},
	}

	fmt.Println("")

	PresentCommands(commands)
	return SelectCommand(r, commands)
}

type AskMeQuestionCommand struct {
	question Question
}

func (cmd *AskMeQuestionCommand) Prompt() {
	cmd.question = GetRandomQuestion()
	PrintBlock(PrintBlockOptions{Title: cmd.question.Question, Message: "Enter your answer..."})
}

func (cmd *AskMeQuestionCommand) Name() string {
	return "Ask me a question"
}

func (cmd *AskMeQuestionCommand) HandleInput(r *bufio.Reader) (Command, error) {
	answerInput, err := ReadUserInput(r)
	if err != nil {
		return nil, err
	}

	PrintBlock(PrintBlockOptions{Title: "You answered", Message: answerInput})
	PrintBlock(PrintBlockOptions{Title: "The answer should be:", Message: cmd.question.Answer})
	PrintBlock(PrintBlockOptions{Title: "Did you answer correctly? Be honest! (y/n)"})

	var answerMarkedAsCorrect bool
promptLoop:
	for {
		answerCorrectInput, err := ReadUserInput(r)
		if err != nil {
			return nil, err
		}
		switch answerCorrectInput {
		case "y":
			answerMarkedAsCorrect = true
			break promptLoop
		case "n":
			answerMarkedAsCorrect = false
			break promptLoop
		default:
			fmt.Println("Please enter either y or n")
			continue
		}
	}

	if answerMarkedAsCorrect {
		PrintBlock(PrintBlockOptions{Title: "Well done!"})
	} else {
		PrintBlock(PrintBlockOptions{Title: "No worries, you'll do better next time..."})
	}

	return AfterAskMeQuestionCommand{}, nil
}

type AfterAskMeQuestionCommand struct{}

func (cmd AfterAskMeQuestionCommand) Name() string {
	return "What would you like to do next?"
}

func (cmd AfterAskMeQuestionCommand) Prompt() {
	PrintBlock(PrintBlockOptions{Title: "What would you like to do next?", Message: "Select an option"})
}

func (cmd AfterAskMeQuestionCommand) HandleInput(r *bufio.Reader) (Command, error) {
	commands := []Command{
		&AskMeQuestionCommand{},
		&StartCommand{},
	}
	fmt.Println("")
	PresentCommands(commands)
	return SelectCommand(r, commands)
}
