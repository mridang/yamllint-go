package truthy

import (
	"sort"
	"testing"

	"github.com/mridang/yamllint-go/internal/config"
	"github.com/mridang/yamllint-go/internal/linter"
	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func check(t *testing.T, content string, confYAML string, problems ...types.Problem) {
	t.Helper()

	var parsed map[string]any
	require.NoError(t, yaml.Unmarshal([]byte(confYAML), &parsed))

	cfg := &config.Config{Rules: map[string]config.RuleConfig{}}

	if v, ok := parsed["truthy"].(string); ok && v == "disable" {
		cfg.Rules["truthy"] = config.RuleConfig{Level: "disable"}
	} else if v, ok := parsed["truthy"].(map[string]any); ok {
		cfg.Rules["truthy"] = config.RuleConfig{Level: "error", Config: v}
	} else {
		cfg.Rules["truthy"] = config.RuleConfig{Level: "error"}
	}

	l := linter.New(cfg)
	l.RegisterRule("truthy", &Runner{})

	actual, err := l.Run([]byte(content))
	require.NoError(t, err)

	sort.Slice(actual, func(i, j int) bool {
		if actual[i].Line != actual[j].Line {
			return actual[i].Line < actual[j].Line
		}
		return actual[i].Column < actual[j].Column
	})

	require.Equal(t, len(problems), len(actual))
	for i := range problems {
		require.Equal(t, problems[i].Line, actual[i].Line)
		require.Equal(t, problems[i].Column, actual[i].Column)
	}
}

func TestTruthy_Disabled(t *testing.T) {
	check(t, "---\n1: True\n", "truthy: disable")
	check(t, "---\nTrue: 1\n", "truthy: disable")
}

func TestTruthy_Enabled(t *testing.T) {
	conf := "truthy: enable\ndocument-start: disable\n"
	check(t,
		"---\n1: True\nTrue: 1\n",
		conf,
		types.Problem{Line: 2, Column: 4},
		types.Problem{Line: 3, Column: 1},
	)
}

func TestTruthy_CustomAllowed(t *testing.T) {
	conf := "truthy:\n  allowed-values: [\"yes\", \"no\"]\n"
	check(t,
		"---\nkey1: true\nkey2: Yes\nkey3: false\nkey4: no\n",
		conf,
		types.Problem{Line: 2, Column: 7},
		types.Problem{Line: 3, Column: 7},
		types.Problem{Line: 4, Column: 7},
	)
}

func TestTruthy_CheckKeysDisabled(t *testing.T) {
	conf := "truthy:\n  allowed-values: []\n  check-keys: false\n"
	check(t,
		"---\nYES: 0\nYes: 0\nyes: 0\n",
		conf,
	)
}
