package repository

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	GitExtension = ".justdoit"
	HeadFile     = "HEAD"
	DescFile     = "description"
	ConfigFile   = "config"
	ObjectDir    = "objects"
)

type GitRepository struct {
	WorkTree string       // The path to the repository.
	GitDir   string       // The path to the .git directory.
	Config   *viper.Viper // The configuration file.
}

// initializeGitRepo initializes a Git repository.
//
// Parameters:
// - path: The path to the repository.
// - force: A boolean indicating whether to force the initialization.
//
// Returns:
// - A pointer to a GitRepository struct containing the repository paths and configuration.
// - An error if any of the initialization operations fail.
func initializeGitRepo(path string, force bool) (*GitRepository, error) {
	repo := GitRepository{
		WorkTree: path,
		GitDir:   filepath.Join(path, GitExtension),
		Config:   viper.New(),
	}

	if !force {
		isDir, err := isDir(repo.GitDir)
		if err != nil {
			return nil, err
		}

		if !isDir {
			return nil, fmt.Errorf("'%s' is not a git repository", path)
		}
	}

	repo.Config.SetConfigName(ConfigFile)
	repo.Config.AddConfigPath(repo.GitDir)
	repo.Config.SetConfigType("ini")

	if err := readConfig(&repo, force); err != nil {
		return nil, err
	}
	return &repo, nil
}

// readConfig reads the configuration for the Git repository.
//
// Parameters:
// - repository: A pointer to a GitRepository struct containing the repository paths and configuration.
// - force: A boolean indicating whether to force the reading of the configuration.
//
// Returns:
// - An error if the configuration file cannot be read or if the repository format version is unsupported.
func readConfig(repo *GitRepository, force bool) error {
	if err := repo.Config.ReadInConfig(); err != nil {
		if !force {
			return fmt.Errorf("failed to read config file: %s", err)
		}
	} else {
		if !force {
			version := repo.Config.GetInt("core.repositoryformatversion")
			if version != 0 {
				return fmt.Errorf("unsupported repositoryformatversion %d", version)
			}
		}
	}
	return nil
}

// CreateGitRepository creates a new Git repository at the specified path.
//
// Parameters:
// - path: The path where the Git repository should be created.
//
// Returns:
// - A pointer to a GitRepository struct containing the repository paths and configuration.
// - An error if any of the repository creation operations fail.
func CreateGitRepository(path string) (*GitRepository, error) {
	repo, err := initializeGitRepo(path, true)
	if err != nil {
		return nil, err
	}

	if err := ensureValidRepoExists(repo); err != nil {
		return nil, err
	}

	if err := createInitialDirectories(repo); err != nil {
		return nil, err
	}

	if err := createGitFiles(repo); err != nil {
		return nil, err
	}

	config := repoDefaultConfig()
	config.SetConfigFile(GetGitFilePath(repo, false, ConfigFile))

	if err := config.WriteConfig(); err != nil {
		return nil, err
	}
	return repo, nil
}

// ensureValidRepoExists checks if the Git repository exists and is valid.
//
// Parameters:
// - repository: A pointer to a GitRepository struct containing the repository paths.
//
// Returns:
// - An error if the repository is not valid or if any of the directory operations fail.
func ensureValidRepoExists(repo *GitRepository) error {
	if pathExists(repo.GitDir) {
		isDir, err := isDir(repo.WorkTree)
		if err != nil || !isDir {
			return fmt.Errorf("'%s' is not a directory", repo.WorkTree)
		}

		dirs, err := listDir(repo.GitDir)
		if err != nil || len(dirs) > 0 {
			return fmt.Errorf("'%s' is not an empty directory", repo.GitDir)
		}
	} else {
		if err := os.MkdirAll(repo.WorkTree, 0755); err != nil {
			return err
		}
	}
	return nil
}

// createInitialDirectories creates the initial directory structure, like branches, objects, refs/tags, refs/heads
//
// Parameters:
// - repository: A pointer to a GitRepository struct containing the repository paths.
//
// Returns:
// - An error if any of the directory creation operations fail.
func createInitialDirectories(repo *GitRepository) error {
	if _, err := EnsureGitDirExists(repo, true, "branches"); err != nil {
		return err
	}

	if _, err := EnsureGitDirExists(repo, true, "objects"); err != nil {
		return err
	}

	if _, err := EnsureGitDirExists(repo, true, "refs", "tags"); err != nil {
		return err
	}

	if _, err := EnsureGitDirExists(repo, true, "refs", "heads"); err != nil {
		return err
	}

	return nil
}

// createGitFiles creates the initial files for a Git repository.
//
// Parameters:
// - repository: A pointer to a GitRepository struct containing the repository paths.
//
// Returns:
// - An error if any of the file creation operations fail.
func createGitFiles(repo *GitRepository) error {
	// .git/description
	descriptionPath := GetGitFilePath(repo, false, DescFile)
	descriptionContent := "Unnamed repository; edit this file 'description' to name the repository.\n"
	if err := os.WriteFile(descriptionPath, []byte(descriptionContent), 0644); err != nil {
		return err
	}

	// .git/HEAD
	headPath := GetGitFilePath(repo, false, HeadFile)
	headContent := "ref: refs/heads/master\n"
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		return err
	}

	return nil
}

// repoDefaultConfig creates and returns a default configuration for a Git repository.
//
// Returns:
// - A pointer to a viper.Viper instance containing the default configuration.
func repoDefaultConfig() *viper.Viper {
	config := viper.New()

	config.SetConfigType("ini")
	config.Set("core.repositoryformatversion", "0")
	config.Set("core.filemode", "false")
	config.Set("core.bare", "false")

	return config
}
