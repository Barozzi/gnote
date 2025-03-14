package cmd

import (
	"gnote/config"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBuildDayArgs(t *testing.T) {
	testCases := []struct {
		name          string
		date          time.Time
		expectedDay   string
		showTimesheet bool
		showWednesday bool
		showExpense   bool
	}{
		{
			name:          "Friday",
			date:          time.Date(1970, time.January, 2, 12, 0, 0, 0, time.UTC),
			expectedDay:   "Friday, 2 January 1970\n",
			showTimesheet: true,
			showWednesday: false,
			showExpense:   false,
		},
		{
			name:          "Wednesday",
			date:          time.Date(1970, time.January, 7, 12, 0, 0, 0, time.UTC),
			expectedDay:   "Wednesday, 7 January 1970\n",
			showTimesheet: false,
			showWednesday: true,
			showExpense:   false,
		},
		{
			name:          "Last Wednesday of January 1970",
			date:          time.Date(1970, time.January, 28, 12, 0, 0, 0, time.UTC),
			expectedDay:   "Wednesday, 28 January 1970\n",
			showTimesheet: false,
			showWednesday: true,
			showExpense:   true,
		},
		{
			name:          "Last Friday of month",
			date:          time.Date(1971, time.December, 31, 12, 0, 0, 0, time.UTC),
			expectedDay:   "Friday, 31 December 1971\n",
			showTimesheet: true,
			showWednesday: false,
			showExpense:   true,
		},
		{
			name:          "Unix Birthday",
			date:          time.Date(1970, time.January, 1, 12, 0, 0, 0, time.UTC),
			expectedDay:   "Thursday, 1 January 1970\n",
			showTimesheet: false,
			showWednesday: false,
			showExpense:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := buildDayArgs(tc.date)

			if args.Day != tc.expectedDay {
				t.Errorf("Expected Day to be %q, but got %q", tc.expectedDay, args.Day)
			}

			if args.ShowTimesheet != tc.showTimesheet {
				t.Errorf("Expected ShowTimesheet to be %t, but got %t", tc.showTimesheet, args.ShowTimesheet)
			}

			if args.ShowWorkingWednesday != tc.showWednesday {
				t.Errorf("Expected ShowWorkingWednesday to be %t, but got %t", tc.showWednesday, args.ShowWorkingWednesday)
			}

			if args.ShowExpenseTodo != tc.showExpense {
				t.Errorf("Expected ShowExpenseTodo to be %t, but got %t", tc.showExpense, args.ShowExpenseTodo)
			}
		})
	}
}

// MockConfigReader allows us to control the config values during testing
type MockConfigReader struct {
	Config config.Config
	Error  error
}

func (m MockConfigReader) ReadConfig() (*config.Config, error) {
	return &m.Config, m.Error
}

// MockEditor allows us to avoid actually opening a file in an editor during testing
type MockEditor struct {
	OpenFileFunc func(string) error
	FilePath     string
}

func TestCreateDayFileQuarterly(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Mock the config reader to return a known configuration
	mockConfig := MockConfigReader{
		Config: config.Config{
			VaultPath: tempDir,
			DayPath:   "days",
		},
		Error: nil,
	}

	config.ReadConfigMock = mockConfig.ReadConfig
	defer func() { config.ReadConfigMock = nil }() // Restore original function

	testCases := []struct {
		name           string
		time           time.Time
		expectedFolder string
		expectedFile   string
	}{
		{
			name:           "Q1 1970",
			time:           time.Date(1970, time.January, 25, 12, 0, 0, 0, time.UTC),
			expectedFolder: filepath.Join(tempDir, "days", "1970_Q1"),
			expectedFile:   filepath.Join(tempDir, "days", "1970_Q1", "1-25-1970.md"),
		},
		{
			name:           "Q2 1970",
			time:           time.Date(1970, time.April, 15, 12, 0, 0, 0, time.UTC),
			expectedFolder: filepath.Join(tempDir, "days", "1970_Q2"),
			expectedFile:   filepath.Join(tempDir, "days", "1970_Q2", "4-15-1970.md"),
		},
		{
			name:           "Q3 1970",
			time:           time.Date(1970, time.August, 10, 12, 0, 0, 0, time.UTC),
			expectedFolder: filepath.Join(tempDir, "days", "1970_Q3"),
			expectedFile:   filepath.Join(tempDir, "days", "1970_Q3", "8-10-1970.md"),
		},
		{
			name:           "Q4 1970",
			time:           time.Date(1970, time.December, 1, 12, 0, 0, 0, time.UTC),
			expectedFolder: filepath.Join(tempDir, "days", "1970_Q4"),
			expectedFile:   filepath.Join(tempDir, "days", "1970_Q4", "12-1-1970.md"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testArgs := DayArgs{
				Day:                  "Test Day\n",
				ShowTimesheet:        false,
				ShowWorkingWednesday: false,
				ShowExpenseTodo:      false,
			}

			filePath, err := createDayFile(testArgs, tc.time)
			if err != nil {
				t.Fatalf("createDayFile returned an error: %v", err)
			}

			if filePath != tc.expectedFile {
				t.Errorf("Expected file path to be %q, but got %q", tc.expectedFile, filePath)
			}

			// Check if the folder was actually created
			_, err = os.Stat(tc.expectedFolder)
			if os.IsNotExist(err) {
				t.Errorf("Expected folder to be created at %q, but it doesn't exist", tc.expectedFolder)
			}

			// Check if the file was actually created
			_, err = os.Stat(filePath)
			if os.IsNotExist(err) {
				t.Errorf("Expected file to be created at %q, but it doesn't exist", filePath)
			}

			// Clean up temp dir - remove to test manually
			os.RemoveAll(tempDir)
		})
	}
}

func TestCreateDayFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Mock the config reader to return a known configuration
	mockConfig := MockConfigReader{
		Config: config.Config{
			VaultPath: tempDir,
			DayPath:   "days",
		},
		Error: nil,
	}

	config.ReadConfigMock = mockConfig.ReadConfig
	defer func() { config.ReadConfigMock = nil }() // Restore original function

	// Test Time
	testTime := time.Date(2024, time.January, 25, 12, 0, 0, 0, time.UTC)
	testArgs := DayArgs{
		Day:                  "Thursday, 25 January 2024\n",
		ShowTimesheet:        false,
		ShowWorkingWednesday: false,
		ShowExpenseTodo:      false,
	}

	filePath, err := createDayFile(testArgs, testTime)
	if err != nil {
		t.Fatalf("createDayFile returned an error: %v", err)
	}

	expectedFolderPath := filepath.Join(tempDir, "days", "2024_Q1")
	expectedFilePath := filepath.Join(expectedFolderPath, "1-25-2024.md")

	if filePath != expectedFilePath {
		t.Errorf("Expected file path to be %q, but got %q", expectedFilePath, filePath)
	}

	// Check if the folder was actually created
	_, err = os.Stat(expectedFolderPath)
	if os.IsNotExist(err) {
		t.Errorf("Expected folder to be created at %q, but it doesn't exist", expectedFolderPath)
	}

	// Check if the file was actually created
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("Expected file to be created at %q, but it doesn't exist", filePath)
	}

	// Clean up temp dir
	os.RemoveAll(tempDir)
}
