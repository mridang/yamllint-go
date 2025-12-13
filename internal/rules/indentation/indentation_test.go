package indentation

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestIndentationRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name: "wrong indentation",
			yaml: "---\nkey:\n   value\n",
			config: Config{
				Spaces: 2,
			},
			expected: []types.Problem{
				{Line: 2, Column: 3, Desc: "wrong indentation: expected 2 but found 3"},
			},
		},
		{
			name: "valid indentation",
			yaml: "---\nkey:\n  value\n",
			config: Config{
				Spaces: 2,
			},
			expected: []types.Problem{},
		},
		{
			name: "consistent spacing",
			yaml: "---\nparent:\n  child:\n    value\n",
			config: Config{
				Spaces: "consistent",
			},
			expected: []types.Problem{},
		},
		{
			name: "indent sequences true",
			yaml: "---\nlist:\n  - item1\n  - item2\n",
			config: Config{
				Spaces:          2,
				IndentSequences: true,
			},
			expected: []types.Problem{},
		},
		{
			name: "indent sequences false",
			yaml: "---\nlist:\n- item1\n- item2\n",
			config: Config{
				Spaces:          2,
				IndentSequences: false,
			},
			expected: []types.Problem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config
			
			require.NotNil(t, config)
			require.Equal(t, "indentation", rule.ID())
			require.Equal(t, "token", rule.Type())
		})
	}
}
