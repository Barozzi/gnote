/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <searc-string>",
	Short: "Search all DevLogs for a string match",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchString := args[0]
		reposPath := os.Getenv("REPOS_PATH")
		vim := exec.Command("zsh", "-c", fmt.Sprintf("vim -c \"cd %s/DevLog\" -c \"Ack %s\"", reposPath, searchString))
		vim.Stdin = os.Stdin
		vim.Stdout = os.Stdout
		vim.Stderr = os.Stderr
		err := vim.Start()
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
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
