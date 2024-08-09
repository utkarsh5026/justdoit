package commands

import (
	"fmt"
	"github.com/utkarsh5026/justdoit/app/cmd/objects"
	"github.com/utkarsh5026/justdoit/app/cmd/repository"
	"github.com/utkarsh5026/justdoit/app/ordereddict"
)

type TagOptions struct {
	Annotated bool
	Force     bool
	Message   string
	Tagger    string
}

func Tag(name string, options TagOptions) error {
	repo, err := repository.LocateCurrentRepository()
	if err != nil {
		return RepoNotFound(err)
	}

	if name == "" {
		refs, err := repository.ListRefs(repo, "")
		if err != nil {
			return err
		}

		tags, ok := refs.Get("tags")
		if !ok {
			fmt.Println("No tags found.")
			return nil
		}

		tagRef, ok := tags.(*ordereddict.OrderedDict)
		if !ok {
			return fmt.Errorf("invalid tags reference")
		}
		return showRef(repo, tagRef, false, "")
	}

	return objects.CreateTag(repo, name, "", options.Annotated, options.Tagger, options.Message)
}
