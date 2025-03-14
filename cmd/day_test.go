package cmd

import (
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
