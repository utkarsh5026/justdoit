package repository

import "os"

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

// pathExists checks if the given path exists.
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// isFile checks if the given path is a file.
func isFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return !fileInfo.IsDir(), nil
}

// CloseFile  closes the file and panics if an error occurs.
func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
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

	return dir.ReadDir(-1) // Read all files
}
