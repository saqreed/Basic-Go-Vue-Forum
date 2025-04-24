package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "forum"
)

var DB *sqlx.DB

func InitDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Устанавливаем максимальное количество открытых соединений
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	log.Println("Successfully connected to database!")
	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
