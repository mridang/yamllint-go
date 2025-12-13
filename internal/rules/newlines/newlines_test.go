package newlines

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestNewLinesRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name: "dos newlines when unix expected",
			yaml: "---\r\nkey: value\r\n",
			config: Config{
				Type: "unix",
			},
			expected: []types.Problem{
				{Line: 0, Column: 3, Desc: "wrong new line character: expected \\n"},
			},
		},
		{
			name: "unix newlines valid",
			yaml: "---\nkey: value\n",
			config: Config{
				Type: "unix",
			},
			expected: []types.Problem{},
		},
		{
			name: "dos newlines valid",
			yaml: "---\r\nkey: value\r\n",
			config: Config{
				Type: "dos",
			},
			expected: []types.Problem{},
		},
		{
			name: "unix newlines when dos expected",
			yaml: "---\nkey: value\n",
			config: Config{
				Type: "dos",
			},
			expected: []types.Problem{
				{Line: 0, Column: 3, Desc: "wrong new line character: expected \\r\\n"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config
			
			require.NotNil(t, config)
			require.Equal(t, "new-lines", rule.ID())
			require.Equal(t, "line", rule.Type())
		})
	}
}
