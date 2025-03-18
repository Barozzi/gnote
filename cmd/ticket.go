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

type TicketArgs struct {
	Ticket   string
	Tag      string
	Link     string
	Estimate int
}

// UserInputCollector interface
type UserInputCollector interface {
	Collect() (TicketArgs, error)
}

// HuhInputCollector concrete implementation
type HuhInputCollector struct{}

func (h *HuhInputCollector) Collect() (TicketArgs, error) {
	var (
		ticket   string
		tag      string
		link     string
		estimate int
	)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is the ticket number?").
				Value(&ticket),
			huh.NewInput().
				Title("What Tag should this ticket use?").
				Value(&link),
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
		return TicketArgs{}, err
	}
	tag = strings.Replace(link, " ", "_", -1)
	return TicketArgs{ticket, tag, link, estimate}, nil
}

// FileGenerator interface
type FileGenerator interface {
	Generate(ticketArgs TicketArgs) error
}

type TemplateInfo struct {
	template *template.Template
	fpath    string
}

// TodoFileGenerator concrete implementation
type TodoFileGenerator struct{ TemplateInfo }

func (g *TodoFileGenerator) Generate(ticketArgs TicketArgs) error {
	return writeProjectFile(g.template, ticketArgs, g.TemplateInfo.fpath)
}

// DescFileGenerator concrete implementation
type DescFileGenerator struct{ TemplateInfo }

func (g *DescFileGenerator) Generate(ticketArgs TicketArgs) error {
	return writeProjectFile(g.template, ticketArgs, g.TemplateInfo.fpath)
}

// DescFileGenerator concrete implementation
type InvestigationFileGenerator struct{ TemplateInfo }

func (g *InvestigationFileGenerator) Generate(ticketArgs TicketArgs) error {
	return writeProjectFile(g.template, ticketArgs, g.TemplateInfo.fpath)
}

// DescFileGenerator concrete implementation
type EstimateFileGenerator struct{ TemplateInfo }

func (g *EstimateFileGenerator) Generate(ticketArgs TicketArgs) error {
	return writeProjectFile(g.template, ticketArgs, g.TemplateInfo.fpath)
}

// ProjectCreator struct
type ProjectCreator struct {
	projectPath    string
	fileGenerators []FileGenerator
}

func NewProjectCreator(projectPath string, fileGenerators []FileGenerator) *ProjectCreator {
	return &ProjectCreator{projectPath: projectPath, fileGenerators: fileGenerators}
}

func (pc *ProjectCreator) CreateProject(ticketArgs TicketArgs) error {
	err := createProjectFolder(pc.projectPath)
	if err != nil {
		return err
	}

	for _, generator := range pc.fileGenerators {
		err := generator.Generate(ticketArgs)
		if err != nil {
			return err
		}
	}

	return nil
}

// ticketCmd represents the ticket command
var ticketCmd = &cobra.Command{
	Use:   "ticket",
	Short: "A brief description of your command",
	Long:  `A longer description...`,
	Run: func(cmd *cobra.Command, args []string) {
		collector := &HuhInputCollector{} // Or another implementation
		ticketArgs, err := collector.Collect()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Paths are application specific and break the single responsibility principle
		timeNow := time.Now()
		dueDate := timeNow.AddDate(0, 0, ticketArgs.Estimate)
		formattedDueDate := fmt.Sprintf("%s, %d %s %d\n", dueDate.Weekday(), dueDate.Day(), dueDate.Month().String(), dueDate.Year())
		projectPath := fmt.Sprintf("/Users/gb0218/vaults/work/01-projects/%s", ticketArgs.Ticket)
		descPath := fmt.Sprintf("%s/%s.md", projectPath, ticketArgs.Ticket)
		estimatePath := fmt.Sprintf("%s/%s.md", projectPath, strings.TrimRight(formattedDueDate, " \t\n"))
		todoPath := fmt.Sprintf("%s/TODO.md", projectPath)
		investigationPath := fmt.Sprintf("%s/Investigation.md", projectPath)

		// Example with file paths passed into generator
		todoGenerator := &TodoFileGenerator{TemplateInfo{todoTemplate(), todoPath}}
		descGenerator := &DescFileGenerator{TemplateInfo{descTemplate(), descPath}}
		estimateGenerator := &EstimateFileGenerator{TemplateInfo{estimateTemplate(), estimatePath}}
		investigationGenerator := &InvestigationFileGenerator{TemplateInfo{investigationTemplate(), investigationPath}}

		creator := NewProjectCreator(projectPath, []FileGenerator{todoGenerator, descGenerator, estimateGenerator, investigationGenerator})

		err = creator.CreateProject(ticketArgs)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(ticketCmd)
}

func init() {
	rootCmd.AddCommand(ticketCmd)
}

func investigationTemplate() *template.Template {
	const tmpl = `# [[{{.Ticket}}]] - Investigation

## What is the problem you are trying to solve? (5Why)

## Related Code (filenames, or snippets)

## Which docs have you consulted?

## Describe this change from First Principals

### Core Components

### Key Considerations

`
	return template.Must(template.New("DescTemplate").Parse(tmpl))
}

func estimateTemplate() *template.Template {
	const tmpl = `# [[{{.Ticket}}]] - Estimate: {{.Estimate}}/3

## If this ticket was not completed by the date estimated, please describe why.

### Were there interruptions?:

### Was there something that you did not understand?: 

### Does the code require a refactor before continuing?: 

### Do you need help from another team member?: 

`
	return template.Must(template.New("DescTemplate").Parse(tmpl))
}

func descTemplate() *template.Template {
	const tmpl = `---
id: {{.Ticket}} 
aliases: 
tags:
  - '{{.Tag}}'
link: "[[{{.Link}}]]"
---

# [[{{.Ticket}}]]

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
  - [ ] Investigate
  - [ ] Make a feature branch

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
