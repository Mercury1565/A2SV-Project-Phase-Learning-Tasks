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

// NewLibrary creates a new instance of the Library struct.
// It initializes the Books and Members maps and returns a pointer to the Library.
func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

// AddBook adds a new book to the library.
// It takes a book object as a parameter and adds it to the library's collection of books.
func (library *Library) AddBook(book models.Book) {
	library.Books[book.ID] = book
}

// RemoveBook removes a book from the library based on the given bookID.
func (library *Library) RemoveBook(bookID int) {
	delete(library.Books, bookID)
}

// BorrowBook borrows a book from the library for a specific member.
// It takes the book ID and member ID as parameters and returns an error if any issue occurs.
// If the book is not found, it returns an error with the message "book not found".
// If the book is already borrowed, it returns an error with the message "book Unavailable".
// If the member is not found, it returns an error with the message "member not found".
// If the book is borrowed successfully, it updates the book's status to "borrowed" and adds the book to the member's list of borrowed books.
func (library *Library) BorrowBook(bookID int, memberID int) error {

	book, book_exists := library.Books[bookID]
	_, member_exists := library.Members[memberID]

	// check if book exists
	if !book_exists {
		return errors.New("book not found")
	}

	// check book status
	if book.Status == "borrowed" {
		return errors.New("book Unavailable")
	}

	// check if member exists
	if !member_exists {
		return errors.New("member not found")
	}

	// borrrow successfully
	book.Status = "borrowed"
	member := library.Members[memberID]
	member.BooksBorrowed = append(member.BooksBorrowed, book)
	return nil
}

// ReturnBook returns a book to the library by updating its status and removing it from the member's borrowed books list.
// It takes the book ID and member ID as parameters and returns an error if any of the following conditions are met:
// - The book with the given ID does not exist in the library.
// - The book with the given ID is not currently borrowed.
// - The member with the given ID does not exist in the library.
// - The member with the given ID has not borrowed the book with the given ID.
// If the book is successfully returned, its status is updated to "available" and it is removed from the member's borrowed books list.
// If any error occurs, an error message is returned.
func (library *Library) ReturnBook(bookID int, memberID int) error {

	book, book_exists := library.Books[bookID]
	_, member_exists := library.Members[memberID]

	// check if book exists
	if !book_exists {
		return errors.New("book not found")
	}

	// check book status
	if book.Status == "available" {
		return errors.New("book is Not Borrowed")
	}

	// check if member exists
	if !member_exists {
		return errors.New("member not found")
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
	}
	return errors.New("this member hasn't borrowed ")
}

// ListAvailableBooks returns a list of available books in the library.
// It iterates through the library's collection of books and checks the status of each book.
// If the status is "Available", the book is added to the list of available books.
// The function then returns the list of available books.
func (library *Library) ListAvailableBooks() []models.Book {

	var available_books []models.Book

	for _, book := range library.Books {
		if book.Status == "Available" {
			available_books = append(available_books, book)
		}
	}

	return available_books
}

// ListBorrowedBooks returns a list of books borrowed by a member.
// It takes a memberID as input and returns a slice of models.Book.
func (library *Library) ListBorrowedBooks(memberID int) []models.Book {

	return library.Members[memberID].BooksBorrowed
}
