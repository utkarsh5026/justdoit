package commands

import (
	"fmt"
	"github.com/utkarsh5026/justdoit/app/cmd/objects"
	"github.com/utkarsh5026/justdoit/app/cmd/repository"
	"os/user"
	"strconv"
	"time"
)

type LsFilesOptions struct {
	Verbose bool
}

var IndexReadError = func(err error) error {
	return fmt.Errorf("failed to read index: %w", err)
}

func LsFile(opts LsFilesOptions) error {
	repo, err := repository.LocateCurrentRepository()
	if err != nil {
		return RepoNotFound(err)
	}

	index, err := objects.ReadIndex(repo)
	if err != nil {
		return IndexReadError(err)
	}

	if opts.Verbose {
		fmt.Printf("Index version: %d\n, containing %d entries\n", index.Version, len(index.Entries))
	}

	for _, e := range index.Entries {
		fmt.Println(e.Name)

		if opts.Verbose {
			fileType := e.ModeType.String()
			fmt.Printf("\t%s with perms: %o\n", fileType, e.ModePerms)
		}

		fmt.Printf("  created: %s.%d, modified: %s.%d\n",
			time.Unix(e.Ctime[0], e.Ctime[1]).Format(time.RFC3339),
			e.Ctime[1],
			time.Unix(e.Mtime[0], e.Mtime[1]).Format(time.RFC3339),
			e.Mtime[1])

		uid := strconv.Itoa(int(e.UserId))
		gid := strconv.Itoa(int(e.GroupId))

		usr, err := user.LookupId(uid)
		if err != nil {
			fmt.Printf("  user: %s (%d)\n", "unknown", e.UserId)
		} else {
			fmt.Printf("  user: %s (%d)\n", usr.Username, e.UserId)
		}

		group, err := user.LookupGroupId(gid)
		if err != nil {
			fmt.Printf("  group: %s (%d)\n", "unknown", e.GroupId)
		} else {
			fmt.Printf("  group: %s (%d)\n", group.Name, e.GroupId)
		}

		fmt.Printf("  flags: stage=%d assume_valid=%t\n", e.FlagStage, e.FlagAssumeValid)
	}

	return nil
}
