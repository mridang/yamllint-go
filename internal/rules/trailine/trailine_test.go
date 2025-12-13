package trailinglines

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestTrailingLinesRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name:   "too many blank lines at end",
			yaml:   "---\nkey: value\n\n\n",
			config: Config{},
			expected: []types.Problem{
				{Line: 2, Column: 0, Desc: "too many blank lines at end of file (1 > 0)"},
			},
		},
		{
			name:     "valid single trailing newline",
			yaml:     "---\nkey: value\n",
			config:   Config{},
			expected: []types.Problem{},
		},
		{
			name:     "no trailing newline",
			yaml:     "---\nkey: value",
			config:   Config{},
			expected: []types.Problem{},
		},
		{
			name:   "multiple trailing blank lines",
			yaml:   "---\nkey: value\n\n\n\n",
			config: Config{},
			expected: []types.Problem{
				{Line: 2, Column: 0, Desc: "too many blank lines at end of file (1 > 0)"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config

			require.NotNil(t, config)
			require.Equal(t, "trailing-lines", rule.ID())
			require.Equal(t, "line", rule.Type())
		})
	}
}
