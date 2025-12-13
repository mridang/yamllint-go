package keyordering

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeyOrderingRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name:   "wrong key ordering",
			yaml:   "---\nzebra: 1\napple: 2\n",
			config: Config{},
			expected: []types.Problem{
				{Line: 2, Column: 0, Desc: "wrong ordering of key \"apple\" in mapping"},
			},
		},
		{
			name:     "correct key ordering",
			yaml:     "---\napple: 1\nbanana: 2\nzebra: 3\n",
			config:   Config{},
			expected: []types.Problem{},
		},
		{
			name: "ignored keys",
			yaml: "---\nzebra: 1\n_internal: 2\napple: 3\n",
			config: Config{
				IgnoredKeys: []string{"^_.*"},
			},
			expected: []types.Problem{},
		},
		{
			name:   "nested mappings ordering",
			yaml:   "---\nparent:\n  zebra: 1\n  apple: 2\n",
			config: Config{},
			expected: []types.Problem{
				{Line: 3, Column: 2, Desc: "wrong ordering of key \"apple\" in mapping"},
			},
		},
		{
			name:     "alphabetical ordering valid",
			yaml:     "---\na: 1\nb: 2\nc: 3\n",
			config:   Config{},
			expected: []types.Problem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config

			require.NotNil(t, config)
			require.Equal(t, "key-ordering", rule.ID())
			require.Equal(t, "token", rule.Type())
		})
	}
}
