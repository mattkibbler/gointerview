package migrations

import (
	"database/sql"
	"fmt"
)

type migrationFunc func(*sql.DB) error
type migration struct {
	name          string
	migrationFunc migrationFunc
}

var migrations = []migration{
	{name: "createQuestionCategoriesTable", migrationFunc: createQuestionCategoriesTable},
	{name: "createQuestionsTable", migrationFunc: createQuestionsTable},
	{name: "createAnswersTable", migrationFunc: createAnswersTable},
}

func RunMigrations(db *sql.DB) error {
	// First check if migrations table exists, create it if not
	err := checkMigrationsTableExists(db)
	if err != nil {
		return err
	}

	// Fetch list of migrations
	migrationsCount, err := getMigrationsCount(db)
	if err != nil {
		return err
	}

	newMigrationsToRun := migrations[migrationsCount:]

	for index, migration := range newMigrationsToRun {
		fmt.Println("running migration", migration.name)
		err := migration.migrationFunc(db)
		if err != nil {
			return err
		}
		if err := recordMigration(db, migration, migrationsCount+index); err != nil {
			return err
		}
	}

	return nil
}

func recordMigration(db *sql.DB, m migration, index int) error {
	_, err := db.Exec(`INSERT INTO migrations (index_number, name, performed_at) VALUES (?, ?, CURRENT_TIMESTAMP)`, index, m.name)
	return err
}

func checkMigrationsTableExists(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		index_number INTEGER PRIMARY KEY,
		name VARCHAR(250) NOT NULL,
		performed_at  TIMESTAMP
	)`)
	return err
}

func getMigrationsCount(db *sql.DB) (int, error) {
	var migrationsCount int
	query := db.QueryRow("SELECT COUNT(*) as migrationsCount FROM migrations")
	err := query.Scan(&migrationsCount)
	if err != nil {
		return 0, err
	}
	return migrationsCount, nil
}

func createQuestionsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		question TEXT NOT NULL,
		answer TEXT NOT NULL,
		category_id INTEGER,
		FOREIGN KEY (category_id)
		REFERENCES question_categories(id)
		ON DELETE CASCADE
	)`)
	return err
}

func createAnswersTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS answers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		question_id INTEGER,
		given_answer TEXT,
		is_correct BOOLEAN,
		answered_at TIMESTAMP,
		FOREIGN KEY (question_id)
		REFERENCES questions(id)
		ON DELETE CASCADE
	)`)
	return err
}

func createQuestionCategoriesTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS question_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(250)
	)`)
	return err
}
