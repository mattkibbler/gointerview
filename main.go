package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mattkibbler/gointerview/apperrors"
	"github.com/mattkibbler/gointerview/commands"
	"github.com/mattkibbler/gointerview/db/migrations"
	"github.com/mattkibbler/gointerview/output"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Define flags
	dbPath := flag.String("db", "./sqlitestorage/gointerview.db", "Path to the SQLite database")
	// Parse flags
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = migrations.RunMigrations(db)
	if err != nil {
		log.Fatal("could not run migrations: ", err)
	}

	reader := bufio.NewReader(os.Stdin)

	output.PrintBlock(output.PrintBlockOptions{Title: "Welcome to gointerview!!!"})

	var activeCommand commands.Command
	var lastActiveCommand commands.Command

	activeCommand = commands.StartCommand{}

	for {
		lastActiveCommand = activeCommand
		err := activeCommand.Prompt(db)
		if err != nil {
			handleErrorFromCommand(err, &activeCommand, &lastActiveCommand)
			continue
		}

		activeCommand, err = activeCommand.HandleInput(db, reader)
		if err != nil {
			handleErrorFromCommand(err, &activeCommand, &lastActiveCommand)
			continue
		}
	}
}

func handleErrorFromCommand(err error, activeCommand *commands.Command, lastActiveCommand *commands.Command) {
	if err != nil {
		fmt.Println("")

		output.PrintError(fmt.Sprintf("Error: %v", err))
		if _, ok := err.(apperrors.ErrorRequiringRestart); ok {
			*activeCommand = commands.StartCommand{}
		} else {
			*activeCommand = *lastActiveCommand
		}

		// Pause a moment so that the user can see the error before going back to the last command
		time.Sleep(time.Second)
	}
}
