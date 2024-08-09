package commands

import (
	"fmt"
	"github.com/utkarsh5026/justdoit/app/cmd/objects"
	"github.com/utkarsh5026/justdoit/app/cmd/repository"
	"path/filepath"
)

func LsTree(recursive bool, treeSha string) error {
	repo, err := repository.LocateCurrentRepository()
	if err != nil {
		return err
	}
	return lsTree(repo, recursive, treeSha, "")
}

// lsTree lists the contents of a Git tree object.
//
// Parameters:
// - repo: A pointer to the GitRepository object representing the current repository.
// - recursive: A boolean indicating whether to list contents recursively.
// - treeSha: A string representing the SHA-1 hash of the tree object to list.
// - prefix: A string representing the prefix path for the entries.
//
// Returns:
// - An error if any operation fails, otherwise nil.
func lsTree(repo *repository.GitRepository, recursive bool, treeSha string, prefix string) error {

	repo, err := repository.LocateCurrentRepository()
	if err != nil {
		return fmt.Errorf("unable to locate repository: %w", err)
	}

	oman := objects.NewObjectManager(repo)
	sha := oman.FindObject(treeSha, objects.TreeType, true)

	obj, err := oman.ReadObject(sha)
	if err != nil {
		return fmt.Errorf("failed to read tree object: %w", err)
	}

	tree, ok := obj.(*objects.GitTree)
	if !ok {
		return fmt.Errorf("invalid tree object")
	}

	entries := tree.Entries()
	for _, entry := range entries {
		entryType, err := entry.Type()
		if err != nil {
			return fmt.Errorf("failed to get object type: %w", err)
		}

		if !(recursive && entryType == objects.TreeType) {
			printTreeEntry(prefix, entryType, entry)
		} else {
			prefix := filepath.Join(prefix, entry.Name())
			if err := lsTree(repo, recursive, entry.Sha(), prefix); err != nil {
				return err
			}
		}
	}

	return nil
}

// printTreeEntry prints the details of a Git tree entry, int the format:
// <object type> <mode> <sha> <path>
//
// Parameters:
// - prefix: A string representing the prefix path for the entry.
// - objType: The type of the Git object (e.g., blob, tree, commit).
// - entry: A pointer to the GitTreeLeaf object representing the tree entry.
func printTreeEntry(prefix string, objType objects.GitObjectType, entry *objects.GitTreeLeaf) {
	mode := entry.Mode()
	sha := entry.Sha()
	path := filepath.Join(prefix, entry.Name())
	fmt.Printf("%s %s %s %s\n", objType.String(), mode, sha, path)
}