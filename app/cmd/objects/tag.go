package objects

import (
	"fmt"
	"github.com/utkarsh5026/justdoit/app/cmd/repository"
	"github.com/utkarsh5026/justdoit/app/ordereddict"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GitTag represents a Git tag object, which can be either a lightweight tag or an annotated tag.
//
// A lightweight tag is simply a reference to a commit SHA, while an annotated tag contains additional metadata
// such as the tagger's name, a timestamp, and a message.
type GitTag struct {
	// Common fields for both lightweight and annotated tags
	Name   string // The name of the tag
	Object string // SHA-1 hash of the referenced commit

	// Fields specific to annotated tags
	Type      string    // The type of the object referenced by the tag (e.g., "commit")
	Tagger    string    // The name of the person who created the tag
	Timestamp time.Time // The timestamp of the tag
	Message   string    // The message associated with the tag
}

// AnnotationTag creates a new annotated Git tag with the given name, SHA, tagger, and message.
func AnnotationTag(name string, sha string, tagger string, message string) *GitTag {
	return &GitTag{
		Name:      name,
		Object:    sha,
		Type:      CommitType.String(),
		Tagger:    tagger,
		Timestamp: time.Now(),
		Message:   message,
	}
}

// IsAnnotation returns true if the Git tag is an annotated tag (commit type), false otherwise.
func (gt *GitTag) IsAnnotation() bool {
	return gt.Type == CommitType.String()
}

// ToKvlm converts the GitTag object to an OrderedDict representation.
//
// This method creates a new OrderedDict and populates it with the fields of the GitTag object
// if it is an annotated tag. The fields include the object SHA, type, tag name, tagger information,
// and the tag message.
//
// Returns:
// - An *ordereddict.OrderedDict containing the key-value pairs representing the GitTag object.
func (gt *GitTag) ToKvlm() *ordereddict.OrderedDict {
	kvlm := ordereddict.New()
	if gt.IsAnnotation() {
		kvlm.Set("object", []byte(gt.Object))
		kvlm.Set("type", []byte(gt.Type))
		kvlm.Set("tag", []byte(gt.Name))
		kvlm.Set("tagger", []byte(fmt.Sprintf("%s %d +0000", gt.Tagger, gt.Timestamp.Unix())))
		kvlm.Set("", []byte(gt.Message))
	}
	return kvlm
}

// FromKVLM populates the GitTag object from an OrderedDict representation.
//
// This method extracts the fields from the given OrderedDict and assigns them to the corresponding
// fields of the GitTag object. It handles both common fields and annotated tag specific fields.
//
// Parameters:
// - kvlm: An *ordereddict.OrderedDict containing the key-value pairs representing the GitTag object.
//
// Returns:
// - An error if any operation fails, otherwise nil.
func (gt *GitTag) FromKVLM(kvlm *ordereddict.OrderedDict) error {
	if obj, exists := kvlm.Get("object"); exists {
		gt.Object = string(obj.([]byte))
	}

	if tagName, exists := kvlm.Get("tag"); exists {
		gt.Name = string(tagName.([]byte))
	}

	// Parse annotated tag specific fields
	if typ, exists := kvlm.Get("type"); exists {
		gt.Type = string(typ.([]byte))
	}

	if tagger, exists := kvlm.Get("tagger"); exists {
		taggerInfo := string(tagger.([]byte))
		gt.Tagger = taggerInfo

		// Extract timestamp from tagger info
		parts := strings.Fields(taggerInfo)
		if len(parts) >= 2 {
			timestamp, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
			if err == nil {
				gt.Timestamp = time.Unix(timestamp, 0)
			}
		}
	}

	if message, exists := kvlm.Get(""); exists {
		gt.Message = string(message.([]byte))
	}

	return nil
}

type TagObject struct {
	kvlm    *ordereddict.OrderedDict
	tagInfo *GitTag
}

func (to *TagObject) Serialize() ([]byte, error) {
	return KvlmSerialize(to.kvlm)
}

func (to *TagObject) Deserialize(data []byte) error {
	var err error
	to.kvlm = KvlmParse(data, 0, nil)
	return err
}

func (to *TagObject) Format() GitObjectType {
	return TagType
}

func (to *TagObject) SetData(data []byte) {
	to.kvlm = KvlmParse(data, 0, nil)
}

func Tag() *TagObject {
	return &TagObject{
		kvlm: ordereddict.New(),
	}
}

// CreateTag creates a new tag in the Git repository.
//
// This function creates either a lightweight tag or an annotated tag based on the createTagObject parameter.
// It finds the object SHA from the reference, creates the tag object if needed, and writes the tag to the repository.
//
// Parameters:
// - repo: The Git repository object (*repository.GitRepository).
// - name: The name of the tag to create.
// - ref: The reference from which to create the tag.
// - createTagObject: A boolean indicating whether to create an annotated tag object.
// - tagger: The name of the person who created the tag.
// - message: The message associated with the tag.
//
// Returns:
// - An error if any operation fails, otherwise nil.
func CreateTag(repo *repository.GitRepository, name string, ref string, createTagObject bool, tagger string, message string) error {
	om := NewObjectManager(repo)
	sha := om.FindObject(ref, TagType, true)

	if createTagObject {
		tag := Tag()
		tag.tagInfo = AnnotationTag(name, sha, tagger, message)
		tag.kvlm = tag.tagInfo.ToKvlm()
		tagSha, err := om.WriteObject(tag, false)

		if err != nil {
			return err
		}

		return createRef(repo, filepath.Join("tags", name), tagSha)
	} else {
		refName := filepath.Join("tags", name)
		return createRef(repo, refName, sha)
	}
}

// createRef creates a new reference in the Git repository.
//
// This function writes the given SHA to a reference file in the repository.
// The reference file is created in the "refs" directory.
//
// Parameters:
// - repo: The Git repository object *repository.GitRepository.
// - refName: The name of the reference to create.
// - sha: The SHA value to write to the reference file.
//
// Returns:
// - An error if any operation fails, otherwise nil.
func createRef(repo *repository.GitRepository, refName string, sha string) error {
	refPath := filepath.Join("refs", refName)
	path := repository.GetGitFilePath(repo, false, refPath)

	if err := os.WriteFile(path, []byte(sha+"\n"), 0644); err != nil {
		return err
	}
	return nil
}
