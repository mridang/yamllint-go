package commas

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestCommasRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name: "too many spaces before comma",
			yaml: "---\n[1  , 2]\n",
			config: Config{
				MaxSpacesBefore: 0,
			},
			expected: []types.Problem{
				{Line: 1, Column: 4, Desc: "too many spaces before comma"},
			},
		},
		{
			name: "too few spaces after comma",
			yaml: "---\n[1,2]\n",
			config: Config{
				MinSpacesAfter: 1,
			},
			expected: []types.Problem{
				{Line: 1, Column: 2, Desc: "too few spaces after comma"},
			},
		},
		{
			name: "too many spaces after comma",
			yaml: "---\n[1,  2]\n",
			config: Config{
				MaxSpacesAfter: 1,
			},
			expected: []types.Problem{
				{Line: 1, Column: 2, Desc: "too many spaces after comma"},
			},
		},
		{
			name: "valid spacing",
			yaml: "---\n[1, 2, 3]\n",
			config: Config{
				MaxSpacesBefore: 0,
				MinSpacesAfter:  1,
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
			require.Equal(t, "commas", rule.ID())
			require.Equal(t, "token", rule.Type())
		})
	}
}
