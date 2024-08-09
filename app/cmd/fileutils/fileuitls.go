package fileutils

import "os"

// IsDir checks if the given path is a directory.
//
// Parameters:
// - path: The path to check as a string.
//
// Returns:
// - A boolean indicating whether the path is a directory.
// - An error if there is an issue retrieving the file information.
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// PathExists checks if the given path exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsFile checks if the given path is a file.
func IsFile(path string) (bool, error) {
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

// ListDir lists the contents of a directory.
//
// Parameters:
// - path: The path to the directory to be listed.
//
// Returns:
// - A slice of os.DirEntry representing the contents of the directory.
// - An error if there is an issue opening or reading the directory.
func ListDir(path string) ([]os.DirEntry, error) {
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
