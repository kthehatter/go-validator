package validator

import (
	"testing"
)

func TestRemoveSpecialChars(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"Hello!@# World!", "Hello World"},
		{"123#456", "123456"},
		{123, 123}, // Non-string input
	}

	for _, test := range tests {
		result := RemoveSpecialChars(test.input)
		if result != test.expected {
			t.Errorf("RemoveSpecialChars(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"hello world", "Hello World"},
		{"HELLO WORLD", "Hello World"},
		{"çimant", "Çimant"},
		{"hello, world!", "Hello, World!"}, // Handles punctuation correctly
		{"123", "123"},                     // Non-alphabetic input
		{123, 123},                         // Non-string input
	}

	for _, test := range tests {
		result := ToTitleCase(test.input)
		if result != test.expected {
			t.Errorf("ToTitleCase(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"123", 123},
		{123.45, 123},
		{"abc", "abc"}, // Invalid input
	}

	for _, test := range tests {
		result := ToInt(test.input)
		if result != test.expected {
			t.Errorf("ToInt(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"123.45", 123.45},
		{123, 123.0},
		{"abc", "abc"}, // Invalid input
	}

	for _, test := range tests {
		result := ToFloat(test.input)
		if result != test.expected {
			t.Errorf("ToFloat(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	transformer := Truncate(5)

	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"hello world", "hello"},
		{"hi", "hi"},
		{123, 123}, // Non-string input
	}

	for _, test := range tests {
		result := transformer(test.input)
		if result != test.expected {
			t.Errorf("Truncate(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestReplace(t *testing.T) {
	transformer := Replace("foo", "bar")

	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"foo bar", "bar bar"},
		{"hello world", "hello world"},
		{123, 123}, // Non-string input
	}

	for _, test := range tests {
		result := transformer(test.input)
		if result != test.expected {
			t.Errorf("Replace(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}
