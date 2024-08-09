package test

import (
	"testing"
	"time"

	"github.com/utkarsh5026/justdoit/app/cmd/objects"
	"github.com/utkarsh5026/justdoit/app/ordereddict"
)

func TestGitTag(t *testing.T) {
	t.Run("AnnotationTag", func(t *testing.T) {
		name := "v1.0.0"
		sha := "abc123"
		tagger := "John Doe <john@example.com>"
		message := "Release version 1.0.0"

		tag := objects.AnnotationTag(name, sha, tagger, message)

		if tag.Name != name {
			t.Errorf("Expected Name to be %s, got %s", name, tag.Name)
		}
		if tag.Object != sha {
			t.Errorf("Expected Object to be %s, got %s", sha, tag.Object)
		}
		if tag.Type != objects.CommitType.String() {
			t.Errorf("Expected Type to be %s, got %s", objects.CommitType.String(), tag.Type)
		}
		if tag.Tagger != tagger {
			t.Errorf("Expected Tagger to be %s, got %s", tagger, tag.Tagger)
		}
		if tag.Timestamp.IsZero() {
			t.Error("Expected Timestamp to be non-zero")
		}
		if tag.Message != message {
			t.Errorf("Expected Message to be %s, got %s", message, tag.Message)
		}
	})

	t.Run("IsAnnotation", func(t *testing.T) {
		annotatedTag := &objects.GitTag{Type: objects.CommitType.String()}
		lightweightTag := &objects.GitTag{}

		if !annotatedTag.IsAnnotation() {
			t.Error("Expected annotatedTag to be an annotation")
		}
		if lightweightTag.IsAnnotation() {
			t.Error("Expected lightweightTag not to be an annotation")
		}
	})

	t.Run("ToKvlm", func(t *testing.T) {
		tag := &objects.GitTag{
			Name:      "v1.0.0",
			Object:    "abc123",
			Type:      objects.CommitType.String(),
			Tagger:    "John Doe <john@example.com>",
			Timestamp: time.Unix(1623456789, 0),
			Message:   "Release version 1.0.0",
		}

		kvlm := tag.ToKvlm()

		expectedFields := map[string]string{
			"object": "abc123",
			"type":   objects.CommitType.String(),
			"tag":    "v1.0.0",
			"tagger": "John Doe <john@example.com> 1623456789 +0000",
			"":       "Release version 1.0.0",
		}

		for key, expectedValue := range expectedFields {
			if value, exists := kvlm.Get(key); !exists {
				t.Errorf("Expected key %s to exist in KVLM", key)
			} else if string(value.([]byte)) != expectedValue {
				t.Errorf("For key %s, expected %s, got %s", key, expectedValue, string(value.([]byte)))
			}
		}
	})

	t.Run("FromKVLM", func(t *testing.T) {
		t.Run("AnnotatedTag", func(t *testing.T) {
			kvlm := ordereddict.New()
			kvlm.Set("object", []byte("abc123"))
			kvlm.Set("type", []byte(objects.CommitType.String()))
			kvlm.Set("tag", []byte("v1.0.0"))
			kvlm.Set("tagger", []byte("John Doe <john@example.com> 1623456789 +0000"))
			kvlm.Set("", []byte("Release version 1.0.0"))

			tag := &objects.GitTag{}
			err := tag.FromKVLM(kvlm)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			expectedFields := map[string]string{
				"Name":    "v1.0.0",
				"Object":  "abc123",
				"Type":    objects.CommitType.String(),
				"Tagger":  "John Doe <john@example.com> 1623456789 +0000",
				"Message": "Release version 1.0.0",
			}

			for field, expected := range expectedFields {
				switch field {
				case "Name":
					if tag.Name != expected {
						t.Errorf("Expected Name to be %s, got %s", expected, tag.Name)
					}
				case "Object":
					if tag.Object != expected {
						t.Errorf("Expected Object to be %s, got %s", expected, tag.Object)
					}
				case "Type":
					if tag.Type != expected {
						t.Errorf("Expected Type to be %s, got %s", expected, tag.Type)
					}
				case "Tagger":
					if tag.Tagger != expected {
						t.Errorf("Expected Tagger to be %s, got %s", expected, tag.Tagger)
					}
				case "Message":
					if tag.Message != expected {
						t.Errorf("Expected Message to be %s, got %s", expected, tag.Message)
					}
				}
			}

			expectedTime := time.Unix(1623456789, 0)
			if !tag.Timestamp.Equal(expectedTime) {
				t.Errorf("Expected Timestamp to be %v, got %v", expectedTime, tag.Timestamp)
			}
		})

		t.Run("LightweightTag", func(t *testing.T) {
			kvlm := ordereddict.New()
			kvlm.Set("object", []byte("abc123"))

			tag := &objects.GitTag{}
			err := tag.FromKVLM(kvlm)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tag.Object != "abc123" {
				t.Errorf("Expected Object to be abc123, got %s", tag.Object)
			}
			if tag.Type != "" {
				t.Errorf("Expected Type to be empty, got %s", tag.Type)
			}
			if tag.Tagger != "" {
				t.Errorf("Expected Tagger to be empty, got %s", tag.Tagger)
			}
			if !tag.Timestamp.IsZero() {
				t.Errorf("Expected Timestamp to be zero, got %v", tag.Timestamp)
			}
			if tag.Message != "" {
				t.Errorf("Expected Message to be empty, got %s", tag.Message)
			}
		})
	})
}
