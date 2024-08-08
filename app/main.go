package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/utkarsh5026/justdoit/app/cmd"
)

func initCommand() *cobra.Command {
	var repoPath string
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create an empty Git repository or reinitialize an existing one",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			_, err := cmd.CreateGitRepository(repoPath)
			if err != nil {
				return err
			}

			fmt.Println("Initialized empty Git repository in", repoPath)
			return nil
		},
	}

	initCmd.Flags().StringVarP(&repoPath, "path",
		"p", ".", "The path to the repository")
	return initCmd
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "justdoit",
		Short: "It is a simple CLI application to manage your tasks.",
	}

	initCmd := initCommand()
	rootCmd.AddCommand(initCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
