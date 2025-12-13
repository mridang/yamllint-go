package commentsindentation

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestCommentsIndentationRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name:   "comment not indented like content",
			yaml:   "---\nkey: value\n  # wrong indent\n",
			config: Config{},
			expected: []types.Problem{
				{Line: 2, Column: 2, Desc: "comment not indented like content"},
			},
		},
		{
			name:     "valid comment indentation",
			yaml:     "---\nkey: value\n# comment\n",
			config:   Config{},
			expected: []types.Problem{},
		},
		{
			name:     "nested comment indentation",
			yaml:     "---\nparent:\n  child: value\n  # nested comment\n",
			config:   Config{},
			expected: []types.Problem{},
		},
		{
			name:     "comment before content",
			yaml:     "---\n# comment\nkey: value\n",
			config:   Config{},
			expected: []types.Problem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config

			require.NotNil(t, config)
			require.Equal(t, "comments-indentation", rule.ID())
			require.Equal(t, "comment", rule.Type())
		})
	}
}
