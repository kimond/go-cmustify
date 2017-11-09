package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("Should parse status", func(t *testing.T) {
		result := Parse("status playing")
		assert.Equal(t, result["status"], "playing")
	})

	t.Run("Should parse multiple words status", func(t *testing.T) {
		result := Parse("status playing loud")
		assert.Equal(t, result["status"], "playing loud")
	})

	t.Run("Should parse multiple tokens", func(t *testing.T) {
		result := Parse("status playing loud title Yeah")
		assert.Equal(t, result["status"], "playing loud")
		assert.Equal(t, result["title"], "Yeah")
	})
}

func TestFormatMessageBody(t *testing.T) {
	table := []struct {
		name     string
		metaData string
		expected string
	}{
		{"Should format data with title",
			"status playing title Dancing with the stars", "Dancing with the stars"},
		{"Should format data with title and artist",
			"status playing title Dancing with the stars artist Jamo", "Dancing with the stars by Jamo"},
		{"Should format data with title, artist and album",
			"status playing title Dancing with the stars artist Jamo album Best", "Dancing with the stars by Jamo, Best"},
		{"Should format data with no title",
			"status playing", "Unknown"},
	}
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			result := Parse(tt.metaData)
			messageBody := FormatMessageBody(result)
			assert.Equal(t, messageBody, tt.expected)
		})
	}
}

func TestLibNotify(t *testing.T) {
	n := new(libNotify)
	n.send("test", "content")
}

type fakeNotifier struct {
	Called bool
}

func (n *fakeNotifier) send(s, b string) error { n.Called = true; return nil }

func TestHandleData(t *testing.T) {
	t.Run("Should call the notifier", func(t *testing.T) {
		n := new(fakeNotifier)
		HandleData(n, "status playing")
		assert.True(t, n.Called)
	})
}

func TestPrintUsage(t *testing.T) {
	printUsage()
}
