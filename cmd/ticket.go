package cmd

import (
	"fmt"
	"gnote/config"
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
	Generate(ticketArgs TicketArgs, cfg *config.Config) error
}

type TemplateInfo struct {
	template *template.Template
}

// TodoFileGenerator concrete implementation
type TodoFileGenerator struct{ TemplateInfo }

func (g *TodoFileGenerator) Generate(ticketArgs TicketArgs, cfg *config.Config) error {
	projectPath := fmt.Sprintf("%s/%s/%s", cfg.VaultPath, cfg.ProjectsPath, ticketArgs.Ticket)
	todoPath := fmt.Sprintf("%s/TODO.md", projectPath)
	return writeProjectFile(g.template, ticketArgs, todoPath)
}

// DescFileGenerator concrete implementation
type DescFileGenerator struct{ TemplateInfo }

func (g *DescFileGenerator) Generate(ticketArgs TicketArgs, cfg *config.Config) error {
	projectPath := fmt.Sprintf("%s/%s/%s", cfg.VaultPath, cfg.ProjectsPath, ticketArgs.Ticket)
	descPath := fmt.Sprintf("%s/%s.md", projectPath, ticketArgs.Ticket)
	return writeProjectFile(g.template, ticketArgs, descPath)
}

// DescFileGenerator concrete implementation
type InvestigationFileGenerator struct{ TemplateInfo }

func (g *InvestigationFileGenerator) Generate(ticketArgs TicketArgs, cfg *config.Config) error {
	projectPath := fmt.Sprintf("%s/%s/%s", cfg.VaultPath, cfg.ProjectsPath, ticketArgs.Ticket)
	investigationPath := fmt.Sprintf("%s/Investigation.md", projectPath)
	return writeProjectFile(g.template, ticketArgs, investigationPath)
}

// DescFileGenerator concrete implementation
type EstimateFileGenerator struct{ TemplateInfo }

func (g *EstimateFileGenerator) Generate(ticketArgs TicketArgs, cfg *config.Config) error {
	timeNow := time.Now()
	dueDate := timeNow.AddDate(0, 0, ticketArgs.Estimate)
	formattedDueDate := fmt.Sprintf("%s, %d %s %d\n", dueDate.Weekday(), dueDate.Day(), dueDate.Month().String(), dueDate.Year())
	projectPath := fmt.Sprintf("%s/%s/%s", cfg.VaultPath, cfg.ProjectsPath, ticketArgs.Ticket)
	estimatePath := fmt.Sprintf("%s/%s.md", projectPath, strings.TrimRight(formattedDueDate, " \t\n"))
	return writeProjectFile(g.template, ticketArgs, estimatePath)
}

// ProjectCreator struct
type ProjectCreator struct {
	cfg            *config.Config
	fileGenerators []FileGenerator
}

func NewProjectCreator(cfg *config.Config, fileGenerators []FileGenerator) *ProjectCreator {
	return &ProjectCreator{cfg: cfg, fileGenerators: fileGenerators}
}

func (pc *ProjectCreator) CreateProject(ticketArgs TicketArgs) error {
	projectPath := fmt.Sprintf("%s/%s/%s", pc.cfg.VaultPath, pc.cfg.ProjectsPath, ticketArgs.Ticket)
	err := createProjectFolder(projectPath)
	if err != nil {
		return err
	}

	for _, generator := range pc.fileGenerators {
		err := generator.Generate(ticketArgs, pc.cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

// ticketCmd represents the ticket command
var ticketCmd = &cobra.Command{
	Use:   "ticket",
	Short: "Create a new project folder populated with default files",
	Long: `Follow a prompt to create a project folder containing:
  1. Project description file
  2. Investigation file
  3. TODO file
  4. Possibly, Estimate file
  `,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.ReadConfig()
		if err != nil {
			fmt.Println("Error reading config:", err)
			return
		}
		collector := &HuhInputCollector{}
		ticketArgs, err := collector.Collect()
		if err != nil {
			fmt.Println(err)
			return
		}

		todoGenerator := &TodoFileGenerator{TemplateInfo{todoTemplate()}}
		descGenerator := &DescFileGenerator{TemplateInfo{descTemplate()}}
		estimateGenerator := &EstimateFileGenerator{TemplateInfo{estimateTemplate()}}
		investigationGenerator := &InvestigationFileGenerator{TemplateInfo{investigationTemplate()}}

		creator := NewProjectCreator(cfg, []FileGenerator{todoGenerator, descGenerator, estimateGenerator, investigationGenerator})

		err = creator.CreateProject(ticketArgs)
		if err != nil {
			fmt.Println(err)
		}
	},
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
