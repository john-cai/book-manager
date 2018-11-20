package database

import (
	"testing"

	"github.com/john-cai/book-manager/models"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBookByISBN(t *testing.T) {
	db, err := NewTestDB()
	require.NoError(t, err)
	// insert a book and collection

	book := models.Book{
		ISBN:   uuid.New(),
		Title:  "1",
		Author: "abc",
	}
	collection := models.Collection{Name: "collection1"}
	values := []interface{}{
		&book,
		&collection,
		&models.BookCollection{BookISBN: book.ISBN, CollectionID: collection.ID},
	}
	for _, v := range values {
		err := db.db.Insert(v)
		if err != nil {
			panic(err)
		}
	}
	require.NoError(t, db.db.Insert(&models.BookCollection{BookISBN: book.ISBN, CollectionID: collection.ID}))
	b, err := db.GetBookByISBN(book.ISBN)
	require.NoError(t, err)
	assert.Len(t, b.Collections, 1)
	assert.Equal(t, collection.Name, b.Collections[0].Name)
}

func TestGetCollectionByID(t *testing.T) {
	db, err := NewTestDB()
	require.NoError(t, err)
	// insert a book and collection
	book := models.Book{
		ISBN:   uuid.New(),
		Title:  "1",
		Author: "abc",
	}
	collection := models.Collection{Name: "collection1"}
	values := []interface{}{
		&book,
		&collection,
	}
	for _, v := range values {
		err := db.db.Insert(v)
		if err != nil {
			panic(err)
		}
	}
	require.NoError(t, db.db.Insert(&models.BookCollection{BookISBN: book.ISBN, CollectionID: collection.ID}))
	c, err := db.GetCollectionByID(collection.ID)
	require.NoError(t, err)
	assert.Len(t, c.Books, 1)
	assert.Equal(t, book.Title, c.Books[0].Title)

}
