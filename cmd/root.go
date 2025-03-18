/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnote",
	Short: "GNote is a tool for quick note-taking and project organization.",
	Long: `
 ▗▄▄▖▗▖  ▗▖ ▗▄▖▗▄▄▄▖▗▄▄▄▖
▐▌   ▐▛▚▖▐▌▐▌ ▐▌ █  ▐▌   
▐▌▝▜▌▐▌ ▝▜▌▐▌ ▐▌ █  ▐▛▀▀▘
▝▚▄▞▘▐▌  ▐▌▝▚▄▞▘ █  ▐▙▄▄▖
                         
  GNote helps you manage your notes, dev logs, and projects 
by providing commands to quickly create and organize files. 
Use 'gnote [command] --help' for more information about a specific command.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.new-day.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
