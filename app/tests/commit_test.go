package test

import (
	"github.com/utkarsh5026/justdoit/app/cmd/objects"
	"github.com/utkarsh5026/justdoit/app/ordereddict"
	"reflect"
	"testing"
	"time"
)

func TestParseSignature(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    objects.GitSignature
		wantErr bool
	}{
		{
			name:  "Valid signature",
			input: []byte("John Doe <john@example.com> 1623456789 +0100"),
			want: objects.GitSignature{
				Name:  "John Doe",
				Email: "john@example.com",
				When:  time.Unix(1623456789, 0),
			},
			wantErr: false,
		},
		{
			name:    "Invalid signature format",
			input:   []byte("Invalid signature"),
			want:    objects.GitSignature{},
			wantErr: true,
		},
		{
			name:    "Invalid timestamp",
			input:   []byte("John Doe <john@example.com> invalid +0100"),
			want:    objects.GitSignature{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := objects.ParseSignature(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateCommitFromKVLM(t *testing.T) {
	tests := []struct {
		name    string
		kvlm    *ordereddict.OrderedDict
		want    *objects.GitCommit
		wantErr bool
	}{
		{
			name: "Valid commit",
			kvlm: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("tree", []byte("29ff16c9c14e2652b22f8b78bb08a5a07930c147"))
				od.Set("parent", []byte("206941306e8a8af65b66eaaaea388a7ae24d49a0"))
				od.Set("author", []byte("John Doe <john@example.com> 1623456789 +0100"))
				od.Set("committer", []byte("Jane Smith <jane@example.com> 1623456790 +0100"))
				od.Set("", []byte("Implement new feature"))
				return od
			}(),
			want: &objects.GitCommit{
				Tree:    "29ff16c9c14e2652b22f8b78bb08a5a07930c147",
				Parents: []string{"206941306e8a8af65b66eaaaea388a7ae24d49a0"},
				Author: objects.GitSignature{
					Name:  "John Doe",
					Email: "john@example.com",
					When:  time.Unix(1623456789, 0),
				},
				Committer: objects.GitSignature{
					Name:  "Jane Smith",
					Email: "jane@example.com",
					When:  time.Unix(1623456790, 0),
				},
				Message: "Implement new feature",
			},
			wantErr: false,
		},
		{
			name: "Missing tree",
			kvlm: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("author", []byte("John Doe <john@example.com> 1623456789 +0100"))
				od.Set("committer", []byte("Jane Smith <jane@example.com> 1623456790 +0100"))
				od.Set("", []byte("Implement new feature"))
				return od
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "Missing author",
			kvlm: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("tree", []byte("29ff16c9c14e2652b22f8b78bb08a5a07930c147"))
				od.Set("committer", []byte("Jane Smith <jane@example.com> 1623456790 +0100"))
				od.Set("", []byte("Implement new feature"))
				return od
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "Missing committer",
			kvlm: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("tree", []byte("29ff16c9c14e2652b22f8b78bb08a5a07930c147"))
				od.Set("author", []byte("John Doe <john@example.com> 1623456789 +0100"))
				od.Set("", []byte("Implement new feature"))
				return od
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "Multiple parents",
			kvlm: func() *ordereddict.OrderedDict {
				od := ordereddict.New()
				od.Set("tree", []byte("29ff16c9c14e2652b22f8b78bb08a5a07930c147"))
				od.Set("parent", [][]byte{
					[]byte("206941306e8a8af65b66eaaaea388a7ae24d49a0"),
					[]byte("f7e1cf3b22eb0df3270e58e333e138c97a596ca3"),
				})
				od.Set("author", []byte("John Doe <john@example.com> 1623456789 +0100"))
				od.Set("committer", []byte("Jane Smith <jane@example.com> 1623456790 +0100"))
				od.Set("", []byte("Merge branch 'feature'"))
				return od
			}(),
			want: &objects.GitCommit{
				Tree: "29ff16c9c14e2652b22f8b78bb08a5a07930c147",
				Parents: []string{
					"206941306e8a8af65b66eaaaea388a7ae24d49a0",
					"f7e1cf3b22eb0df3270e58e333e138c97a596ca3",
				},
				Author: objects.GitSignature{
					Name:  "John Doe",
					Email: "john@example.com",
					When:  time.Unix(1623456789, 0),
				},
				Committer: objects.GitSignature{
					Name:  "Jane Smith",
					Email: "jane@example.com",
					When:  time.Unix(1623456790, 0),
				},
				Message: "Merge branch 'feature'",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := objects.CreateCommitFromKVLM(tt.kvlm)
			if (err != nil) != tt.wantErr {
				t.Errorf("createCommitFromKVLM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createCommitFromKVLM() = %v, want %v", got, tt.want)
			}
		})
	}
}
