package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectToDB() {

	db_password := os.Getenv("DB_PASSWORD")
	db, err := sql.Open("mysql", "root:"+db_password+"@tcp(localhost:3306)/auth")

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connected")
	DB = db
}
