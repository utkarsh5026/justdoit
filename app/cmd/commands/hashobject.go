package commands

import (
	"fmt"
	"github.com/utkarsh5026/justdoit/app/cmd/objects"
	"github.com/utkarsh5026/justdoit/app/cmd/repository"
)

// HashObject hashes a file and optionally writes the object to the repository.
//
// This function locates the current repository if the write option is enabled,
// determines the object type from the provided string, and hashes the file.
// If the write option is enabled, the object is written to the repository.
//
// eg : git hash-object -w README.md  | git hash-object -w README.md -t blob
//
// Parameters:
// - filePath: A string representing the path to the file to be hashed.
// - objectType: A string representing the type of the object (e.g., "blob", "tree", "commit").
// - write: A boolean indicating whether to write the object to the repository.
//
// Returns:
// - An error if any operation fails, otherwise nil.
func HashObject(filePath string, objectType string, write bool) error {
	var repo *repository.GitRepository
	var err error

	if write {
		repo, err = repository.LocateCurrentRepository()
	}

	if err != nil {
		return RepoNotFound(err)
	}

	om := objects.NewObjectManager(repo)
	obType, err := objects.TypeFromString(objectType)

	if err != nil {
		return ObjectTypeError(err)
	}

	hash, err := om.HashObject(filePath, obType, write)
	if err != nil {
		return fmt.Errorf("failed to hash object: %w", err)
	}

	fmt.Println(hash)
	return nil
}
