package main

import (
	"math/rand"
	"time"
)

type Question struct {
	Question string
	Answer   string
}

var questions = []Question{
	{Question: "What's my name?", Answer: "Bob"},
	{Question: "Where am I?", Answer: "Here"},
	{Question: "Why is this?", Answer: "Because"},
}

func GetRandomQuestion() Question {
	// Seed the random number generator
	seed := rand.NewSource(time.Now().Unix())
	r := rand.New(seed) // initialize local pseudorandom generator

	randomIndex := r.Intn(len(questions))

	return questions[randomIndex]
}
