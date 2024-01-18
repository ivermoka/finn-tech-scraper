package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"github.com/joho/godotenv"
)

var db *sql.DB // kansje bad practise, men lettest

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func Connect() {
	connStr := fmt.Sprintf("postgresql://%s:%s@ep-twilight-cell-35826753.eu-central-1.aws.neon.tech/job-scraper?sslmode=require", os.Getenv("NEONUSER"), os.Getenv("NEONPASS"))
	fmt.Println("Connection String:", connStr)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	var version string
	if err := db.QueryRow("select version()").Scan(&version); err != nil {
		panic(err)
	}

	fmt.Printf("version=%s\n", version)

	_, err = db.Exec("TRUNCATE TABLE technology_counts")
	if err != nil {
		fmt.Println("Error emptying technology_counts", err)
	}
	_, err = db.Exec("ALTER SEQUENCE technology_counts_id_seq RESTART WITH 1;")
	if err != nil {
		fmt.Println("Error resetting technology_counts", err)
	}
}

func uploadToDB(tech string, count int) {
	if db == nil {
		fmt.Println("Database connection is not established.")
		return
	}

	stmt, err := db.Prepare("INSERT INTO technology_counts (tech, count) VALUES ($1, $2)")
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(tech, count)
	if err != nil {
		fmt.Println("Error executing SQL statement:", err)
		return
	}
}
