package main

import (
	"strings"
)

func freqCount(s string) map[string]int {
	splitted := strings.Split(s, " ")
	freq := make(map[string]int)

	for _, word := range splitted {
		if word == " "{
			continue
		}
		word = strings.ToLower(word)
		freq[word] += 1
	}
	return freq
}

func isPalindrome(s string) bool {
	left := 0
	right := len(s) - 1

	for left < right {
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}

	return true
}

func main() {

}
