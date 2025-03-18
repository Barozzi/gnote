package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"gnote/config"
	"os/exec"

	"github.com/spf13/cobra"
)

type DayArgs struct {
	Day                  string
	ShowTimesheet        bool
	ShowWorkingWednesday bool
	ShowExpenseTodo      bool
}

type Editor interface {
	OpenFile(string) error
}

type NvimEditor struct{}

func (n NvimEditor) OpenFile(filePath string) error {
	vim := exec.Command("zsh", "-c", fmt.Sprintf("nvim %s", filePath))
	vim.Stdin = os.Stdin
	vim.Stdout = os.Stdout
	vim.Stderr = os.Stderr
	err := vim.Start()
	if err != nil {
		return fmt.Errorf("Vim failed to start correctly: %w", err)
	}
	err = vim.Wait()
	if err != nil {
		return fmt.Errorf("Vim failed to exit correctly: %w", err)
	}
	return nil
}

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "Create a new DevLog for the current day.",
	Run: func(cmd *cobra.Command, args []string) {
		timeNow := time.Now()
		dayArgs := buildDayArgs(timeNow)

		filePath, err := createDayFile(dayArgs, timeNow)
		if err != nil {
			fmt.Printf("Failed to create day file: %s\n", err)
			os.Exit(1)
		}

		editor := NvimEditor{}
		if err := editor.OpenFile(filePath); err != nil {
			fmt.Printf("Failed to open file in editor: %s\n", err)
			os.Exit(1)
		}
	},
}

func buildDayArgs(timeNow time.Time) DayArgs {
	formattedDay := fmt.Sprintf("%s, %d %s %d\n", timeNow.Weekday(), timeNow.Day(), timeNow.Month().String(), timeNow.Year())
	showTimesheet := timeNow.Weekday() == time.Friday
	showWorkingWednesday := timeNow.Weekday() == time.Wednesday
	showExpenseTodo := isLastWeekdayOfMonth(timeNow)
	return DayArgs{
		Day:                  formattedDay,
		ShowTimesheet:        showTimesheet,
		ShowWorkingWednesday: showWorkingWednesday,
		ShowExpenseTodo:      showExpenseTodo,
	}
}

func createDayFile(args DayArgs, timeNow time.Time) (string, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		return "", err
	}

	year := timeNow.Year()
	quarter := getQuarter(timeNow)
	quarterFolder := fmt.Sprintf("%d_Q%d", year, quarter)

	folderPath := filepath.Join(cfg.VaultPath, cfg.DayPath, quarterFolder)

	// Create the folder if it doesn't exist
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", err
	}

	filePath := filepath.Join(folderPath, fmt.Sprintf("%d-%d-%d.md", timeNow.Month(), timeNow.Day(), timeNow.Year()))

	_, err = os.Stat(filePath)
	if err == nil {
		// File exists, don't overwrite
		return filePath, nil
	}

	// I'm expecting a file-does-not-exist error; If its not that kind of error lets raise it up the stack
	if !os.IsNotExist(err) {
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	newDayT := newDayTemplate()
	err = newDayT.Execute(file, args)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func getQuarter(date time.Time) int {
	month := date.Month()
	switch {
	case month >= time.January && month <= time.March:
		return 1
	case month >= time.April && month <= time.June:
		return 2
	case month >= time.July && month <= time.September:
		return 3
	case month >= time.October && month <= time.December:
		return 4
	default:
		return 0 // Should never happen
	}
}

func isLastWeekdayOfMonth(date time.Time) bool {
	// Get the last day of the month
	// Using the zeroth day trick to get the last day of month
	lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())

	// Check if the given date is one of the last weekdays of the month
	for i := 0; i < 7; i++ {
		currentDay := lastDay.AddDate(0, 0, -i)
		// Only consider weekdays (Monday to Friday)
		if currentDay.Weekday() >= time.Monday && currentDay.Weekday() <= time.Friday {
			if currentDay.Weekday() == date.Weekday() &&
				currentDay.Year() == date.Year() &&
				currentDay.Month() == date.Month() &&
				currentDay.Day() == date.Day() {
				return true
			}
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(dayCmd)
}

func newDayTemplate() *template.Template {
	const newDayTemplate = `# {{.Day}}

## Morning Checklist

- [ ] check email
- [ ] check calendar
- [ ] check Slack
- [ ] check home todo
{{- if .ShowTimesheet }}
- [ ] time sheet
{{-  end }}
{{- if .ShowWorkingWednesday}}
- [ ] working Wednesday
{{-  end }}
{{- if .ShowExpenseTodo}}
- [ ] WFH expenses in Concur
{{-  end }}

## What do you want to accomplish today?

- [ ]
`

	return template.Must(template.New("newDayTemplate").Parse(newDayTemplate))
}
