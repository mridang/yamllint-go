package docstart

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestDocumentStartRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name: "missing document start",
			yaml: "key: value\n",
			config: Config{
				Present: true,
			},
			expected: []types.Problem{
				{Line: 0, Column: 0, Desc: "missing document start \"---\""},
			},
		},
		{
			name: "valid document start",
			yaml: "---\nkey: value\n",
			config: Config{
				Present: true,
			},
			expected: []types.Problem{},
		},
		{
			name: "forbidden document start",
			yaml: "---\nkey: value\n",
			config: Config{
				Present: false,
			},
			expected: []types.Problem{
				{Line: 0, Column: 0, Desc: "found forbidden document start \"---\""},
			},
		},
		{
			name: "no document start when not required",
			yaml: "key: value\n",
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
			require.Equal(t, "document-start", rule.ID())
			require.Equal(t, "token", rule.Type())
		})
	}
}
