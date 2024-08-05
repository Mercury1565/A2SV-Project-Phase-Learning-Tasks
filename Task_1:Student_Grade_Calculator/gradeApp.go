package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var name string
var grades = make(map[string]float64)

func displayGrade(name string) {
	// display name
	fmt.Printf("\nThe grades for the student %q are:\n", name)

	// display each grade
	var sum float64 = 0
	for subject, grade := range grades {
		fmt.Printf("Subject: %s-----Grade: %g\n", subject, grade)
		sum += grade
	}

	// display average
	var avg float64 = sum / float64(len(grades))
	fmt.Printf("\nThe average grade: %v\n", avg)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// get student name
	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Printf("Hello, %v \n", name)

	for {
		// get subject
		fmt.Print("\nEnter subject, Press 'Enter' to finish: ")
		subject, _ := reader.ReadString('\n')
		subject = strings.TrimSpace(subject)

		// abort for empty input
		if subject == "" {
			displayGrade(name)
			break
		}

		// get subject
		fmt.Printf("Enter grade for the subject %q: ", subject)

		grade, _ := reader.ReadString('\n')
		grade = strings.TrimSpace(grade)

		num_grade, err := strconv.ParseFloat(grade, 64)

		for {
			if err != nil || num_grade < 0 || num_grade > 100 {
				// invalid grade input

				fmt.Printf("\n%v is an invalid grade input !!!\n", grade)
				fmt.Printf("Please enter valid grade for the subject %q: ", subject)

				grade, _ = reader.ReadString('\n')
				grade = strings.TrimSpace(grade)
				num_grade, err = strconv.ParseFloat(grade, 64)
			} else {
				grades[subject] = num_grade
				break
			}
		}

		grades[subject] = num_grade
	}
}
