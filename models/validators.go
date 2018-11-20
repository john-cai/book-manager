package models

import "github.com/john-cai/book-manager/responder"

func ValidateBook(book Book) []responder.Error {
	var errs []responder.Error
	if book.ISBN == "" {
		errs = append(errs, responder.Error{
			Field:   "isbn",
			Message: "required",
		})
	}
	if book.Title == "" {
		errs = append(errs, responder.Error{
			Field:   "title",
			Message: "required",
		})
	}
	if book.Author == "" {
		errs = append(errs, responder.Error{
			Field:   "author",
			Message: "required",
		})
	}
	//TODO: rest of logic
	return errs
}

func ValidateCollection(collection Collection) []responder.Error {
	//TODO: implement me
	return []responder.Error{}
}
