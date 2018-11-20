package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"

	"github.com/john-cai/book-manager/models"
	"github.com/john-cai/book-manager/responder"
)

func (s *Server) AddBook(w http.ResponseWriter, r *http.Request) {
	var err error
	var book models.Book
	if err = json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not read request"))
		return
	}
	var validationErrs []responder.Error
	if validationErrs = models.ValidateBook(book); len(validationErrs) > 0 {
		if err = responder.RespondErrors(w, validationErrs, http.StatusBadRequest); err != nil {
			log.Errorf("error when responding with 400 error: %v", err)
		}
		return
	}

	if err = s.database.AddBook(&book); err != nil {
		if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
			log.Errorf("error when responding with 500 error %v", err)
		}
		return
	}

	if err = responder.RespondSingle(w, &book, http.StatusCreated); err != nil {
		log.Errorf("error when responding with 201 error %v", err)
	}
}

func (s *Server) ViewBook(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	isbn := vars["isbn"]

	var book *models.Book
	if book, err = s.database.GetBookByISBN(isbn); err != nil {
		if err == pg.ErrNoRows {
			if err = responder.RespondError(w, "", "", http.StatusNotFound); err != nil {
				log.Errorf("error when responding with 400 error: %v", err)
			}
			return
		}
		if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
			log.Errorf("error when responding with 500 error: %v", err)
		}
		return
	}
	if err = responder.RespondSingle(w, &book, http.StatusOK); err != nil {
		log.Errorf("error when responding with 200 error: ", err)
	}
}

func (s *Server) EditBook(w http.ResponseWriter, r *http.Request) {
	var err error
	var book models.Book

	if err = json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not read request"))
		return
	}
	vars := mux.Vars(r)
	isbn := vars["isbn"]
	book.ISBN = isbn

	var validationErrs []responder.Error
	if validationErrs = models.ValidateBook(book); len(validationErrs) > 0 {
		if err = responder.RespondErrors(w, validationErrs, http.StatusBadRequest); err != nil {
			log.Errorf("error when responding with 400 error: %v", err)
		}
		return
	}

	if err = s.database.UpdateBook(&book); err != nil {
		if err == pg.ErrNoRows {
			if err = responder.RespondError(w, "", "", http.StatusNotFound); err != nil {
				log.Errorf("error when responding with 400 error: %v", err)
			}
			return
		}
		log.Errorf("error when updating book: %v", err)
		if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
			log.Errorf("error when responding with 500 error %v", err)
		}
		return
	}

	if err = responder.RespondSingle(w, &book, http.StatusOK); err != nil {
		log.Errorf("error when responding with 201 error %v", err)
	}
}

func (s *Server) RemoveBook(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	isbn := vars["isbn"]

	if err = s.database.DeleteBookByISBN(isbn); err != nil {
		if err == pg.ErrNoRows {
			if err = responder.RespondError(w, "", "", http.StatusNotFound); err != nil {
				log.Errorf("error when responding with 400 error: %v", err)
			}
			return
		}

		if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
			log.Errorf("error when responding with 500 error: %v", err)
		}
		return
	}

	if err = responder.Respond(w, http.StatusOK); err != nil {
		log.Errorf("error when responding with 200 error: ", err)
	}
}

func (s *Server) AddCollection(w http.ResponseWriter, r *http.Request) {
	var err error
	var collection models.Collection
	if err = json.NewDecoder(r.Body).Decode(&collection); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not read request"))
		return
	}
	var validationErrs []responder.Error
	if validationErrs = models.ValidateCollection(collection); len(validationErrs) > 0 {
		if err = responder.RespondErrors(w, validationErrs, http.StatusBadRequest); err != nil {
			log.Errorf("error when responding with 400 error: %v", err)
		}
		return
	}

	if err = s.database.AddCollection(&collection); err != nil {
		if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
			log.Errorf("error when responding with 500 error %v", err)
		}
		return
	}

	if err = responder.RespondSingle(w, &collection, http.StatusCreated); err != nil {
		log.Errorf("error when responding with 201 error %v", err)
	}

}

func (s *Server) ViewCollection(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	collectionID, err := strconv.Atoi(vars["collection_id"])
	if err != nil {
		responder.RespondError(w, "bad value", "collection_id", http.StatusBadRequest)
		return
	}

	var collection *models.Collection
	if collection, err = s.database.GetCollectionByID(collectionID); err != nil {
		if err == pg.ErrNoRows {
			if err = responder.RespondError(w, "", "", http.StatusNotFound); err != nil {
				log.Errorf("error when responding with 400 error: %v", err)
			}
			return
		}

		if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
			log.Errorf("error when responding with 500 error: %v", err)
		}
		return
	}

	if err = responder.RespondSingle(w, &collection, http.StatusOK); err != nil {
		log.Errorf("error when responding with 200 error: ", err)
	}
}

