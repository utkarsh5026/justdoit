package commands

import (
	"fmt"
	"github.com/utkarsh5026/justdoit/app/cmd/fileutils"
	"github.com/utkarsh5026/justdoit/app/cmd/objects"
	"github.com/utkarsh5026/justdoit/app/cmd/repository"
	"os"
	"path/filepath"
)

// Checkout checks out a specific commit to the given path.
//
// This function locates the current repository, reads the commit object,
// and checks out the tree associated with the commit to the specified path.
// It ensures the path is a directory and is empty before proceeding.
//
// Parameters:
// - commit: The commit SHA to check out.
// - path: The destination path where the commit should be checked out.
//
// Returns:
// - An error if any operation fails, otherwise nil.
func Checkout(commit string, path string) error {
	repo, err := repository.LocateCurrentRepository()
	if err != nil {
		return RepoNotFound(err)
	}

	om := objects.NewObjectManager(repo)
	commitSha := om.FindObject(commit, objects.CommitType, false)

	object, err := om.ReadObject(commitSha)
	if err != nil {
		return ObjectReadError(err)
	}

	commitObj, ok := object.(*objects.CommitObject)
	if ok {
		treeSha := commitObj.GetCommit().Tree
		object, err = om.ReadObject(treeSha)

		if err != nil {
			return ObjectReadError(err)
		}
	}

	pathExists := fileutils.PathExists(path)
	if pathExists {
		isDir, err := fileutils.IsDir(path)
		if err != nil {
			return err
		}

		if !isDir {
			return fmt.Errorf("'%s' is not a directory", path)
		}

		dirs, err := fileutils.ListDir(path)
		if err != nil {
			return err
		}
		if len(dirs) > 0 {
			return fmt.Errorf("'%s' is not empty", path)
		}
	}

	return checkoutTree(om, object.(*objects.GitTree), path)
}

// checkoutTree recursively checks out a Git tree object to the specified path.
//
// This function reads the entries of a Git tree object and processes each entry based on its type.
// It handles directories (trees), files (blobs), commits, and tags.
//
// Parameters:
// - om: A pointer of objects.ObjectManager used to read objects from the repository.
// - tree: The *objects.GitTree object representing the tree to be checked out.
// - path: The destination path where the tree should be checked out.
//
// Returns:
// - An error if any operation fails, otherwise nil.
func checkoutTree(om *objects.ObjectManager, tree *objects.GitTree, path string) error {
	entries := tree.Entries()
	for _, entry := range entries {
		object, err := om.ReadObject(entry.Sha())
		if err != nil {
			return ObjectReadError(err)
		}

		dest := filepath.Join(path, entry.Name())

		switch object.Format() {
		case objects.TreeType:
			if err := os.Mkdir(dest, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			err = checkoutTree(om, object.(*objects.GitTree), dest)

		case objects.BlobType:
			data, err := object.Serialize()
			if err != nil {
				return SerializationError(err)
			}

			err = os.WriteFile(dest, data, 0644)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
		case objects.CommitType:
		case objects.TagType:
		}
	}

	return nil
}
