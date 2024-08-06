package services

import (
	"Task_3-Console_Based_Library_Management_System/models"
	"errors"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct {
	Books   map[int]models.Book   // [book_id]   -> Book
	Members map[int]models.Member // [member_id] -> Member
}

func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

func (library *Library) AddBook(book models.Book) {
	// add new book
	library.Books[book.ID] = book
}

func (library *Library) RemoveBook(bookID int) {
	// remove a book
	delete(library.Books, bookID)
}

func (library *Library) BorrowBook(bookID int, memberID int) error {
	// borrow a book

	book, book_exists := library.Books[bookID]
	_, member_exists := library.Members[memberID]

	// check if book exists
	if !book_exists {
		return errors.New("Book not found")
	}

	// check book status
	if book.Status == "borrowed" {
		return errors.New("Book Unavailable!")
	}

	// check if member exists
	if !member_exists {
		return errors.New("Member not found")
	}

	// borrrow successfully
	book.Status = "borrowed"
	member := library.Members[memberID]
	member.BooksBorrowed = append(member.BooksBorrowed, book)
	return nil
}

func (library *Library) ReturnBook(bookID int, memberID int) error {
	// return a book

	book, book_exists := library.Books[bookID]
	_, member_exists := library.Members[memberID]

	// check if book exists
	if !book_exists {
		return errors.New("Book not found")
	}

	// check book status
	if book.Status == "available" {
		return errors.New("Book is Not Borrowed!")
	}

	// check if member exists
	if !member_exists {
		return errors.New("Member not found")
	}

	// return successfully
	book.Status = "available"
	for idx, curr_book := range library.Members[memberID].BooksBorrowed {
		if curr_book.ID == book.ID {
			member := library.Members[memberID]

			arr := member.BooksBorrowed
			arr = append(arr[:idx], arr[idx+1:]...)
			member.BooksBorrowed = arr
			return nil
		}
		return errors.New("This member hasn't borrowed ")
	}
	return nil
}

func (library *Library) ListAvailableBooks() []models.Book {
	// get all books with status "Available"

	var available_books []models.Book

	for _, book := range library.Books {
		if book.Status == "Available" {
			available_books = append(available_books, book)
		}
	}

	return available_books
}

func (library *Library) ListBorrowedBooks(memberID int) []models.Book {
	// get all books borrowed by member with ID 'memberID'

	return library.Members[memberID].BooksBorrowed
}
