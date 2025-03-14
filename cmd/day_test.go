package cmd

import (
	"testing"
	"time"

	"gnote/config"
)

// MockConfigReader allows us to control the config values during testing
type MockConfigReader struct {
	Config config.Config
	Error  error
}

func (m MockConfigReader) ReadConfig() (config.Config, error) {
	return m.Config, m.Error
}

// MockEditor allows us to avoid actually opening a file in an editor during testing
type MockEditor struct {
	OpenFileFunc func(string) error
	FilePath     string
}

func (m *MockEditor) OpenFile(filePath string) error {
	m.FilePath = filePath
	if m.OpenFileFunc != nil {
		return m.OpenFileFunc(filePath)
	}
	return nil
}

func TestBuildDayArgs(t *testing.T) {
	friday := time.Date(1970, time.January, 2, 12, 0, 0, 0, time.UTC)
	args := buildDayArgs(friday)
	if friday.Weekday().String() != "Friday" {
		t.Error("Expected day to be Friday but got ", friday.Weekday().String())
	}
	if !args.ShowTimesheet {
		t.Error("Expected ShowTimesheet to be true on Friday")
	}

	// Test case 2: Wednesday, should show WorkingWednesday
	wednesday := time.Date(1970, time.January, 7, 12, 0, 0, 0, time.UTC)
	args = buildDayArgs(wednesday)
	if wednesday.Weekday().String() != "Wednesday" {
		t.Error("Expected day to be Wednesday but got ", wednesday.Weekday().String())
	}
	if !args.ShowWorkingWednesday {
		t.Error("Expected ShowWorkingWednesday to be true on Wednesday")
	}

	// Test case 3: Last Wednesday of the month, should show ExpenseTodo
	lastWednesday := time.Date(1970, time.January, 28, 12, 0, 0, 0, time.UTC) // Last Wednesday in Jan 2024
	args = buildDayArgs(lastWednesday)
	if lastWednesday.Weekday().String() != "Wednesday" {
		t.Error("Expected day to be Wednesday but got ", lastWednesday.Weekday().String())
	}
	if !args.ShowExpenseTodo {
		t.Error("Expected ShowExpenseTodo to be true on the last Wednesday of the month")
	}

	// Test case 4: Check date formatting
	unixBirthday := time.Date(1970, time.January, 1, 12, 0, 0, 0, time.UTC)
	expectedDay := "Thursday, 1 January 1970\n"
	args = buildDayArgs(unixBirthday)
	if args.Day != expectedDay {
		t.Errorf("Expected day to be %q, but got %q", expectedDay, args.Day)
	}
}
