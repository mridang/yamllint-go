package octalvalues

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestOctalValuesRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name: "implicit octal forbidden",
			yaml: "---\nvalue: 0755\n",
			config: Config{
				ForbidImplicitOctal: true,
			},
			expected: []types.Problem{
				{Line: 1, Column: 7, Desc: "forbidden implicit octal value \"0755\""},
			},
		},
		{
			name: "explicit octal forbidden",
			yaml: "---\nvalue: 0o755\n",
			config: Config{
				ForbidExplicitOctal: true,
			},
			expected: []types.Problem{
				{Line: 1, Column: 7, Desc: "forbidden explicit octal value \"0o755\""},
			},
		},
		{
			name: "implicit octal allowed",
			yaml: "---\nvalue: 0755\n",
			config: Config{
				ForbidImplicitOctal: false,
			},
			expected: []types.Problem{},
		},
		{
			name: "explicit octal allowed",
			yaml: "---\nvalue: 0o755\n",
			config: Config{
				ForbidExplicitOctal: false,
			},
			expected: []types.Problem{},
		},
		{
			name: "decimal value not flagged",
			yaml: "---\nvalue: 755\n",
			config: Config{
				ForbidImplicitOctal: true,
				ForbidExplicitOctal: true,
			},
			expected: []types.Problem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config

			require.NotNil(t, config)
			require.Equal(t, "octal-values", rule.ID())
			require.Equal(t, "token", rule.Type())
		})
	}
}
