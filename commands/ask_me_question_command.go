package commands

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mattkibbler/gointerview/apperrors"
	"github.com/mattkibbler/gointerview/data"
	"github.com/mattkibbler/gointerview/input"
	"github.com/mattkibbler/gointerview/output"
)

type AskMeQuestionCommand struct {
	question data.Question
}

func (cmd *AskMeQuestionCommand) Prompt(db *sql.DB) error {
	return nil
}

func (cmd *AskMeQuestionCommand) Name() string {
	return "Ask me a question"
}

func (cmd AskMeQuestionCommand) GetQuestion(db *sql.DB, r *bufio.Reader) (*data.Question, error) {
	questionCategories, err := data.GetQuestionCategories(db)
	if err != nil {
		return nil, err
	}
	var question data.Question
	if len(questionCategories) > 0 {
		output.PrintBlock(output.PrintBlockOptions{Title: "Select a question category"})
		fmt.Println("")
		var categoryIds = []int{0} // add in 0 as an option, which will be the "any category" option
		fmt.Println("0. Any category")
		for _, cat := range questionCategories {
			output.TypewriterPrint(fmt.Sprintf("%v. %v\n", cat.ID, cat.Name))
			categoryIds = append(categoryIds, cat.ID)
		}
		fmt.Print("\n")
		selectedNumericOption, err := input.SelectNumericOption(r, categoryIds)
		if err != nil {
			return nil, err
		}
		if selectedNumericOption == 0 {
			output.PrintBlock(output.PrintBlockOptions{Title: "Getting a question for any category..."})
			question, err = data.GetRandomQuestion(db)
		} else {
			var selectedCategory data.QuestionCategory
			for _, cat := range questionCategories {
				if cat.ID == selectedNumericOption {
					selectedCategory = cat
				}
			}
			output.PrintBlock(output.PrintBlockOptions{Title: fmt.Sprintf("Getting a question for the '%v' category...\n", selectedCategory.Name)})
			question, err = data.GetRandomQuestionForCategory(db, selectedNumericOption)
		}
		if err == sql.ErrNoRows {
			return nil, errors.New("sorry, there are no questions for this category")
		} else if err != nil {
			return nil, err
		}
	} else {
		question, err = data.GetRandomQuestion(db)
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrorRequiringRestart{Message: "no questions available"}
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

	output.TypewriterPrint(question.Question)
	fmt.Print("\n\n")
	output.TypewriterPrint("Enter your answer...")
	fmt.Print("\n\n")

	answerInput, err := input.ReadUserInput(r)
	if err != nil {
		return nil, err
	}

	output.PrintBlock(output.PrintBlockOptions{Title: "You answered", Message: answerInput})
	time.Sleep(time.Second)
	output.PrintBlock(output.PrintBlockOptions{Title: "The answer should be:", Message: question.Answer})
	time.Sleep(time.Second)

	lastAnswer, err := data.GetLastAnswer(db, question.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get previous answer to question: %v", err)
	}
	if lastAnswer.GivenAnswer != "" {
		output.PrintBlock(output.PrintBlockOptions{Title: "Your last answer was:", Message: lastAnswer.GivenAnswer})
		time.Sleep(time.Second)
	}

	output.PrintBlock(output.PrintBlockOptions{Title: "Did you answer correctly this time? Be honest! (y/n)"})

	fmt.Print("\n")

	var answerMarkedAsCorrect bool
promptLoop:
	for {
		answerCorrectInput, err := input.ReadUserInput(r)
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
			fmt.Print("\n")
			continue
		}
	}

	err = data.RecordAnswer(db, cmd.question.ID, answerInput, answerMarkedAsCorrect)
	if err != nil {
		return nil, fmt.Errorf("could not record answer: %v", err)
	}

	fmt.Println("")

	if answerMarkedAsCorrect {
		output.TypewriterPrint("Well done!")
	} else {
		output.TypewriterPrint("No worries, you'll do better next time...")
	}

	fmt.Println("")

	time.Sleep(time.Second)

	return AfterAskMeQuestionCommand{}, nil
}
