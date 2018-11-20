package database

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/john-cai/book-manager/models"
)

// Database is a top level database object we can tie database methods to
type Database struct {
	db *pg.DB
}

// New creates a new database object
func New(user, addr, database string) (*Database, error) {

	db := pg.Connect(&pg.Options{
		User:     user,
		Addr:     addr,
		Database: database,
	})
	return &Database{
		db: db,
	}, nil
}

func NewTestDB() (*Database, error) {
	return New("postgres", "localhost:5432", "bookmanager_test")
}

func ignoreSoftDeletedBookCollections(q *orm.Query) (*orm.Query, error) {
	return q.Where("book_collections.deleted_at is null"), nil
}

func (d *Database) GetBookByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	if err := d.db.Model(&book).Where("isbn = ?", isbn).Relation("Collections").First(); err != nil {
		return nil, err
	}
	return &book, nil
}

func sliceToPGArray(s []string) string {
	var b bytes.Buffer
	b.WriteString("array[")
	for i, v := range s {
		s[i] = fmt.Sprint("'%s'", v)
	}
	b.WriteString(strings.Join(s, ","))
	b.WriteString("]")
	return b.String()
}

func (d *Database) GetBooks(isbn, title, author string, publishedYear int, genres []string) ([]models.Book, error) {
	var books []models.Book
	q := d.db.Model(&books)
	if isbn != "" {
		q = q.Where("isbn = ?", isbn)
	}
	if title != "" {
		q = q.Where("title = ?", title)
	}
	if author != "" {
		q = q.Where("author = ?", author)
	}
	if publishedYear != 0 {
		q = q.Where("published_at = ?", time.Date(publishedYear, 0, 0, 0, 0, 0, 0, time.UTC))
	}

	if genres != nil && len(genres) > 0 {
		q = q.Where(fmt.Sprintf("metadata->'genres' ?| %s", sliceToPGArray(genres)))
	}
	if err := q.Relation("Collections").Select(); err != nil {
		return nil, err
	}
	return books, nil
}

func (d *Database) AddBook(b *models.Book) error {
	return d.db.Insert(b)
}

func (d *Database) UpdateBook(b *models.Book) error {
	return d.db.Update(b)
}

func (d *Database) DeleteBookByISBN(isbn string) error {
	return d.db.Delete(&models.Book{
		ISBN: isbn,
	})
}

func (d *Database) GetCollectionByID(id int) (*models.Collection, error) {
	var collection models.Collection
	if err := d.db.Model(&collection).Where("id = ?", id).Relation("Books", ignoreSoftDeletedBookCollections).First(); err != nil {
		return nil, err
	}
	return &collection, nil
}
func (d *Database) GetAllCollections() ([]models.Collection, error) {
	var collections []models.Collection
	if err := d.db.Model(&collections).Relation("Books", ignoreSoftDeletedBookCollections).Select(); err != nil {
		return nil, err
	}
	return collections, nil
}
func (d *Database) AddCollection(c *models.Collection) error {
	return d.db.Insert(c)
}

func (d *Database) UpdateCollection(c *models.Collection) error {
	return d.db.Update(c)
}

func (d *Database) DeleteCollectionByID(id int) error {
	return d.db.Delete(&models.Collection{ID: id})
}

func (d *Database) AddBookToCollection(b *models.Book, c *models.Collection) error {
	if b.ISBN == "" {
		return errors.New("book isbn missing")
	}
	if c.ID == 0 {
		return errors.New("collection id missing")
	}
	return d.db.Insert(&models.BookCollection{
		BookISBN:     b.ISBN,
		CollectionID: c.ID,
	})
}

func (d *Database) RemoveBookFromCollection(b *models.Book, c *models.Collection) error {
	if b.ISBN == "" {
		return errors.New("book isbn missing")
	}
	if c.ID == 0 {
		return errors.New("collection id missing")
	}
	return d.db.Delete(&models.BookCollection{
		BookISBN:     b.ISBN,
		CollectionID: c.ID,
	})
}
