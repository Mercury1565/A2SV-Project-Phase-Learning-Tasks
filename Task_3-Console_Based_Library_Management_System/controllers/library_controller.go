package controllers

import (
	"Task_3-Console_Based_Library_Management_System/models"
	"Task_3-Console_Based_Library_Management_System/services"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var library services.LibraryManager

func init() {
	// 'init' is automatically called when the package is imported
	library = services.NewLibrary()
}

func AnnounceInstructions() {
	fmt.Println("-------------------------------------------------------------------")
	fmt.Println("Welcome!, follow the instructions to continue\n")

	fmt.Println("Press 1: Add a new Book")
	fmt.Println("Press 2: Remove an existing Book")
	fmt.Println("Press 3: Borrow a Book")
	fmt.Println("Press 4: Return a Book")
	fmt.Println("Press 5: List all available Books")
	fmt.Println("Press 6: List all borrowed books by a member.")
	fmt.Println("Press 0: Abort")
}

func GetInstructions() int {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter an instruction: ")
	instruction, _ := reader.ReadString('\n')
	instruction = strings.TrimSpace(instruction)
	num_instruction, err := strconv.Atoi(instruction)

	for err != nil || num_instruction < 0 || num_instruction > 6 {
		fmt.Println("Invalid Instruction !")
		fmt.Print("Enter a valid instruction: ")

		instruction, _ := reader.ReadString('\n')
		instruction = strings.TrimSpace(instruction)
		num_instruction, err = strconv.Atoi(instruction)
	}

	return num_instruction
}

func BookInputHandler() models.Book {
	reader := bufio.NewReader(os.Stdin)

	// get id
	fmt.Print("Enter Book ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)
	num_id, err := strconv.Atoi(id)

	for err != nil {
		fmt.Println("Invalid Book ID !")
		fmt.Print("Enter a valid Book ID: ")

		id, _ := reader.ReadString('\n')
		id = strings.TrimSpace(id)
		num_id, err = strconv.Atoi(id)
	}

	// get title
	fmt.Print("Enter Book Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	// get author
	fmt.Print("Enter Book Author: ")
	author, _ := reader.ReadString('\n')
	author = strings.TrimSpace(author)

	s := models.Book{ID: num_id, Title: title, Author: author, Status: "Available"}
	return s
}

func MemberInputHandler() models.Member {
	reader := bufio.NewReader(os.Stdin)

	// get id
	fmt.Print("Enter Member ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)
	num_id, err := strconv.Atoi(id)

	for err != nil {
		fmt.Println("Invalid Member id !")
		fmt.Print("Enter a valid Member ID: ")

		id, _ := reader.ReadString('\n')
		id = strings.TrimSpace(id)
		num_id, err = strconv.Atoi(id)
	}

	// get name
	fmt.Print("Enter Member Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// get borrowed books
	var borrowedBooks []models.Book

	for {
		fmt.Println("Press '+' to add a book. Press 'Enter' to finsh")
		action, _ := reader.ReadString('\n')

		if action == "" {
			break
		}

		borrowedBooks = append(borrowedBooks, BookInputHandler())
	}

	return models.Member{ID: num_id, Name: name, BooksBorrowed: borrowedBooks}
}

func IDInputHandler(target string) int {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter a %q ID: ", target)
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)
	num_id, err := strconv.Atoi(id)

	for err != nil {
		fmt.Printf("Invalid %q id!\n", target)
		fmt.Printf("Enter a valid %q ID: ", target)

		id, _ := reader.ReadString('\n')
		id = strings.TrimSpace(id)
		num_id, err = strconv.Atoi(id)
	}

	return num_id
}
func BookListDisplayer(myBookList []models.Book) {
	for _, book := range myBookList {
		fmt.Printf("Book ID: %v | Book Title: %v | Book Author: %v | Book Status: %v\n", book.ID, book.Title, book.Author, book.Status)
	}
}

func MainMenu() {
	for {
		AnnounceInstructions()
		instruction := GetInstructions()

		if instruction == 0 {
			break
		}

		switch instruction {
		case 1:
			// add new book
			library.AddBook(BookInputHandler())
		case 2:
			// remove a book
			library.RemoveBook(IDInputHandler("Book"))
		case 3:
			// borrow a book
			fmt.Println(library.BorrowBook(IDInputHandler("Book"), IDInputHandler("Member")))
		case 4:
			// return a book
			fmt.Println(library.ReturnBook(IDInputHandler("Book"), IDInputHandler("Member")))
		case 5:
			// list available books
			BookListDisplayer(library.ListAvailableBooks())
		case 6:
			// list borrowed books
			BookListDisplayer(library.ListBorrowedBooks(IDInputHandler("Member")))
		}
	}
}
