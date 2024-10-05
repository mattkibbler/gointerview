package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mattkibbler/gointerview/db/migrations"
	_ "github.com/mattn/go-sqlite3"
)

type errorRequiringRestart struct {
	Message string
}

func (err errorRequiringRestart) Error() string {
	return fmt.Sprintf("error requiring restart: %v", err.Message)
}

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

	PrintBlock(PrintBlockOptions{Title: "Welcome to gointerview!!!"})

	var activeCommand Command
	var lastActiveCommand Command

	activeCommand = StartCommand{}

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

func handleErrorFromCommand(err error, activeCommand *Command, lastActiveCommand *Command) {
	if err != nil {
		fmt.Println("")
		fmt.Println("!!!")
		fmt.Println("Error:", err)
		fmt.Println("!!!")
		if _, ok := err.(errorRequiringRestart); ok {
			*activeCommand = StartCommand{}
		} else {
			*activeCommand = *lastActiveCommand
		}

		// Pause a moment so that the user can see the error before going back to the last command
		time.Sleep(time.Second)
	}
}
