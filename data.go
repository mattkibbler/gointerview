package main

import (
	"database/sql"
)

type Question struct {
	ID       int
	Question string
	Answer   string
}

type Answer struct {
	ID          int
	GivenAnswer string
	IsCorrect   bool
	AnsweredAt  string
}

func GetRandomQuestion(db *sql.DB) (Question, error) {
	var question Question
	err := db.QueryRow(`SELECT id as ID, question as Question, answer as Answer FROM questions ORDER BY RANDOM() LIMIT 1`).Scan(&question.ID, &question.Question, &question.Answer)
	if err != nil {
		return question, err
	}
	return question, nil
}

func RecordAnswer(db *sql.DB, questionId int, answerText string, isCorrect bool) error {
	_, err := db.Exec(`INSERT INTO answers (
		question_id, given_answer, is_correct, answered_at
	) VALUES (?, ?, ?, CURRENT_TIMESTAMP)`, questionId, answerText, isCorrect)

	return err
}

func GetLastAnswer(db *sql.DB, questionId int) (Answer, error) {
	var answer Answer
	err := db.QueryRow(`SELECT
		id as ID,
		given_answer as GivenAnswer,
		is_correct as IsCorrect,
		answered_at as AnsweredAt
		FROM answers ORDER BY answered_at DESC LIMIT 1
	`).Scan(&answer.ID, &answer.GivenAnswer, &answer.IsCorrect, &answer.AnsweredAt)
	// We won't consider an empty result to be an error, just return the empty answer
	if err == sql.ErrNoRows {
		return answer, nil
	} else if err != nil {
		return answer, err
	}
	return answer, nil
}
