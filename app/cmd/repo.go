package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	GitExtension = ".git"
	HeadFile     = "HEAD"
	DescFile     = "description"
	ConfigFile   = "config"
)

type GitRepository struct {
	WorkTree string       // The path to the repository.
	GitDir   string       // The path to the .git directory.
	Config   *viper.Viper // The configuration file.
}

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

	repo.Config.SetConfigName("config")
	repo.Config.AddConfigPath(repo.GitDir)
	repo.Config.SetConfigType("ini")

	if err := readConfig(&repo, force); err != nil {
		return nil, err
	}
	return &repo, nil
}

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
	config.SetConfigFile(repoFile(repo, false, ConfigFile))

	if err := config.WriteConfig(); err != nil {
		return nil, err
	}
	return repo, nil
}

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
// - repo: A pointer to a GitRepository struct containing the repository paths.
//
// Returns:
// - An error if any of the directory creation operations fail.
func createInitialDirectories(repo *GitRepository) error {
	if _, err := repoDir(repo, true, "branches"); err != nil {
		return err
	}

	if _, err := repoDir(repo, true, "objects"); err != nil {
		return err
	}

	if _, err := repoDir(repo, true, "refs", "tags"); err != nil {
		return err
	}

	if _, err := repoDir(repo, true, "refs", "heads"); err != nil {
		return err
	}

	return nil
}

// createGitFiles creates the initial files for a Git repository.
//
// Parameters:
// - repo: A pointer to a GitRepository struct containing the repository paths.
//
// Returns:
// - An error if any of the file creation operations fail.
func createGitFiles(repo *GitRepository) error {
	// .git/description
	descriptionPath := repoFile(repo, false, DescFile)
	descriptionContent := "Unnamed repository; edit this file 'description' to name the repository.\n"
	if err := os.WriteFile(descriptionPath, []byte(descriptionContent), 0644); err != nil {
		return err
	}

	// .git/HEAD
	headPath := repoFile(repo, false, HeadFile)
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

	config.Set("core.repositoryformatversion", "0")
	config.Set("core.filemode", "false")
	config.Set("core.bare", "false")

	return config
}
