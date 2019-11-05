package main

import (
	"books-list/controllers"
	"books-list/driver"
	_ "encoding/json"

	// "fmt"
	"books-list/models"
	"database/sql"
	"log"
	"net/http"
	_ "strconv"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var db *sql.DB

var books []models.Book

func init() {
	gotenv.Load("go.env")

}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type WithCORS struct {
	r *mux.Router
}

func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		res.Header().Set("Access-Control-Allow-Origin", origin)
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	// Stop here for a Preflighted OPTIONS request.
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(res, req)
}

func main() {
	db = driver.ConnectDB()
	router := mux.NewRouter()
	controller := controllers.Controller{}
	router.HandleFunc("/books", controller.GetBooks(db)).Methods("GET")
	router.HandleFunc("/book/{id}", controller.GetBook(db)).Methods("GET")
	router.HandleFunc("/books", controller.AddBook(db)).Methods("POST")
	router.HandleFunc("/books", controller.UpdateBook(db)).Methods("PUT")
	router.HandleFunc("/book/{id}", controller.RemoveBook(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8900", &WithCORS{router}))
}
