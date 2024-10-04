package main

import "fmt"

type PrintBlockOptions struct {
	Title   string
	Message string
}

func PrintBlock(opts PrintBlockOptions) {
	fmt.Println("")
	fmt.Println("#", opts.Title)
	if opts.Message != "" {
		fmt.Println(opts.Message)
	}
}

func PresentCommands(commands []Command) {
	for i, cmd := range commands {
		fmt.Printf("%v. %v\n", i+1, cmd.Name())
	}
}
