package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateRepoPath constructs a file path by joining the repository path with additional paths.
//
// Parameters:
// - repository: A pointer to a GitRepository struct containing the repository paths.
// - paths: A variadic parameter representing additional path segments to be joined.
//
// Returns:
// - A string representing the combined file path.
func CreateRepoPath(repo *GitRepository, paths ...string) string {
	paths = append([]string{repo.GitDir}, paths...)
	return filepath.Join(paths...)
}

// EnsureGitDirExists constructs a directory path within a repository and optionally creates the directory.
//
// Parameters:
// - repository: A pointer to a GitRepository struct containing the repository paths.
// - mkdir: A boolean indicating whether to create the directory if it does not exist.
// - paths: A variadic parameter representing additional path segments to be joined.
//
// Returns:
// - A string representing the combined directory path. If an error occurs, an empty string is returned.
// - An error if there is an issue checking the directory status or creating the directory
func EnsureGitDirExists(repo *GitRepository, mkdir bool, paths ...string) (string, error) {
	path := CreateRepoPath(repo, paths...)
	pathExists := true

	if _, err := os.Stat(path); os.IsNotExist(err) {
		pathExists = false
	}

	if pathExists {
		isDir, err := isDir(path)
		if err != nil || !isDir {
			return "", err
		} else {
			return path, nil
		}
	}

	if mkdir {
		if err := os.MkdirAll(path, 0755); err != nil {
			return "", err
		}
		return path, nil
	}

	return "", nil
}

// GetGitFilePath constructs a file path within a repository and optionally creates the necessary directories.
//
// Parameters:
// - repository: A pointer to a GitRepository struct containing the repository paths.
// - mkdir: A boolean indicating whether to create the directory if it does not exist.
// - paths: A variadic parameter representing additional path segments to be joined.
//
// Returns:
// - A string representing the combined file path. If an error occurs, an empty string is returned
func GetGitFilePath(repo *GitRepository, mkdir bool, paths ...string) string {
	dirPath := paths[:len(paths)-1]
	_, err := EnsureGitDirExists(repo, mkdir, dirPath...)
	if err != nil {
		return ""
	}
	return CreateRepoPath(repo, paths...)
}

// LocateGitRepository searches for a Git repository directory starting from a given path.
//
// Parameters:
// - startPath: A string representing the starting path for the search.
// - required: A boolean indicating whether the Git repository is required.
//
// Returns:
// - A pointer to a GitRepository struct if a Git repository is found.
// - An error if there is an issue with the search or if the repository is required but not found.
func LocateGitRepository(startPath string, required bool) (*GitRepository, error) {
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return nil, err
	}

	gitPath := filepath.Join(absPath, GitExtension)
	if pathExists(gitPath) {
		return initializeGitRepo(absPath, false)
	}

	parentPath := filepath.Dir(absPath)
	if parentPath == absPath {
		if required {
			return nil, fmt.Errorf("no git directory found")
		}
		return nil, nil
	}

	return LocateGitRepository(parentPath, required)
}

func LocateCurrentRepository() (*GitRepository, error) {
	return LocateGitRepository(".", true)
}
