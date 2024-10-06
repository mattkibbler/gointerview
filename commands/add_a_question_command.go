package commands

import (
	"bufio"
	"database/sql"
	"fmt"
	"time"

	"github.com/mattkibbler/gointerview/data"
	"github.com/mattkibbler/gointerview/input"
	"github.com/mattkibbler/gointerview/output"
)

type AddAQuestionCommand struct{}

func (cmd AddAQuestionCommand) Name() string {
	return "Add a question"
}

func (cmd AddAQuestionCommand) Prompt(db *sql.DB) error {
	output.PrintBlock(output.PrintBlockOptions{Title: "Add a question"})
	return nil
}

func (cmd AddAQuestionCommand) HandleInput(db *sql.DB, r *bufio.Reader) (Command, error) {
	fmt.Println("Enter question text")
	fmt.Println("")
	questionInput, err := input.ReadUserInput(r)
	if err != nil {
		return nil, err
	}
	fmt.Println("Enter answer text")
	fmt.Println("")
	answerInput, err := input.ReadUserInput(r)
	if err != nil {
		return nil, err
	}

	// Use an int pointer so we can set it to nil, resulting in the field being null in the database
	// otherwise we would set the category as 0
	var categoryId *int
	questionCategories, err := data.GetQuestionCategories(db)
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

		selectedNumericOption, err := input.SelectNumericOption(r, categoryIds)
		if err != nil {
			return nil, err
		}
		if selectedNumericOption == 0 {
			categoryId = nil
		} else {
			categoryId = &selectedNumericOption
		}
	}

	err = data.CreateQuestion(db, questionInput, answerInput, categoryId)
	if err != nil {
		return nil, err
	}

	fmt.Println("")
	fmt.Println("Question added!")
	time.Sleep(time.Second)

	return StartCommand{}, nil
}
