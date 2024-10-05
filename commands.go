package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"time"
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
		&AddAQuestionCommand{},
		AddACategoryCommand{},
	}

	fmt.Println("")

	PresentCommands(commands)
	return SelectCommand(r, commands)
}

type AskMeQuestionCommand struct {
	question Question
}

func (cmd *AskMeQuestionCommand) Prompt(db *sql.DB) error {
	return nil
}

func (cmd *AskMeQuestionCommand) Name() string {
	return "Ask me a question"
}

func (cmd AskMeQuestionCommand) GetQuestion(db *sql.DB, r *bufio.Reader) (*Question, error) {
	questionCategories, err := GetQuestionCategories(db)
	if err != nil {
		return nil, err
	}
	var question Question
	if len(questionCategories) > 0 {
		PrintBlock(PrintBlockOptions{Title: "Select a question category"})

		var categoryIds = []int{0} // add in 0 as an option, which will be the "any category" option
		fmt.Println("0. Any category")
		for _, cat := range questionCategories {
			fmt.Printf("%v. %v\n", cat.ID, cat.Name)
			categoryIds = append(categoryIds, cat.ID)
		}

		selectedNumericOption, err := SelectNumericOption(r, categoryIds)
		if err != nil {
			return nil, err
		}
		if selectedNumericOption == 0 {
			fmt.Println()
			fmt.Println("Getting a question for any category...")
			question, err = GetRandomQuestion(db)
		} else {
			fmt.Println()
			var selectedCategory QuestionCategory
			for _, cat := range questionCategories {
				if cat.ID == selectedNumericOption {
					selectedCategory = cat
				}
			}
			fmt.Printf("Getting a question for the '%v' category...\n", selectedCategory.Name)
			question, err = GetRandomQuestionForCategory(db, selectedNumericOption)
		}
		if err == sql.ErrNoRows {
			return nil, errors.New("sorry, there are no questions for this category")
		} else if err != nil {
			return nil, err
		}
	} else {
		question, err = GetRandomQuestion(db)
		if err == sql.ErrNoRows {
			return nil, errorRequiringRestart{Message: "no questions available"}
		} else if err != nil {
			return nil, err
		}
	}

	return &question, nil
}

func (cmd *AskMeQuestionCommand) HandleInput(db *sql.DB, r *bufio.Reader) (Command, error) {
	question, err := cmd.GetQuestion(db, r)
	if err != nil {
		return nil, err
	}

	PrintBlock(PrintBlockOptions{Title: question.Question, Message: "Enter your answer..."})

	answerInput, err := ReadUserInput(r)
	if err != nil {
		return nil, err
	}

	PrintBlock(PrintBlockOptions{Title: "You answered", Message: answerInput})
	time.Sleep(time.Second)
	PrintBlock(PrintBlockOptions{Title: "The answer should be:", Message: question.Answer})
	time.Sleep(time.Second)

	lastAnswer, err := GetLastAnswer(db, question.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get previous answer to question: %v", err)
	}
	if lastAnswer.GivenAnswer != "" {
		PrintBlock(PrintBlockOptions{Title: "You previously answered:", Message: lastAnswer.GivenAnswer})
		time.Sleep(time.Second)
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

	time.Sleep(time.Second)

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

type AddAQuestionCommand struct{}

func (cmd AddAQuestionCommand) Name() string {
	return "Add a question"
}

func (cmd AddAQuestionCommand) Prompt(db *sql.DB) error {
	PrintBlock(PrintBlockOptions{Title: "Add a question"})
	return nil
}

func (cmd AddAQuestionCommand) HandleInput(db *sql.DB, r *bufio.Reader) (Command, error) {
	fmt.Println("Enter question text")
	fmt.Println("")
	questionInput, err := ReadUserInput(r)
	if err != nil {
		return nil, err
	}
	fmt.Println("Enter answer text")
	fmt.Println("")
	answerInput, err := ReadUserInput(r)
	if err != nil {
		return nil, err
	}

	// Use an int pointer so we can set it to nil, resulting in the field being null in the database
	// otherwise we would set the category as 0
	var categoryId *int
	questionCategories, err := GetQuestionCategories(db)
	if err != nil {
		return nil, err
	}
	if len(questionCategories) > 0 {
		var categoryIds = []int{0} // add in 0 as an option, which will be the "no category" option
		fmt.Println("0. No category")
		for _, cat := range questionCategories {
			fmt.Printf("%v. %v\n", cat.ID, cat.Name)
			categoryIds = append(categoryIds, cat.ID)
		}

		selectedNumericOption, err := SelectNumericOption(r, categoryIds)
		if err != nil {
			return nil, err
		}
		if selectedNumericOption == 0 {
			categoryId = nil
		} else {
			categoryId = &selectedNumericOption
		}
	}

	err = CreateQuestion(db, questionInput, answerInput, categoryId)
	if err != nil {
		return nil, err
	}

	fmt.Println("")
	fmt.Println("Question added!")
	time.Sleep(time.Second)

	return StartCommand{}, nil
}

type AddACategoryCommand struct{}

func (cmd AddACategoryCommand) Name() string {
	return "Add a category"
}

func (cmd AddACategoryCommand) Prompt(db *sql.DB) error {
	PrintBlock(PrintBlockOptions{Title: "Add a category"})
	return nil
}

func (cmd AddACategoryCommand) HandleInput(db *sql.DB, r *bufio.Reader) (Command, error) {
	categoryName, err := ReadUserInput(r)
	if err != nil {
		return nil, err
	}
	err = CreateQuestionCategory(db, categoryName)
	if err != nil {
		return nil, err
	}

	fmt.Println("")
	fmt.Printf("Category '%v' added\n", categoryName)
	time.Sleep(time.Second)

	return StartCommand{}, nil
}
