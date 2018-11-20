package models

import (
	"time"

	"github.com/go-pg/pg/orm"
)

func init() {
	orm.RegisterTable((*Collection)(nil))
	orm.RegisterTable((*Book)(nil))
	orm.RegisterTable((*BookCollection)(nil))
}

type Book struct {
	ISBN        string       `sql:"isbn,pk",json:"isbn"`
	Title       string       `json:"title"`
	Author      string       `json:"author"`
	Description string       `json:"description"`
	PublishedAt time.Time    `json:"published_at"`
	Metadata    Metadata     `json:"metadata"`
	Collections []Collection `pg:"many2many:book_collections,joinFK:collection_id", json:"collections"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   time.Time    `pg:",soft_delete" json:"deleted_at"`
}

func (b *Book) BeforeInsert(db orm.DB) error {
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	return nil
}

type Collection struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Books       []Book    `pg:"many2many:book_collections,joinFK:book_isbn", json:"books"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `pg:",soft_delete" json:"deleted_at"`
}

func (c *Collection) BeforeInsert(db orm.DB) error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	return nil
}

type BookCollection struct {
	BookISBN     string    `sql:"book_isbn,pk"`
	CollectionID int       `sql:"collection_id,pk"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `pg:",soft_delete" json:"deleted_at"`
}

type Metadata struct {
	Genres []string `json:"genres"`
}
