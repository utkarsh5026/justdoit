package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/utkarsh5026/justdoit/app/cmd/commands"
	"github.com/utkarsh5026/justdoit/app/cmd/repository"
)

func initCommand() *cobra.Command {
	var repoPath string
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create an empty Git repository or reinitialize an existing one",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			_, err := repository.CreateGitRepository(repoPath)
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
		Long: "The 'cat-file' command provides content of repository objects. " +
			"It can be used to display the content of a blob, commit, tag, or tree object.",
		RunE: func(command *cobra.Command, args []string) error {
			options := commands.CatFileOptions{
				Type: objectType != "",
			}
			return commands.CatFile(objectType, options)
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
			return commands.HashObject(filePath, objectType, write)
		},
	}

	hashObjectCmd.Flags().StringVarP(&objectType, "type", "t", "blob", "Specify the type (blob, commit, tag, tree)")

	hashObjectCmd.Flags().BoolVarP(&write, "write", "w", false, "Actually write the object into the database")

	hashObjectCmd.Flags().StringVarP(&filePath, "path", "p", "", "Read object from <file>")

	_ = hashObjectCmd.MarkFlagRequired("path")

	return hashObjectCmd
}

func logCommand() *cobra.Command {
	var commit string

	logCmd := &cobra.Command{
		Use:   "log",
		Short: "Display history of a given commit.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				commit = args[0]
			} else {
				commit = "HEAD"
			}
			// Add your logic here to handle the log command
			return nil
		},
	}

	logCmd.Flags().StringVarP(&commit, "commit", "c", "HEAD", "Commit to start at")

	return logCmd
}

func lsTreeCommand() *cobra.Command {

	var recursive bool
	var treeSha string

	lsTreeCmd := &cobra.Command{
		Use:   "ls-tree",
		Short: "List the contents of a tree object",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.LsTree(recursive, treeSha)
		},
	}

	lsTreeCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recurse into sub-trees")

	lsTreeCmd.Flags().StringVarP(&treeSha, "tree", "t", "HEAD", "The tree to list")
	return lsTreeCmd
}

func checkoutCommand() *cobra.Command {
	var commit string
	var path string
	checkoutCmd := &cobra.Command{
		Use:   "checkout",
		Short: "Checkout a commit inside of a directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.Checkout(commit, path)
		},
	}

	checkoutCmd.Flags().StringVarP(&commit, "commit", "c", "", "The commit or tree to checkout.")
	_ = checkoutCmd.MarkFlagRequired("commit")

	checkoutCmd.Flags().StringVarP(&path, "path", "p", "", "The EMPTY directory to checkout on.")
	_ = checkoutCmd.MarkFlagRequired("path")

	return checkoutCmd
}

func showRefCommand() *cobra.Command {
	showRefCmd := &cobra.Command{
		Use:   "show-ref",
		Short: "List references in a local repository.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.ShowRef()
		},
	}

	return showRefCmd
}

func tagCommand() *cobra.Command {
	var createTagObject bool
	var name string
	var object string

	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "List and create tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				name = args[0]
			}
			if len(args) > 1 {
				object = args[1]
			} else {
				object = "HEAD"
			}
			// Add your logic here to handle the tag command
			return nil
		},
	}

	tagCmd.Flags().BoolVarP(&createTagObject, "create_tag_object", "a", false, "Whether to create a tag object")
	tagCmd.Flags().StringVarP(&name, "name", "n", "", "The new tag's name")
	tagCmd.Flags().StringVarP(&object, "object", "o", "HEAD", "The object the new tag will point to")

	return tagCmd
}

func lsFilesCommand() *cobra.Command {
	var verbose bool
	lsFilesCmd := &cobra.Command{
		Use:   "ls-files",
		Short: "List files in the index",
		RunE: func(cmd *cobra.Command, args []string) error {
			options := commands.LsFilesOptions{
				Verbose: verbose,
			}
			return commands.LsFile(options)
		},
	}

	lsFilesCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show additional information about each file")

	return lsFilesCmd
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "justdoit",
		Short: "It is a simple CLI application to manage your tasks.",
	}

	initCmd := initCommand()
	catFileCmd := catFileCommand()
	hashObjCmd := hashObjectCommand()
	logCmd := logCommand()
	lsTreeCmd := lsTreeCommand()
	checkoutCmd := checkoutCommand()
	showRefCmd := showRefCommand()
	tagCmd := tagCommand()
	lsFilesCmd := lsFilesCommand()
	rootCmd.AddCommand(initCmd,
		catFileCmd,
		hashObjCmd,
		logCmd,
		lsTreeCmd,
		checkoutCmd,
		showRefCmd,
		tagCmd,
		lsFilesCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