func (s *Server) EditCollection(w http.ResponseWriter, r *http.Request) {
	var err error
	var collection models.Collection
	if err = json.NewDecoder(r.Body).Decode(&collection); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not read request"))
		return
	}
	vars := mux.Vars(r)
	collectionName := vars["collection_name"]
	collection.Name = collectionName

	var validationErrs []responder.Error
	if validationErrs = models.ValidateCollection(collection); len(validationErrs) > 0 {
		if err = responder.RespondErrors(w, validationErrs, http.StatusBadRequest); err != nil {
			log.Errorf("error when responding with 400 error: %v", err)
		}
		return
	}

	if err = s.database.UpdateCollection(&collection); err != nil {
		if err == pg.ErrNoRows {
			if err = responder.RespondError(w, "", "", http.StatusNotFound); err != nil {
				log.Errorf("error when responding with 400 error: %v", err)
			}
			return
		}

		if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
			log.Errorf("error when responding with 500 error %v", err)
		}
		return
	}

	if err = responder.RespondSingle(w, &collection, http.StatusCreated); err != nil {
		log.Errorf("error when responding with 201 error %v", err)
	}
}

func (s *Server) RemoveCollection(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)

	collectionID, err := strconv.Atoi(vars["collection_id"])
	if err != nil {
		responder.RespondError(w, "bad value", "collection_id", http.StatusBadRequest)
		return
	}
	var collection models.Collection
	if err = json.NewDecoder(r.Body).Decode(&collection); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not read request"))
		return
	}
	if err = s.database.DeleteCollectionByID(collectionID); err != nil {
		if err == pg.ErrNoRows {
			if err = responder.RespondError(w, "", "", http.StatusNotFound); err != nil {
				log.Errorf("error when responding with 400 error: %v", err)
			}
			return
		}

		if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
			log.Errorf("error when responding with 500 error: %v", err)
		}
		return
	}

	if err = responder.Respond(w, http.StatusOK); err != nil {
		log.Errorf("error when responding with 200 error: ", err)
	}
}

type AddBooksPayload struct {
	BooksToAdd []string `json:"books_to_add"`
}

func sliceToMap(books []models.Book) map[string]struct{} {
	m := make(map[string]struct{})
	for _, book := range books {
		m[book.ISBN] = struct{}{}
	}
	return m
}

func (s *Server) AddBooksToCollection(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	collectionID, err := strconv.Atoi(vars["collection_id"])
	if err != nil {
		responder.RespondError(w, "bad value", "collection_id", http.StatusBadRequest)
		return
	}

	var payload AddBooksPayload
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responder.RespondError(w, "could not read request", "", http.StatusBadRequest)
		return
	}

	var collection *models.Collection
	if collection, err = s.database.GetCollectionByID(collectionID); err != nil {
		if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
			log.Errorf("error when responding with 500 error: %v", err)
		}
		return
	}
	bookMap := sliceToMap(collection.Books)

	for _, newBook := range payload.BooksToAdd {
		if _, ok := bookMap[newBook]; !ok {
			book, err := s.database.GetBookByISBN(newBook)
			if err != nil {
				if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
					log.Errorf("error when responding with 500 error: %v", err)
				}
				return
			}
			if err = s.database.AddBookToCollection(book, collection); err != nil {
				log.Errorf("error when adding book to collection: %v", err)
				if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
					log.Errorf("error when responding with 500 error: %v", err)
				}
				return
			}
		}
	}
	if err = responder.Respond(w, http.StatusOK); err != nil {
		log.Errorf("error when responding with 500 error: %v", err)
	}
}

type RemoveBooksPayload struct {
	BooksToRemove []string `json:"books_to_remove"`
}

func (s *Server) RemoveBooksFromCollection(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	collectionID, err := strconv.Atoi(vars["collection_id"])
	if err != nil {
		responder.RespondError(w, "bad value", "collection_id", http.StatusBadRequest)
		return
	}

	var payload RemoveBooksPayload
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responder.RespondError(w, "could not read request", "", http.StatusBadRequest)
		return
	}

	var collection *models.Collection
	if collection, err = s.database.GetCollectionByID(collectionID); err != nil {
		if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
			log.Errorf("error when responding with 500 error: %v", err)
		}
		return
	}

	bookMap := sliceToMap(collection.Books)

	for _, bookToRemove := range payload.BooksToRemove {
		if _, ok := bookMap[bookToRemove]; ok {
			book, err := s.database.GetBookByISBN(bookToRemove)
			if err != nil {
				if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
					log.Errorf("error when responding with 500 error: %v", err)
				}
				return
			}
			if err = s.database.RemoveBookFromCollection(book, collection); err != nil {
				log.Errorf("error when adding book to collection: %v", err)
				if err = responder.RespondError(w, "something went wrong", "", http.StatusInternalServerError); err != nil {
					log.Errorf("error when responding with 500 error: %v", err)
				}
				return
			}
		}
	}
	if err = responder.Respond(w, http.StatusOK); err != nil {
		log.Errorf("error when responding with 500 error: %v", err)
	}
}
