package server

import (
	"os"

	"github.com/gorilla/mux"
	"github.com/john-cai/book-manager/database"
)

type Server struct {
	*mux.Router
	database *database.Database
}

func NewServer() *Server {
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresAddr := os.Getenv("POSTGRES_ADDR")
	postgresDB := os.Getenv("POSTGRES_DB")
	database, err := database.New(
		postgresUser,
		postgresAddr,
		postgresDB,
	)
	if err != nil {
		panic(err)
	}

	s := &Server{
		Router:   mux.NewRouter(),
		database: database,
	}
	s.configureRoutes()

	return s
}

func (s *Server) configureRoutes() {
	s.HandleFunc("/books", s.AddBook).Methods("POST")
	s.HandleFunc("/books/{isbn}", s.ViewBook).Methods("GET")
	s.HandleFunc("/books/{isbn}", s.EditBook).Methods("PUT")
	s.HandleFunc("/books/{isbn}", s.RemoveBook).Methods("DELETE")

	s.HandleFunc("/collections", s.AddCollection).Methods("POST")
	s.HandleFunc("/collections/{collection_id}", s.ViewCollection).Methods("GET")
	s.HandleFunc("/collections/{collection_id}", s.EditCollection).Methods("PUT")
	s.HandleFunc("/collections/{collection_id}", s.RemoveCollection).Methods("DELETE")
	s.HandleFunc("/collections/{collection_id}/addbooks", s.AddBooksToCollection).Methods("POST")
	s.HandleFunc("/collections/{collection_id}/removebooks", s.RemoveBooksFromCollection).Methods("POST")
}
