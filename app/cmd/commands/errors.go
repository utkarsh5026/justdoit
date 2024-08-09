package commands

import (
	"fmt"
	"github.com/utkarsh5026/justdoit/app/cmd/objects"
)

var RepoNotFound = func(err error) error {
	return fmt.Errorf("unable to locate repository: %w", err)
}

var ObjectMismatchError = func(expected, actual objects.GitObjectType) error {
	return fmt.Errorf("object type mismatch: expected %s, got %s", expected, actual)
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
