/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

/*
* This does:
* - create new folder in work/01-projects
* - create main note in new folder
* - create an Investigation note
* - create an Estimate note with a DATE as the name of the file, so that when I pass the date I can be reminded to reevaluate my approach.
* - create an TODO note
 */

type TicketArgs struct {
	Ticket   string
	Estimate int
}

// ticketCmd represents the ticket command
var ticketCmd = &cobra.Command{
	Use:   "ticket",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
  and usage of using your command. For example:

  Cobra is a CLI library for Go that empowers applications.
  This application is a tool to generate the needed files
  to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ticket   string
			estimate int
		)
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("What is the ticket number?").
					Value(&ticket), // store the chosen option in the "burger" variable
				huh.NewSelect[int]().
					Title("How much work will this take?").
					Options(
						huh.NewOption("None", 0),
						huh.NewOption("A little", 1),
						huh.NewOption("A lot", 3),
					).
					Value(&estimate),
			),
		)
		err := form.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		makeTicket(TicketArgs{ticket, estimate})
	},
}

func init() {
	rootCmd.AddCommand(ticketCmd)
}

func makeTicket(ticketArgs TicketArgs) error {
	timeNow := time.Now()
	dueDate := timeNow.AddDate(0, 0, ticketArgs.Estimate)
	formattedDueDate := fmt.Sprintf("%s, %d %s %d\n", dueDate.Weekday(), dueDate.Day(), dueDate.Month().String(), dueDate.Year())
	projectPath := fmt.Sprintf("/Users/gb0218/vaults/work/01-projects/%s", ticketArgs.Ticket)
	descPath := fmt.Sprintf("%s/%s.md", projectPath, ticketArgs.Ticket)
	estimatePath := fmt.Sprintf("%s/%s.md", projectPath, strings.TrimRight(formattedDueDate, " \t\n"))
	todoPath := fmt.Sprintf("%s/TODO.md", projectPath)
	investigationPath := fmt.Sprintf("%s/Investigation.md", projectPath)

	err := createProjectFolder(projectPath)
	if err != nil {
		return err
	}
	err = writeProjectFile(todoTemplate(), ticketArgs, todoPath)
	if err != nil {
		return err
	}
	err = writeProjectFile(descTemplate(), ticketArgs, descPath)
	if err != nil {
		return err
	}
	err = writeProjectFile(estimateTemplate(), ticketArgs, estimatePath)
	if err != nil {
		return err
	}
	err = writeProjectFile(investigationTemplate(), ticketArgs, investigationPath)
	if err != nil {
		return err
	}
	return nil
}

func investigationTemplate() *template.Template {
	const tmpl = `# [[{{.Ticket}}]] - Investigation

## Related Code

## Services involved

## Tests Needed

`
	return template.Must(template.New("DescTemplate").Parse(tmpl))
}

func estimateTemplate() *template.Template {
	const tmpl = `# [[{{.Ticket}}]] - Estimate: {{.Estimate}}/3

## If this ticket was not completed by the date estimated, please describe why.

Were there interruptions?:

Was there something that you did not understand?: 

Does the code require a refactor before continuing?: 

Do you need help from another team member?: 

`
	return template.Must(template.New("DescTemplate").Parse(tmpl))
}

func descTemplate() *template.Template {
	const tmpl = `# [[{{.Ticket}}]]

## Branch

gb/your-branch-name-here

## Description

`
	return template.Must(template.New("DescTemplate").Parse(tmpl))
}

func todoTemplate() *template.Template {
	const todoTemplate = `# [[{{.Ticket}}]] - TODO

## TODO

  - [ ] Describe work to be done
  - [ ] Make a feature branch
  - [ ] 

`
	return template.Must(template.New("TodoTemplate").Parse(todoTemplate))
}

func writeProjectFile(ticketT *template.Template, ticketArgs TicketArgs, fpath string) error {
	var file *os.File
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		file, _ = os.Create(fpath)
		defer file.Close()
		return ticketT.Execute(file, ticketArgs)
	} else {
		return nil
	}
}

func createProjectFolder(projectPath string) error {
	err := os.Mkdir(projectPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
