package cmd

import (
	"fmt"
	"gnote/config"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// archiveCmd represents the archive command
var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a project",
	Long:  `Moves a project folder from the projects directory to the archive directory, organized by quarter.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.ReadConfig()
		if err != nil {
			fmt.Println("Error reading config:", err)
			return
		}

		projectsPath := filepath.Join(cfg.VaultPath, cfg.ProjectsPath)
		archivePath := filepath.Join(cfg.VaultPath, cfg.ArchivesPath)

		// Get the current year and quarter
		timeNow := time.Now()
		year := timeNow.Year()
		quarter := getQuarter(timeNow)
		quarterFolder := fmt.Sprintf("%d_Q%d", year, quarter)

		// Create the quarter folder in the archive path
		quarterArchivePath := filepath.Join(archivePath, quarterFolder)
		if _, err := os.Stat(quarterArchivePath); os.IsNotExist(err) {
			err = os.MkdirAll(quarterArchivePath, 0755)
			if err != nil {
				fmt.Println("Error creating archive quarter directory:", err)
				return
			}
		}

		// List project folders
		projectFolders, err := listProjectFolders(projectsPath)
		if err != nil {
			fmt.Println("Error listing project folders:", err)
			return
		}

		if len(projectFolders) == 0 {
			fmt.Println("No projects found to archive.")
			return
		}

		// User selection
		var selectedFolder string
		form := huh.NewForm(
			huh.NewGroup( // Wrap the select in a group
				huh.NewSelect[string]().
					Title("Select project to archive:").
					Options(generateHuhOptions(projectFolders)...).
					Value(&selectedFolder),
			),
		)
		err = form.Run()

		if err != nil {
			fmt.Println("Error during selection:", err)
			return
		}

		// Move folder
		sourcePath := filepath.Join(projectsPath, selectedFolder)
		destPath := filepath.Join(quarterArchivePath, selectedFolder) // Use quarterArchivePath

		err = os.Rename(sourcePath, destPath)
		if err != nil {
			fmt.Printf("Error moving project '%s': %v\n", selectedFolder, err)
			return
		}

		fmt.Printf("Project '%s' archived successfully to '%s'\n", selectedFolder, quarterArchivePath)
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}

func listProjectFolders(projectsPath string) ([]string, error) {
	var folders []string
	files, err := os.ReadDir(projectsPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}

	return folders, nil
}

func generateHuhOptions(folders []string) []huh.Option[string] {
	options := make([]huh.Option[string], len(folders))
	for i, folder := range folders {
		options[i] = huh.NewOption(folder, folder)
	}
	return options
}
