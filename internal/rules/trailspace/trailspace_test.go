// file: internal/rules/trailspace/trailspace_test.go
package trailingspaces

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

	var parsedConf map[string]any
	err := yaml.Unmarshal([]byte(confYAML), &parsedConf)
	require.NoError(t, err)

	cfg := &config.Config{
		Rules: map[string]config.RuleConfig{},
	}

	if v, ok := parsedConf["trailing-spaces"].(string); ok && v == "disable" {
		cfg.Rules["trailing-spaces"] = config.RuleConfig{Level: "disable"}
	} else {
		cfg.Rules["trailing-spaces"] = config.RuleConfig{Level: "error"}
	}

	l := linter.New(cfg)
	l.RegisterRule("trailing-spaces", &Runner{})

	actual, err := l.Run([]byte(content))
	require.NoError(t, err)

	sort.Slice(actual, func(i, j int) bool {
		if actual[i].Line != actual[j].Line {
			return actual[i].Line < actual[j].Line
		}
		return actual[i].Column < actual[j].Column
	})

	sort.Slice(problems, func(i, j int) bool {
		if problems[i].Line != problems[j].Line {
			return problems[i].Line < problems[j].Line
		}
		return problems[i].Column < problems[j].Column
	})

	require.Equal(t, len(problems), len(actual), "Wrong number of problems")
	for i := range problems {
		require.Equal(t, problems[i].Line, actual[i].Line, "Problem %d: line mismatch", i)
		require.Equal(t, problems[i].Column, actual[i].Column, "Problem %d: column mismatch", i)
		require.Contains(t, actual[i].Desc, problems[i].Desc)
	}
}

func TestTrailingSpaces_Disabled(t *testing.T) {
	conf := `trailing-spaces: disable`
	check(t, "", conf)
	check(t, "\n", conf)
	check(t, "    \n", conf)
	check(t, "---\nsome: text \n", conf)
}

func TestTrailingSpaces_Enabled(t *testing.T) {
	conf := `trailing-spaces: enable`
	check(t, "", conf)
	check(t, "\n", conf)
	check(t, "    \n", conf,
		types.Problem{Line: 1, Column: 1, Desc: "trailing spaces"},
	)
	check(t, "\t\t\t\n", conf,
		types.Problem{Line: 1, Column: 1, Desc: "trailing spaces"},
	)
	check(t, "---\nsome: text \n", conf,
		types.Problem{Line: 2, Column: 11, Desc: "trailing spaces"},
	)
	check(t, "---\nsome: text\t\n", conf,
		types.Problem{Line: 2, Column: 11, Desc: "trailing spaces"},
	)
}

func TestTrailingSpaces_DOS(t *testing.T) {
	conf := `trailing-spaces: enable`
	check(t, "---\r\nsome: text\r\n", conf)
	check(t, "---\r\nsome: text \r\n", conf,
		types.Problem{Line: 2, Column: 11, Desc: "trailing spaces"},
	)
}
