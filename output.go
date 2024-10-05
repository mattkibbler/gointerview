package main

import "fmt"

const BLUE_TEXT_STYLE = "\033[34m"
const MAGENTA_BG_STYLE = "\033[45m"
const BLACK_BG_STYLE = "\033[40m"
const RESET_TEXT_STYLE = "\033[0m"
const UNDERLINE_TEXT_STYLE = "\033[4m"
const REVERSED_TEXT_AND_BG_STYLE = "\033[7m"

type PrintBlockOptions struct {
	Title   string
	Message string
}

func PrintBlock(opts PrintBlockOptions) {
	fmt.Println("")
	fmt.Printf("%v# %v%v\n", REVERSED_TEXT_AND_BG_STYLE, opts.Title, RESET_TEXT_STYLE)
	if opts.Message != "" {
		fmt.Println(opts.Message)
	}
}

func PresentCommands(commands []Command) {
	for i, cmd := range commands {
		fmt.Printf("%v. %v\n", i+1, cmd.Name())
	}
}
