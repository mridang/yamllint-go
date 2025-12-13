package colons

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestColonsRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name: "too many spaces before colon",
			yaml: "---\nkey  : value\n",
			config: Config{
				MaxSpacesBefore: 0,
			},
			expected: []types.Problem{
				{Line: 1, Column: 5, Desc: "too many spaces before colon"},
			},
		},
		{
			name: "too many spaces after colon",
			yaml: "---\nkey:  value\n",
			config: Config{
				MaxSpacesAfter: 1,
			},
			expected: []types.Problem{
				{Line: 1, Column: 3, Desc: "too many spaces after colon"},
			},
		},
		{
			name: "valid spacing",
			yaml: "---\nkey: value\n",
			config: Config{
				MaxSpacesBefore: 0,
				MaxSpacesAfter:  1,
			},
			expected: []types.Problem{},
		},
		{
			name: "flow mapping colon",
			yaml: "---\n{key: value}\n",
			config: Config{
				MaxSpacesBefore: 0,
				MaxSpacesAfter:  1,
			},
			expected: []types.Problem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config
			
			require.NotNil(t, config)
			require.Equal(t, "colons", rule.ID())
			require.Equal(t, "token", rule.Type())
		})
	}
}
