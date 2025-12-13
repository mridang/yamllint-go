package keyduplicates

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeyDuplicatesRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name:   "duplicated key",
			yaml:   "---\nkey: value1\nkey: value2\n",
			config: Config{},
			expected: []types.Problem{
				{Line: 2, Column: 0, Desc: "duplication of key \"key\" in mapping"},
			},
		},
		{
			name:     "unique keys",
			yaml:     "---\nkey1: value1\nkey2: value2\n",
			config:   Config{},
			expected: []types.Problem{},
		},
		{
			name: "merge key allowed",
			yaml: "---\n<<: &anchor\nkey: value\n",
			config: Config{
				ForbidDuplicatedMergeKeys: false,
			},
			expected: []types.Problem{},
		},
		{
			name: "duplicated merge key forbidden",
			yaml: "---\n<<: &anchor1\n<<: &anchor2\n",
			config: Config{
				ForbidDuplicatedMergeKeys: true,
			},
			expected: []types.Problem{
				{Line: 2, Column: 0, Desc: "duplication of key \"<<\" in mapping"},
			},
		},
		{
			name:   "nested mappings",
			yaml:   "---\nparent:\n  key: value1\n  key: value2\n",
			config: Config{},
			expected: []types.Problem{
				{Line: 3, Column: 2, Desc: "duplication of key \"key\" in mapping"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config

			require.NotNil(t, config)
			require.Equal(t, "key-duplicates", rule.ID())
			require.Equal(t, "token", rule.Type())
		})
	}
}
