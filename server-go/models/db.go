package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	DBAdmin  *sql.DB
	DBBuyer  *sql.DB
	DBSeller *sql.DB
)

func init() {

	godotenv.Load()
	var err error

	DBAdmin, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		os.Getenv("DB_USER_ADM"),
		os.Getenv("DB_PASS_ADM"),
		os.Getenv("DB_HOST"),
		os.Getenv("MYSQL_DB"),
	))
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}

	DBBuyer, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		os.Getenv("DB_USER_BUYER"),
		os.Getenv("DB_PASS_BUYER"),
		os.Getenv("DB_HOST"),
		os.Getenv("MYSQL_DB"),
	))
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}

	DBSeller, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		os.Getenv("DB_USER_SELLER"),
		os.Getenv("DB_PASS_SELLER"),
		os.Getenv("DB_HOST"),
		os.Getenv("MYSQL_DB"),
	))
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}

	DBAdmin.SetConnMaxLifetime(100)
	DBAdmin.SetMaxIdleConns(10)
	DBAdmin.SetMaxOpenConns(10)

	DBBuyer.SetConnMaxLifetime(100)
	DBBuyer.SetMaxIdleConns(10)
	DBBuyer.SetMaxOpenConns(10)

	DBSeller.SetConnMaxLifetime(100)
	DBSeller.SetMaxIdleConns(10)
	DBSeller.SetMaxOpenConns(10)
}
