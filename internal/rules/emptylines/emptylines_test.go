package emptylines

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestEmptyLinesRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name: "too many blank lines",
			yaml: "---\nkey: value\n\n\n\nother: value\n",
			config: Config{
				Max: 2,
			},
			expected: []types.Problem{
				{Line: 4, Column: 0, Desc: "too many blank lines (3 > 2)"},
			},
		},
		{
			name: "valid blank lines",
			yaml: "---\nkey: value\n\n\nother: value\n",
			config: Config{
				Max: 2,
			},
			expected: []types.Problem{},
		},
		{
			name: "too many blank lines at start",
			yaml: "\n\n---\nkey: value\n",
			config: Config{
				MaxStart: 0,
			},
			expected: []types.Problem{
				{Line: 0, Column: 0, Desc: "too many blank lines (2 > 0)"},
			},
		},
		{
			name: "too many blank lines at end",
			yaml: "---\nkey: value\n\n\n",
			config: Config{
				MaxEnd: 0,
			},
			expected: []types.Problem{
				{Line: 3, Column: 0, Desc: "too many blank lines (2 > 0)"},
			},
		},
		{
			name: "valid document structure",
			yaml: "---\nkey: value\n\nother: value\n",
			config: Config{
				Max:      2,
				MaxStart: 0,
				MaxEnd:   0,
			},
			expected: []types.Problem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config
			
			require.NotNil(t, config)
			require.Equal(t, "empty-lines", rule.ID())
			require.Equal(t, "line", rule.Type())
		})
	}
}
