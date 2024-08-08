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

func catFileCommand() *cobra.Command {
	var objectType string
	var object string

	catFileCmd := &cobra.Command{
		Use:   "cat-file",
		Short: "Provide content of repository objects",
		RunE: func(command *cobra.Command, args []string) error {
			repo, err := cmd.LocateGitRepository(".", true)
			if err != nil {
				return fmt.Errorf("unable to locate repository: %w", err)
			}

			om := cmd.NewObjectManager(repo)
			obj, err := om.ReadObject(object)
			if err != nil {
				return fmt.Errorf("failed to read object: %w", err)
			}

			data, err := obj.Serialize()
			if err != nil {
				return fmt.Errorf("failed to serialize object: %w", err)
			}

			fmt.Println(string(data))
			return nil
		},
	}

	catFileCmd.Flags().StringVarP(&objectType, "type", "t", "", "Specify the type (blob, commit, tag, tree)")
	_ = catFileCmd.MarkFlagRequired("type")

	catFileCmd.Flags().StringVarP(&object, "object", "o", "", "The object to display")
	_ = catFileCmd.MarkFlagRequired("object")

	return catFileCmd
}

func hashObjectCommand() *cobra.Command {
	var objectType string
	var write bool
	var filePath string

	hashObjectCmd := &cobra.Command{
		Use:   "hash-object",
		Short: "Compute object ID and optionally creates a blob from a file",
		RunE: func(command *cobra.Command, args []string) error {
			var repo *cmd.GitRepository
			var err error
			if write {
				repo, err = cmd.LocateGitRepository(".", true)
				if err != nil {
					return fmt.Errorf("unable to locate repository: %w", err)
				}
			}

			om := cmd.NewObjectManager(repo)
			obType, err := cmd.TypeFromString(objectType)

			if err != nil {
				return fmt.Errorf("invalid object type: %w", err)
			}

			hash, err := om.HashObject(filePath, obType, write)
			if err != nil {
				return fmt.Errorf("failed to hash object: %w", err)
			}

			fmt.Println(hash)
			return nil
		},
	}

	hashObjectCmd.Flags().StringVarP(&objectType, "type", "t", "blob", "Specify the type (blob, commit, tag, tree)")

	hashObjectCmd.Flags().BoolVarP(&write, "write", "w", false, "Actually write the object into the database")

	hashObjectCmd.Flags().StringVarP(&filePath, "path", "p", "", "Read object from <file>")

	_ = hashObjectCmd.MarkFlagRequired("path")

	return hashObjectCmd
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "justdoit",
		Short: "It is a simple CLI application to manage your tasks.",
	}

	initCmd := initCommand()
	catFileCmd := catFileCommand()
	hashObjCmd := hashObjectCommand()
	rootCmd.AddCommand(initCmd, catFileCmd, hashObjCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
