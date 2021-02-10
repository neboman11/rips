package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const databaseUser = "catpalooza"
const databasePassword = "catpalooza"
const databaseHost = "www.nesbitt.rocks"
const databaseName = "catpalooza"
const databaseTable = "photos"
const sqlQuery = "SELECT * FROM " + databaseTable + " ORDER BY RAND() LIMIT 1;"

var db *sql.DB // Database connection pool.

type databaseRow struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Photo []byte `json:"photo"`
	Size  uint32 `json:"size"`
}

func main() {
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	// GETs
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/random", getRandomPicture)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Home Page!")
}

func getRandomPicture(w http.ResponseWriter, r *http.Request) {
	config := mysql.Config{
		User:                 databaseUser,
		Passwd:               databasePassword,
		Net:                  "tcp",
		Addr:                 databaseHost,
		DBName:               databaseName,
		AllowNativePasswords: true,
	}
	configString := config.FormatDSN()
	var err error
	db, err = sql.Open("mysql", configString)
	if err != nil {
		fmt.Fprintf(w, "Failed to connect to database: %s", err)
		return
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	photo, err := queryPhoto(r.Context())
	if err != nil {
		fmt.Fprintf(w, "Failed to query database: %s", err)
		return
	}

	body, err := json.Marshal(photo)
	fmt.Fprintf(w, "%s", body)
}

func queryPhoto(ctx context.Context) (databaseRow, error) {
	var photo databaseRow
	response := db.QueryRowContext(ctx, sqlQuery)
	err := response.Scan(&photo.ID, &photo.Name, &photo.Photo, &photo.Size)
	if err != nil {
		return photo, err
	}
	return photo, nil
}
