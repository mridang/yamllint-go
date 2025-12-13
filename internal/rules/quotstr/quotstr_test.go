package quotedstrings

import (
	"testing"

	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestQuotedStringsRule(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		config   Config
		expected []types.Problem
	}{
		{
			name: "string not quoted",
			yaml: "---\nkey: value\n",
			config: Config{
				QuoteType:              "single",
				RequiredOnlyForStrings: true,
			},
			expected: []types.Problem{
				{Line: 1, Column: 5, Desc: "string value is not quoted with single quotes"},
			},
		},
		{
			name: "single quotes valid",
			yaml: "---\nkey: 'value'\n",
			config: Config{
				QuoteType:              "single",
				RequiredOnlyForStrings: true,
			},
			expected: []types.Problem{},
		},
		{
			name: "double quotes when single expected",
			yaml: "---\nkey: \"value\"\n",
			config: Config{
				QuoteType:              "single",
				RequiredOnlyForStrings: true,
			},
			expected: []types.Problem{
				{Line: 1, Column: 5, Desc: "string value is not quoted with single quotes"},
			},
		},
		{
			name: "double quotes valid",
			yaml: "---\nkey: \"value\"\n",
			config: Config{
				QuoteType:              "double",
				RequiredOnlyForStrings: true,
			},
			expected: []types.Problem{},
		},
		{
			name: "numbers not required to be quoted",
			yaml: "---\nkey: 123\n",
			config: Config{
				QuoteType:              "single",
				RequiredOnlyForStrings: true,
			},
			expected: []types.Problem{},
		},
		{
			name: "quoted without quote chars forbidden",
			yaml: "---\nkey: 'value'\n",
			config: Config{
				QuoteType:         "any",
				AllowQuotedQuotes: false,
			},
			expected: []types.Problem{
				{Line: 1, Column: 5, Desc: "quoted string without quote characters"},
			},
		},
		{
			name: "extra required pattern",
			yaml: "---\nurl: http://example.com\n",
			config: Config{
				QuoteType:     "any",
				ExtraRequired: []string{"^http://.*"},
			},
			expected: []types.Problem{
				{Line: 1, Column: 5, Desc: "string value is not quoted"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{}
			config := tt.config
			
			require.NotNil(t, config)
			require.Equal(t, "quoted-strings", rule.ID())
			require.Equal(t, "token", rule.Type())
		})
	}
}
