package cmd

import (
	"fmt"
	"os"
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

		editor := NvimEditor{} // Or inject this as a dependency
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
	showExpenseTodo := isLastWednesday(timeNow)
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

	filePath := fmt.Sprintf("%s/%s/2024_Q4/%d-%d-%d.md", cfg.VaultPath, cfg.DayPath, timeNow.Month(), timeNow.Day(), timeNow.Year())

	_, err = os.Stat(filePath)
	if err == nil {
		// File exists, don't overwrite (or handle differently)
		return filePath, nil
	}

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

func isLastWednesday(date time.Time) bool {
	// Get the last day of the month
	lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())

	// Find the last Wednesday
	for lastDay.Weekday() != time.Wednesday {
		lastDay = lastDay.AddDate(0, 0, -1)
	}

	// Compare the given date with the last Wednesday
	return date.Year() == lastDay.Year() &&
		date.Month() == lastDay.Month() &&
		date.Day() == lastDay.Day()
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
