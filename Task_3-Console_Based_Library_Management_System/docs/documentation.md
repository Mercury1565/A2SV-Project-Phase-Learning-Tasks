# Library Management System Documentation

## Project Overview
The Library Management System is a simple console-based application implemented in Go which allows users to manage a collection of books and track borrowed books by library members.

## Folder Structure

library_management/
├── main.go
├── controllers/
│ └── library_controller.go
├── models/
│ └── book.go
│ └── member.go
├── services/
│ └── library_service.go
├── docs/
│ └── documentation.md
└── go.mod

## Structs

### Book
Defines a `Book` struct with the following fields:
- `ID` (int): Unique identifier for the book.
- `Title` (string): Title of the book.
- `Author` (string): Author of the book.
- `Status` (string): Status of the book, either "Available" or "Borrowed".

### Member
Defines a `Member` struct with the following fields:
- `ID` (int): Unique identifier for the member.
- `Name` (string): Name of the member.
- `BorrowedBooks` ([]Book): Slice to hold borrowed books.

## Interfaces

### LibraryManager
Defines a `LibraryManager` interface with the following methods:
- `AddBook(book Book)`
- `RemoveBook(bookID int)`
- `BorrowBook(bookID int, memberID int) error`
- `ReturnBook(bookID int, memberID int) error`
- `ListAvailableBooks() []Book`
- `ListBorrowedBooks(memberID int) []Book`

## Implementation

### Library
Implements the `LibraryManager` interface in a `Library` struct. The `Library` struct has:
- `Books` (map[int]Book): Stores all books with the book ID as the key.
- `Members` (map[int]Member): Stores all members with the member ID as the key.

## Console Interaction
Numbers from 0 upto 6 are assigned for different commands in the console:
- Press 0: Abort
- Press 1: Add a new Book
- Press 2: Remove an existing Book
- Press 3: Borrow a Book
- Press 4: Return a Book
- Press 5: List all available Books
- Press 6: List all borrowed books by a member

## Usage
To run the project, simply execute the main Go file:
```go run main.go```

## Some Example Usages

### Add Book

Welcome!, follow the instructions to continue

Press 1: Add a new Book
Press 2: Remove an existing Book
Press 3: Borrow a Book
Press 4: Return a Book
Press 5: List all available Books
Press 6: List all borrowed books by a member.
Press 0: Abort
Enter an instruction: 1
Enter Book ID: 1
Enter Book Title: The Enemy
Enter Book Author: David Aspinov
Book Added Successfully

## List Available Books

Welcome!, follow the instructions to continue

Press 1: Add a new Book
Press 2: Remove an existing Book
Press 3: Borrow a Book
Press 4: Return a Book
Press 5: List all available Books
Press 6: List all borrowed books by a member.
Press 0: Abort
Enter an instruction: 5
Book ID: 1 | Book Title: The Enemy | Book Author: David Aspinov | Book Status: Available
Book ID: 2 | Book Title: Feker Eske Mekaber | Book Author: Hadis Alemayehu | Book Status: Available

