package linelength

import (
	"sort"
	"strings"
	"testing"

	"github.com/mridang/yamllint-go/internal/config"
	"github.com/mridang/yamllint-go/internal/linter"
	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func check(
	t *testing.T,
	content string,
	confYAML string,
	problems ...types.Problem,
) {
	t.Helper()

	var parsed map[string]any
	require.NoError(t, yaml.Unmarshal([]byte(confYAML), &parsed))

	cfg := &config.Config{Rules: map[string]config.RuleConfig{}}

	if v, ok := parsed["line-length"]; ok {
		if s, ok := v.(string); ok && s == "disable" {
			cfg.Rules["line-length"] = config.RuleConfig{Level: "disable"}
		} else {
			cfg.Rules["line-length"] = config.RuleConfig{
				Level:  "error",
				Config: v.(map[string]any),
			}
		}
	}

	l := linter.New(cfg)
	l.RegisterRule("line-length", &Runner{})

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

	require.Equal(t, len(problems), len(actual))
	for i := range problems {
		require.Equal(t, problems[i].Line, actual[i].Line)
		require.Equal(t, problems[i].Column, actual[i].Column)
		require.Contains(t, actual[i].Desc, problems[i].Desc)
	}
}

func TestLineLength_Default(t *testing.T) {
	conf := `
line-length: {max: 80}
empty-lines: disable
new-line-at-end-of-file: disable
document-start: disable
`

	// PASS (non-breakable words allowed by default)
	check(t, strings.Repeat("a", 81), conf)
	check(t, "---\n"+strings.Repeat("a", 81)+"\n", conf)

	// FAIL (breakable)
	check(
		t,
		strings.Repeat("aaaa ", 16)+"z",
		conf,
		types.Problem{Line: 1, Column: 81, Desc: "line too long"},
	)
}

func TestLineLength_Max10(t *testing.T) {
	conf := `
line-length: {max: 10}
new-line-at-end-of-file: disable
`

	check(t, "---\nABCD EFGHI", conf)
	check(
		t,
		"---\nABCD EFGHIJ",
		conf,
		types.Problem{Line: 2, Column: 11, Desc: "line too long"},
	)
}

func TestLineLength_NonBreakableWords(t *testing.T) {
	conf := `line-length: {max: 20, allow-non-breakable-words: true}`

	check(t, "---\n"+strings.Repeat("A", 30)+"\n", conf)

	check(
		t,
		"---\nlong_line: http://localhost/very/very/long/url\n",
		conf,
		types.Problem{Line: 2, Column: 21, Desc: "line too long"},
	)
}

func TestLineLength_InlineMappings(t *testing.T) {
	conf := `line-length: {max: 20, allow-non-breakable-inline-mappings: true}`

	// PASS
	check(
		t,
		"---\nlong_line: http://localhost/very/very/long/url\n",
		conf,
	)

	// FAIL
	check(
		t,
		"---\nlong line: http://localhost/very/very/long/url\n",
		conf,
		types.Problem{Line: 2, Column: 21, Desc: "line too long"},
	)
}
