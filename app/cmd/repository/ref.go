package repository

import (
	"github.com/utkarsh5026/justdoit/app/cmd/fileutils"
	"github.com/utkarsh5026/justdoit/app/ordereddict"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ListRefs(repo *GitRepository, path string) (*ordereddict.OrderedDict, error) {
	var err error
	if path == "" {
		path, err = EnsureGitDirExists(repo, false, "refs")
		if err != nil {
			return nil, err
		}
	}

	files, err := fileutils.ListDir(path)

	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return strings.Compare(files[i].Name(), files[j].Name()) < 0
	})

	refs := ordereddict.New()
	for _, file := range files {
		can := filepath.Join(path, file.Name())
		if file.IsDir() {
			subRefs, err := ListRefs(repo, can)
			if err != nil {
				return nil, err
			}

			refs.Set(file.Name(), subRefs)
		} else {
			ref, err := resolveRef(repo, can)
			if err != nil {
				return nil, err
			}

			refs.Set(file.Name(), ref)
		}
	}

	return refs, nil
}

// resolveRef resolves a reference file to its corresponding commit SHA or another reference.
//
// This function reads the content of the reference file and follows any symbolic references
// (i.e., references that start with "ref: ") recursively until it finds the actual commit SHA.
//
// Parameters:
// - repo: The Git repository object.
// - refFile: The path to the reference file.
//
// Returns:
// - A string containing the resolved reference (commit SHA or another reference).
// - An error if any operation fails.
func resolveRef(repo *GitRepository, refFile string) (string, error) {
	path := GetGitFilePath(repo, false, refFile)
	isFile, err := fileutils.IsFile(path)

	var ref string
	if err != nil {
		return ref, err
	}

	if !isFile {
		return ref, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ref, err
	}

	if len(data) > 0 {
		data = data[:len(data)-1]
	}

	if strings.HasPrefix(string(data), "ref: ") {
		ref = strings.TrimPrefix(string(data), "ref: ")
		return resolveRef(repo, ref)
	}

	return string(data), nil
}
