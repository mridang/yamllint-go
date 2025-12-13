package docend

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestDocumentEndRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name: "missing document end",
			yaml: "---\nkey: value\n",
			config: Config{
				Present: true,
			},
			expected: []types.Problem{
				{Line: 2, Column: 0, Desc: "missing document end \"...\""},
			},
		},
		{
			name: "valid document end",
			yaml: "---\nkey: value\n...\n",
			config: Config{
				Present: true,
			},
			expected: []types.Problem{},
		},
		{
			name: "forbidden document end",
			yaml: "---\nkey: value\n...\n",
			config: Config{
				Present: false,
			},
			expected: []types.Problem{
				{Line: 2, Column: 0, Desc: "found forbidden document end \"...\""},
			},
		},
		{
			name: "no document end required when false",
			yaml: "---\nkey: value\n",
			config: Config{
				Present: false,
			},
			expected: []types.Problem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config

			require.NotNil(t, config)
			require.Equal(t, "document-end", rule.ID())
			require.Equal(t, "token", rule.Type())
		})
	}
}
