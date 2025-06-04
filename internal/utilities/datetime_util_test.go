package utilities

import (
	"github.com/go-playground/assert/v2"
	"testing"
	"time"
)

func TestAddDay(t *testing.T) {
	layout := "2006-01-02 15:04"
	currentDateTime := "2023-06-10 13:45"
	addedDays := "2023-06-12 13:45"

	// Parse the MySQL datetime string into a time.Time object
	parsedTime, err := time.Parse(layout, currentDateTime)
	if err != nil {
		t.Error(err)
		return
	}

	addedDateTime := AddDay(parsedTime, 2)
	parsedAddedDateTime, err := time.Parse(layout, addedDateTime.Format(layout))
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, parsedAddedDateTime.Format(layout), addedDays)
}
