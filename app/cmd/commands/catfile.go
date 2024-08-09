package commands

import (
	"fmt"
	"github.com/utkarsh5026/justdoit/app/cmd/objects"
	"github.com/utkarsh5026/justdoit/app/cmd/repository"
)

type CatFileOptions struct {
	Type bool // A boolean indicating whether to display the object type.
}

// CatFile provides content or type information for a Git object in the current repository.
// eg : git cat-file -t <object>, git cat-file -p <object>
// Parameters:
// - object: A string representing the SHA-1 hash of the object to read.
// - options: A CatFileOptions struct that specifies whether to print the type of the object.
//
// Returns:
// - An error if any operation fails, otherwise nil.
func CatFile(object string, options CatFileOptions) error {
	repo, err := repository.LocateCurrentRepository()
	if err != nil {
		return RepoNotFound(err)
	}

	om := objects.NewObjectManager(repo)

	commitObj, err := om.ReadObject(object)
	if err != nil {
		return ObjectReadError(err)
	}

	if options.Type {
		fmt.Println(commitObj.Format().String())
		return nil
	}

	data, err := commitObj.Serialize()
	if err != nil {
		return SerializationError(err)
	}

	fmt.Println(string(data))
	return nil
}
