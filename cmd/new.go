/*
Copyright © 2023 Greg Barozzi barozzi@github.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

type NewDayArgs struct {
	Day                  string // Formatted Day for the header
	ShowTimesheet        bool
	ShowWorkingWednesday bool
	ShowExpenseTodo      bool
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new DevLog and open it in the default editor",
	Run: func(cmd *cobra.Command, args []string) {
		timeNow := time.Now()
		formattedDay := fmt.Sprintf("%s, %d %s %d\n", timeNow.Weekday(), timeNow.Day(), timeNow.Month().String(), timeNow.Year())
		showTimesheet := timeNow.Weekday() == time.Friday
		showWorkingWednesday := timeNow.Weekday() == time.Wednesday
		showExpenseTodo := isLastWednesday(timeNow)
		newDayArgs := NewDayArgs{formattedDay, showTimesheet, showWorkingWednesday, showExpenseTodo}

		newDayT := newDayTemplate()

		filePath, err := writeNewDay(newDayT, newDayArgs, timeNow)
		if err != nil {
			fmt.Println("Failed to write new-day template to file")
			os.Exit(1)
		}
		vim := exec.Command("zsh", "-c", fmt.Sprintf("nvim %s", filePath))
		vim.Stdin = os.Stdin
		vim.Stdout = os.Stdout
		vim.Stderr = os.Stderr
		err = vim.Start()
		if err != nil {
			fmt.Printf("Vim failed to start correctly: %s\n", err)
			os.Exit(1)
		}
		err = vim.Wait()
		if err != nil {
			fmt.Printf("Vim failed to exit correctly: %s\n", err)
			os.Exit(1)
		}
	},
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
	rootCmd.AddCommand(newCmd)
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

## What is preventing you?

- 

## What do you want to accomplish today?

- [ ]
`

	return template.Must(template.New("newDayTemplate").Parse(newDayTemplate))
}

func writeNewDay(newDayT *template.Template, newDayArgs NewDayArgs, timeNow time.Time) (newDataFilePath string, fileWriteError error) {
	// TODO - make the folder rotation be automated eg Q2 to Q3
	filePath := fmt.Sprintf("/Users/gb0218/vaults/work/00-dev-log/2024_Q4/%d-%d-%d.md", timeNow.Month(), timeNow.Day(), timeNow.Year())
	fmt.Printf("Filepath: %s\n", filePath)
	var file *os.File

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, _ = os.Create(filePath)
		defer file.Close()
		return filePath, newDayT.Execute(file, newDayArgs)
	} else {
		return filePath, nil
	}
}
