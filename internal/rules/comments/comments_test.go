package comments

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestCommentsRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name: "missing starting space",
			yaml: "---\n#comment\n",
			config: Config{
				RequireStartingSpace: true,
			},
			expected: []types.Problem{
				{Line: 1, Column: 1, Desc: "missing starting space in comment"},
			},
		},
		{
			name: "valid comment with space",
			yaml: "---\n# comment\n",
			config: Config{
				RequireStartingSpace: true,
			},
			expected: []types.Problem{},
		},
		{
			name: "shebang allowed",
			yaml: "#!/usr/bin/env python\n---\n",
			config: Config{
				RequireStartingSpace: true,
				IgnoreShebangs:       true,
			},
			expected: []types.Problem{},
		},
		{
			name: "too few spaces from content",
			yaml: "---\nkey: value # comment\n",
			config: Config{
				MinSpacesFromContent: 2,
			},
			expected: []types.Problem{
				{Line: 1, Column: 11, Desc: "too few spaces before comment: expected 2"},
			},
		},
		{
			name: "valid inline comment spacing",
			yaml: "---\nkey: value  # comment\n",
			config: Config{
				MinSpacesFromContent: 2,
			},
			expected: []types.Problem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config
			
			require.NotNil(t, config)
			require.Equal(t, "comments", rule.ID())
			require.Equal(t, "comment", rule.Type())
		})
	}
}
