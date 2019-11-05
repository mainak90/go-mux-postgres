package main

import (
	"encoding/json"
	// "fmt"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var db *sql.DB

var books []Book

func init() {
	gotenv.Load("go.env")

}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))

	logFatal(err)

	db, err = sql.Open("postgres", pgUrl)
	logFatal(err)
	err = db.Ping()
	logFatal(err)
	log.Println(pgUrl)

	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBooks).Methods("PUT")
	router.HandleFunc("/book/{id}", deleteBooks).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8900", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var book Book
	books := []Book{}
	rows, err := db.Query("select * from books")
	logFatal(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(&books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.Vars(r)
	bookID, _ := strconv.Atoi(params["id"])
	rows := db.QueryRow("Select * from books where id=$1", bookID)
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)
	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var bookID int
	json.NewDecoder(r.Body).Decode(&book)
	err := db.QueryRow("insert into books (title, author, year) values ($1, $2, $3) RETURNING id", &book.Title, &book.Author, &book.Year).Scan(&bookID)
	logFatal(err)
	json.NewEncoder(w).Encode(bookID)
}

func updateBooks(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	res, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, &book.ID)
	rowsAffected, err := res.RowsAffected()
	logFatal(err)
	json.NewEncoder(w).Encode(rowsAffected)
}

func deleteBooks(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID, _ := strconv.Atoi(params["id"])
	rows, err := db.Exec("delete from books where id=$1", bookID)
	rowsAffected, err := rows.RowsAffected()
	logFatal(err)
	json.NewEncoder(w).Encode(rowsAffected)
}
