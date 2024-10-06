package output

import (
	"fmt"
	"time"
)

const BLUE_TEXT_STYLE = "\033[34m"
const MAGENTA_BG_STYLE = "\033[45m"
const BLACK_BG_STYLE = "\033[40m"
const RED_BG_STYLE = "\033[41m"
const WHITE_TEXT_STYLE = "\033[37m"
const BRIGHT_WHITE_TEXT_RED_BG_STYLE = "\033[1;37;41m"
const RESET_TEXT_STYLE = "\033[0m"
const UNDERLINE_TEXT_STYLE = "\033[4m"
const REVERSED_TEXT_AND_BG_STYLE = "\033[7m"

type PrintBlockOptions struct {
	Title   string
	Message string
}

func PrintBlock(opts PrintBlockOptions) {
	fmt.Println("")
	// Set text style
	fmt.Print(REVERSED_TEXT_AND_BG_STYLE)
	TypewriterPrint(opts.Title)
	fmt.Print("\n")
	// Clear text style
	fmt.Print(RESET_TEXT_STYLE)
	if opts.Message != "" {
		fmt.Print("\n")
		TypewriterPrint(opts.Message)
		fmt.Print("\n")
	}
}

func PrintError(message string) {
	// Set text style
	fmt.Print(BRIGHT_WHITE_TEXT_RED_BG_STYLE)
	TypewriterPrint(message)
	// Clear text style
	fmt.Print(RESET_TEXT_STYLE)
	fmt.Print("\n")
}

func TypewriterPrint(message string) {
	for i := 0; i < len(message); i++ {
		fmt.Printf("%c", message[i])
		time.Sleep(10 * time.Millisecond)
	}
}
