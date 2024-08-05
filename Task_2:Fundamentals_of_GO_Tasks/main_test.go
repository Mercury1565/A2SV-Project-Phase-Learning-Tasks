package main

import (
	"reflect"
	"testing"
)

func TestFreqCount(t *testing.T) {
	input := "Hello world hello"
	expected := map[string]int{
		"hello": 2,
		"world": 1,
	}

	result := freqCount(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestIsPalindrome(t *testing.T) {
	input := "abcddcba"
	expected := true
	result := isPalindrome(input)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	input = "hello"
	expected = false
	result = isPalindrome(input)

	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
