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
	Day           string // Formatted Day for the header
	ShowTimesheet bool
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new DevLog and open it in the default editor",
	Run: func(cmd *cobra.Command, args []string) {
		timeNow := time.Now()
		formattedDay := fmt.Sprintf("%s, %d %s %d\n", timeNow.Weekday(), timeNow.Day(), timeNow.Month().String(), timeNow.Year())
		showTimesheet := timeNow.Weekday() == time.Friday
		newDayArgs := NewDayArgs{formattedDay, showTimesheet}

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

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func newDayTemplate() *template.Template {
	const newDayTemplate = `# {{.Day}}

## Morning Checklist

- [ ] check email
- [ ] check calendar
- [ ] check Slack
{{- if .ShowTimesheet }}
- [ ] time sheet
{{-  end }}

## Resistance

Anything that slows you down or you find frustrating

- 

## TODO

- [ ]
`

	return template.Must(template.New("newDayTemplate").Parse(newDayTemplate))
}

func writeNewDay(newDayT *template.Template, newDayArgs NewDayArgs, timeNow time.Time) (newDataFilePath string, fileWriteError error) {
	// reposPath := os.Getenv("REPOS_PATH")
	// TODO - make the folder rotation be automated eg Q2 to Q3
	filePath := fmt.Sprintf("/Users/gb0218/vaults/work/00-dev-log/2024_Q3/%d-%d-%d.md", timeNow.Month(), timeNow.Day(), timeNow.Year())
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
