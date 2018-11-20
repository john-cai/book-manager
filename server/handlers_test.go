package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/john-cai/book-manager/database"
	"github.com/john-cai/book-manager/models"
	"github.com/john-cai/book-manager/responder"
)

func setUpTestServer(t *testing.T) *Server {
	db, err := database.NewTestDB()
	require.NoError(t, err)
	s := &Server{
		database: db,
		Router:   mux.NewRouter(),
	}
	s.configureRoutes()
	return s
}

func TestAddBook(t *testing.T) {
	s := setUpTestServer(t)
	testCases := []struct {
		input        models.Book
		responseCode int
	}{
		{input: models.Book{ISBN: uuid.New(), Title: "Jungle Book", Author: "Rudyard Kipling", PublishedAt: time.Date(1894, 0, 0, 0, 0, 0, 0, time.UTC)}, responseCode: http.StatusCreated},
		{input: models.Book{ISBN: "", Title: "2001: A Space Odyssey", Author: "abc", PublishedAt: time.Date(1994, 0, 0, 0, 0, 0, 0, time.UTC)}, responseCode: http.StatusBadRequest},
		{input: models.Book{ISBN: uuid.New(), Title: "", Author: "Rudyard Kipling", PublishedAt: time.Date(1894, 0, 0, 0, 0, 0, 0, time.UTC)}, responseCode: http.StatusBadRequest},
	}

	for _, testCase := range testCases {
		rec := httptest.NewRecorder()
		var b bytes.Buffer
		json.NewEncoder(&b).Encode(&testCase.input)
		req := httptest.NewRequest(http.MethodPost, "/books", &b)
		s.ServeHTTP(rec, req)
		assert.Equal(t, testCase.responseCode, rec.Result().StatusCode)
	}
}

func TestEditBook(t *testing.T) {
	s := setUpTestServer(t)
	testCases := []struct {
		original     models.Book
		edited       models.Book
		responseCode int
	}{
		{
			original:     models.Book{ISBN: uuid.New(), Title: "Jungle Book", Author: "Rudyard Kipling", PublishedAt: time.Date(1894, 0, 0, 0, 0, 0, 0, time.UTC), CreatedAt: time.Now()},
			edited:       models.Book{Title: "Jungle Book II", Author: "Rudyard Kipling", PublishedAt: time.Date(1899, 0, 0, 0, 0, 0, 0, time.UTC)},
			responseCode: http.StatusOK,
		},
		{
			original:     models.Book{ISBN: uuid.New(), Title: "2001: A Space Odyssey", Author: "abc", PublishedAt: time.Date(1994, 0, 0, 0, 0, 0, 0, time.UTC), CreatedAt: time.Now()},
			edited:       models.Book{Title: "2018: A Space Odyssey", Author: "abc", PublishedAt: time.Date(1994, 0, 0, 0, 0, 0, 0, time.UTC)},
			responseCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		// add book to the system
		rec := httptest.NewRecorder()
		var b bytes.Buffer
		json.NewEncoder(&b).Encode(&testCase.original)
		req := httptest.NewRequest(http.MethodPost, "/books", &b)
		s.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Result().StatusCode)
		b.Reset()
		testCase.edited.CreatedAt = testCase.original.CreatedAt
		json.NewEncoder(&b).Encode(&testCase.edited)

		// modify book
		rec = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/books/%s", testCase.original.ISBN), &b)
		s.ServeHTTP(rec, req)

		if rec.Result().StatusCode != testCase.responseCode {
			var errorResponse responder.ErrorResponse
			require.NoError(t, json.NewDecoder(rec.Result().Body).Decode(&errorResponse))
		}

		require.Equal(t, testCase.responseCode, rec.Result().StatusCode)

		// ensure book has changed on the backend
		rec = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books/%s", testCase.original.ISBN), nil)
		s.ServeHTTP(rec, req)
		var editedBook models.Book

		require.NoError(t, json.NewDecoder(rec.Result().Body).Decode(&editedBook))
		assert.Equal(t, testCase.original.ISBN, editedBook.ISBN)
		assert.Equal(t, testCase.edited.Title, editedBook.Title)
		assert.Equal(t, testCase.edited.Author, editedBook.Author)
		assert.Equal(t, testCase.edited.Description, editedBook.Description)
		assert.Equal(t, testCase.edited.PublishedAt.Year(), editedBook.PublishedAt.Year())
	}
}

