package commands

import "fmt"

var RepoNotFound = func(err error) error {
	return fmt.Errorf("unable to locate repository: %w", err)
}

var InvalidTreeObject = func(err error) error {
	return fmt.Errorf("invalid tree object: %w", err)
}

var ObjectReadError = func(err error) error {
	return fmt.Errorf("failed to read object: %w", err)
}

var ObjectTypeError = func(err error) error {
	return fmt.Errorf("invalid object type: %w", err)
}

var SerializationError = func(err error) error {
	return fmt.Errorf("failed to serialize object: %w", err)
}
