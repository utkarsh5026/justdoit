package commands

import (
	"fmt"
	"github.com/utkarsh5026/justdoit/app/cmd/repository"
	"github.com/utkarsh5026/justdoit/app/ordereddict"
)

func ShowRef() error {
	repo, err := repository.LocateCurrentRepository()
	if err != nil {
		return RepoNotFound(err)
	}

	refs, err := repository.ListRefs(repo, "")
	return showRef(repo, refs, true, "")
}

func showRef(repo *repository.GitRepository, refs *ordereddict.OrderedDict, withHash bool, prefix string) error {
	for _, key := range refs.Keys() {
		value, _ := refs.Get(key)
		switch v := value.(type) {
		case string:
			if withHash {
				fmt.Printf("%s %s%s\n", v, prefix, key)
			} else {
				fmt.Printf("%s%s\n", prefix, key)
			}

		case *ordereddict.OrderedDict:
			newPrefix := prefix
			if prefix != "" {
				newPrefix += "/"
			}
			newPrefix += key
			err := showRef(repo, v, withHash, newPrefix)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("invalid ref type: %T", value)
		}
	}

	return nil
}
