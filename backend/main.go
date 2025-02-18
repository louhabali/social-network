package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbpath := "backend/sn/db/socialnetworks.db"
	migrationspath := "backend/sn/db/migrations/sqlite"

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(dbpath), os.ModePerm); err != nil {
		fmt.Println("failed to create directories", err)
		return
	}
	
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		fmt.Println("failed to register driver", err)
		return
	}
	defer db.Close()

	RunMigrations(db, migrationspath)

	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
		w.Write([]byte("Register API"))
	})
	http.ListenAndServe(":8080", nil)
}

func RunMigrations(db *sql.DB, migpath string) {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		fmt.Println("failed to create migration driver", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		migpath, "sqlite3", driver,
	)
	if err != nil {
		fmt.Println("failed to create migration", err)
		return
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("failed to apply migrations", err)
		return
	}

	fmt.Println("Migrations applied successfully")
}
