package models

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func InsertData(text, fileName string) error {
	godotenv.Load()
	dsn := os.Getenv("DSN")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	query := "INSERT INTO test(text, filename) VALUES(?, ?);"
	_, err = db.Exec(query, text, fileName)
	if err != nil {
		return err
	}
	return nil
}

type Info struct {
	Id       int
	Text     string
	Filename string
}

func GetData() ([]Info, error) {
	godotenv.Load()
	dsn := os.Getenv("DSN")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Info
	for rows.Next() {
		var info Info
		err := rows.Scan(&info.Id, &info.Text, &info.Filename)
		if err != nil {
			return nil, err
		}
		data = append(data, info)
	}
	return data, nil
}

type Item struct {
	Info
	ObjectUrl string
}
