package models

type Member struct {
	ID            int
	Name          string
	BooksBorrowed []Book
}
