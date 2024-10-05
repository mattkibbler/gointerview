package main

import (
	"bufio"
	"database/sql"
	"fmt"
)

type Command interface {
	Name() string
	Prompt(*sql.DB) error
	HandleInput(*sql.DB, *bufio.Reader) (Command, error)
}

type StartCommand struct {
}

func (cmd StartCommand) Name() string {
	return "Start screen"
}

func (cmd StartCommand) Prompt(db *sql.DB) error {
	PrintBlock(PrintBlockOptions{Title: "What would you like to do?", Message: "Select an option"})
	return nil
}

func (cmd StartCommand) HandleInput(db *sql.DB, r *bufio.Reader) (Command, error) {
	commands := []Command{
		&AskMeQuestionCommand{},
	}

	fmt.Println("")

	PresentCommands(commands)
	return SelectCommand(r, commands)
}

type AskMeQuestionCommand struct {
	question Question
}

func (cmd *AskMeQuestionCommand) Prompt(db *sql.DB) error {
	var err error
	cmd.question, err = GetRandomQuestion(db)
	if err == sql.ErrNoRows {
		return errorRequiringRestart{Message: "no questions available"}
	} else if err != nil {
		return err
	}
	PrintBlock(PrintBlockOptions{Title: cmd.question.Question, Message: "Enter your answer..."})
	return nil
}

func (cmd *AskMeQuestionCommand) Name() string {
	return "Ask me a question"
}

func (cmd *AskMeQuestionCommand) HandleInput(db *sql.DB, r *bufio.Reader) (Command, error) {
	answerInput, err := ReadUserInput(r)
	if err != nil {
		return nil, err
	}

	PrintBlock(PrintBlockOptions{Title: "You answered", Message: answerInput})
	PrintBlock(PrintBlockOptions{Title: "The answer should be:", Message: cmd.question.Answer})

	lastAnswer, err := GetLastAnswer(db, cmd.question.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get previous answer to question: %v", err)
	}
	if lastAnswer.GivenAnswer != "" {
		PrintBlock(PrintBlockOptions{Title: "You previously answered:", Message: lastAnswer.GivenAnswer})
	}

	PrintBlock(PrintBlockOptions{Title: "Did you answer correctly this time? Be honest! (y/n)"})

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

	err = RecordAnswer(db, cmd.question.ID, answerInput, answerMarkedAsCorrect)
	if err != nil {
		return nil, fmt.Errorf("could not record answer: %v", err)
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

func (cmd AfterAskMeQuestionCommand) Prompt(db *sql.DB) error {
	PrintBlock(PrintBlockOptions{Title: "What would you like to do next?", Message: "Select an option"})
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
