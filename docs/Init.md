# Git-like Implementation in Go

## Introduction

Welcome to our Git-like implementation project in Go! This project aims to recreate some of the basic functionality of Git, helping you understand how Git works under the hood. It's a great way to learn both about Git internals and Go programming.

## What We're Building

We're creating a simplified version of Git that can:
1. Initialize a new repository
2. Set up the basic Git directory structure
3. Create and manage configuration files

## Key Components

### 1. GitRepository Struct

```go
type GitRepository struct {
    WorkTree string
    GitDir   string
    Config   *viper.Viper
}
```

This struct represents our Git repository. It keeps track of:
- The working directory (`WorkTree`)
- The `.git` directory (`GitDir`)
- The repository configuration (`Config`)

### 2. Main Functions

#### Initializing a Repository

```go
func initializeGitRepo(path string, force bool) (*GitRepository, error)
```

This function sets up a new Git repository. It:
- Creates a `GitRepository` struct
- Checks if the directory is already a Git repo (unless `force` is true)
- Sets up the configuration

#### Creating a Repository

```go
func CreateGitRepository(path string) (*GitRepository, error)
```

This function fully creates a new Git repository. It:
- Initializes the repo
- Ensures the repository directory is valid
- Creates the initial directory structure
- Creates necessary Git files
- Sets up default configuration

### 3. Helper Functions

We have several helper functions that handle specific tasks:

- `readConfig`: Reads the repository configuration
- `ensureValidRepoExists`: Checks if the repository directory is valid
- `createInitialDirectories`: Sets up the Git directory structure
- `createGitFiles`: Creates initial Git files (like HEAD)
- `repoDefaultConfig`: Sets up default Git configuration

## How It Works

1. When you create a new repository, `CreateGitRepository` is called.
2. It first initializes the basic structure using `initializeGitRepo`.
3. Then it ensures the repository location is valid with `ensureValidRepoExists`.
4. Next, it creates the directory structure (`createInitialDirectories`) and necessary files (`createGitFiles`).
5. Finally, it sets up the default configuration.

## Learning Points

- **Go Structs**: We use structs to organize our data.
- **Error Handling**: Go's error handling is used throughout the code.
- **File Operations**: We perform various file and directory operations.
- **Configuration Management**: We use the `viper` library to manage Git configs.

## Next Steps

This is just the beginning! To make this more Git-like, we could add features like:
- Committing changes
- Creating branches
- Merging branches
- Handling remote repositories

Feel free to explore the code and experiment with adding new features!