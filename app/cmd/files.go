package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// isDir checks if the given path is a directory.
//
// Parameters:
// - path: The path to check as a string.
//
// Returns:
// - A boolean indicating whether the path is a directory.
// - An error if there is an issue retrieving the file information.
func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// createRepoPath constructs a file path by joining the repository path with additional paths.
//
// Parameters:
// - repo: A pointer to a GitRepository struct containing the repository paths.
// - paths: A variadic parameter representing additional path segments to be joined.
//
// Returns:
// - A string representing the combined file path.
func createRepoPath(repo *GitRepository, paths ...string) string {
	paths = append([]string{repo.GitDir}, paths...)
	return filepath.Join(paths...)
}

// repoDir constructs a directory path within a repository and optionally creates the directory.
//
// Parameters:
// - repo: A pointer to a GitRepository struct containing the repository paths.
// - mkdir: A boolean indicating whether to create the directory if it does not exist.
// - paths: A variadic parameter representing additional path segments to be joined.
//
// Returns:
// - A string representing the combined directory path. If an error occurs, an empty string is returned.
// - An error if there is an issue checking the directory status or creating the directory
func repoDir(repo *GitRepository, mkdir bool, paths ...string) (string, error) {
	path := createRepoPath(repo, paths...)
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

// repoFile constructs a file path within a repository and optionally creates the necessary directories.
//
// Parameters:
// - repo: A pointer to a GitRepository struct containing the repository paths.
// - mkdir: A boolean indicating whether to create the directory if it does not exist.
// - paths: A variadic parameter representing additional path segments to be joined.
//
// Returns:
// - A string representing the combined file path. If an error occurs, an empty string is returned
func repoFile(repo *GitRepository, mkdir bool, paths ...string) string {
	dirPath := paths[:len(paths)-1]
	_, err := repoDir(repo, mkdir, dirPath...)
	if err != nil {
		return ""
	}
	return createRepoPath(repo, paths...)
}

// listDir lists the contents of a directory.
//
// Parameters:
// - path: The path to the directory to be listed.
//
// Returns:
// - A slice of os.DirEntry representing the contents of the directory.
// - An error if there is an issue opening or reading the directory.
func listDir(path string) ([]os.DirEntry, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {
			panic(err)
		}
	}(dir)

	files, err := dir.ReadDir(-1) // Read all files
	if err != nil {
		return nil, err
	}

	var dirContents []os.DirEntry
	for _, file := range files {
		dirContents = append(dirContents, file)
	}

	return dirContents, nil
}

func findRepoDir(startPath string, required bool) (*GitRepository, error) {
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

	return findRepoDir(parentPath, required)
}