func TestDeleteBook(t *testing.T) {
	s := setUpTestServer(t)
	books := []models.Book{
		models.Book{ISBN: uuid.New(), Title: "Lord of the Rings", Author: "J.R.R. Tolkien", PublishedAt: time.Date(1944, 0, 0, 0, 0, 0, 0, time.UTC), CreatedAt: time.Now()},
		models.Book{ISBN: uuid.New(), Title: "A Wrinkle in Time", Author: "Madeline L'engle", PublishedAt: time.Date(1989, 0, 0, 0, 0, 0, 0, time.UTC), CreatedAt: time.Now()},
	}

	for _, book := range books {
		// add book to the system
		rec := httptest.NewRecorder()
		var b bytes.Buffer
		json.NewEncoder(&b).Encode(&book)
		req := httptest.NewRequest(http.MethodPost, "/books", &b)
		s.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Result().StatusCode)
	}

	testCases := []struct {
		isbn     string
		response int
	}{
		{
			isbn:     books[0].ISBN,
			response: http.StatusOK,
		},
		{
			isbn:     books[1].ISBN,
			response: http.StatusOK,
		},
		{
			isbn:     uuid.New(),
			response: http.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/books/%s", testCase.isbn), nil)
		s.ServeHTTP(rec, req)
		assert.Equal(t, testCase.response, rec.Result().StatusCode)
	}
}

func TestAddCollection(t *testing.T) {
	s := setUpTestServer(t)
	testCases := []struct {
		input        models.Collection
		responseCode int
	}{
		{input: models.Collection{Name: "collection1", Description: "a great collection of books"}, responseCode: http.StatusCreated},
		{input: models.Collection{Name: "collection2", Description: "another great collection of books"}, responseCode: http.StatusCreated},
		{input: models.Collection{Name: "collection3", Description: "this one's just alright"}, responseCode: http.StatusCreated},
	}

	for _, testCase := range testCases {
		rec := httptest.NewRecorder()
		var b bytes.Buffer
		json.NewEncoder(&b).Encode(&testCase.input)
		req := httptest.NewRequest(http.MethodPost, "/collections", &b)
		s.ServeHTTP(rec, req)
		assert.Equal(t, testCase.responseCode, rec.Result().StatusCode)
	}
}

func TestAddBooksToCollection(t *testing.T) {
	s := setUpTestServer(t)
	// add a collection
	collection := models.Collection{Name: "collection1", Description: "a great collection of books"}
	rec := httptest.NewRecorder()
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(&collection)
	req := httptest.NewRequest(http.MethodPost, "/collections", &b)
	s.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Result().StatusCode)
	require.NoError(t, json.NewDecoder(rec.Result().Body).Decode(&collection))

	// add some books
	books := []models.Book{
		models.Book{ISBN: uuid.New(), Title: "Lord of the Rings: Return of the King", Author: "J.R.R. Tolkien", PublishedAt: time.Date(1944, 0, 0, 0, 0, 0, 0, time.UTC), CreatedAt: time.Now()},
		models.Book{ISBN: uuid.New(), Title: "A Wrinkle in Time II", Author: "Madeline L'engle", PublishedAt: time.Date(1989, 0, 0, 0, 0, 0, 0, time.UTC), CreatedAt: time.Now()},
	}

	var booksToAdd []string
	for _, book := range books {
		// add book to the system
		rec := httptest.NewRecorder()
		var b bytes.Buffer
		json.NewEncoder(&b).Encode(&book)
		req := httptest.NewRequest(http.MethodPost, "/books", &b)
		s.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Result().StatusCode)

		booksToAdd = append(booksToAdd, book.ISBN)
	}

	addBooksPayload := AddBooksPayload{
		BooksToAdd: booksToAdd,
	}
	rec = httptest.NewRecorder()
	b.Reset()
	json.NewEncoder(&b).Encode(&addBooksPayload)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/collections/%d/addbooks", collection.ID), &b)
	s.ServeHTTP(rec, req)
	var errResponse responder.ErrorResponse
	json.NewDecoder(rec.Result().Body).Decode(&errResponse)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/collections/%d", collection.ID), nil)
	s.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.NoError(t, json.NewDecoder(rec.Result().Body).Decode(&collection))
	assert.Len(t, collection.Books, len(books))

	// just remove one of the books
	rec = httptest.NewRecorder()
	b.Reset()
	removeBooksPayload := RemoveBooksPayload{
		BooksToRemove: []string{booksToAdd[0]},
	}
	json.NewEncoder(&b).Encode(&removeBooksPayload)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/collections/%d/removebooks", collection.ID), &b)
	s.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/collections/%d", collection.ID), nil)
	s.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.NoError(t, json.NewDecoder(rec.Result().Body).Decode(&collection))
	assert.Len(t, collection.Books, len(books)-1)

}
