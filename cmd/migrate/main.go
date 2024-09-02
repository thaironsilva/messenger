package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/thaironsilva/messenger/config"

	_ "github.com/lib/pq"
)

func start_migration(db *sql.DB) {
	query := `
	CREATE SCHEMA IF NOT EXISTS private;
	CREATE TABLE IF NOT EXISTS private.migrations (
		id uuid DEFAULT gen_random_uuid(),
		file_name VARCHAR UNIQUE NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	);
	`
	_, err := db.Exec(query)

	if err != nil {
		panic(err)
	}
}

func get_migrations(db *sql.DB) []string {
	var array []string

	rows, _ := db.Query("SELECT file_name FROM private.migrations")
	defer rows.Close()

	for rows.Next() {
		var migration_row string
		rows.Scan(&migration_row)
		array = append(array, migration_row)
	}

	return array
}

func create_new_migration_files(file_name string) {
	timestamp := time.Now().UTC().Format("20060102150405")
	if err := os.WriteFile("migrations/"+timestamp+"_"+file_name+".up.sql", []byte("-- migration up for "+file_name), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("migrations/"+timestamp+"_"+file_name+".down.sql", []byte("-- migration down for "+file_name), 0755); err != nil {
		panic(err)
	}
}

func migrate_down(l *log.Logger, db *sql.DB, stored_migrations []string) {
	files, err := os.ReadDir("migrations/")
	if err != nil {
		panic(err)
	}

	for i := len(files) - 1; i >= 0; i-- {
		file_name := strings.Split(files[i].Name(), ".")
		if file_name[1] == "down" {
			for _, migration_name := range stored_migrations {
				if file_name[0] == migration_name {
					l.Println("Reverting migration ", file_name[0], "...")
					query, _ := os.ReadFile("migrations/" + files[i].Name())
					if _, err := db.Query(string(query)); err != nil {
						l.Fatal("Failed to revert migration due to error: ", err)
					}
					db.Exec("DELETE FROM private.migrations WHERE file_name = ($1)", file_name[0])
					l.Println("Migration ", file_name[0], " reverted.")
				}
			}
		}
	}
}

func reset_db(db *sql.DB) {
	query := `
	DROP SCHEMA public CASCADE;
	DROP SCHEMA private CASCADE;
	CREATE SCHEMA public;
	GRANT ALL ON SCHEMA public TO postgres;
	GRANT ALL ON SCHEMA public TO public;
	`
	if _, err := db.Exec(query); err != nil {
		panic(err)
	}
}

func migrate_up(l *log.Logger, db *sql.DB, stored_migrations []string) {
	files, err := os.ReadDir("migrations/")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		file_name := strings.Split(file.Name(), ".")
		if file_name[1] == "up" {
			var run_migration bool = true
			for _, migration_name := range stored_migrations {
				if file_name[0] == migration_name {
					run_migration = false
					l.Println("Migration", file_name[0], "already ran.")
					break
				}
			}
			if run_migration {
				l.Println("Running migration ", file_name[0], "...")
				query, _ := os.ReadFile("migrations/" + file.Name())
				if _, err := db.Query(string(query)); err != nil {
					l.Fatal("Failed to run migration due to error: ", err)
				}
				db.Exec("INSERT INTO private.migrations (file_name) VALUES ($1)", file_name[0])
				l.Println("Migration ", file_name[0], " finished.")
			}
		}
	}
}

func main() {
	l := log.New(os.Stdout, "", log.LstdFlags)
	l.SetFlags(0)

	db := config.NewDB()

	defer db.Close()

	start_migration(db)

	stored_migrations := get_migrations(db)

	command := flag.String("command", "up", "migration command (up, down, reset)")
	flag.Parse()

	switch *command {
	case "new":
		create_new_migration_files(os.Args[2])
	case "down":
		l.Println("Running migrations DOWN")
		migrate_down(l, db, stored_migrations)
		l.Println("All migrations successfully reverted.")
	case "reset":
		l.Println("RESETING migrations")
		reset_db(db)
		l.Println("Starting migrations")
		start_migration(db)
		l.Println("Rerunning migrations")
		migrate_up(l, db, []string{})
		l.Println("All migrations are done.")
	default:
		l.Println("Running migrations UP")
		migrate_up(l, db, stored_migrations)
		l.Println("All migrations are done.")
	}
}
