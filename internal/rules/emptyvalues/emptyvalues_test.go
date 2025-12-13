package emptyvalues

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

	if v, ok := parsed["empty-values"].(string); ok && v == "disable" {
		cfg.Rules["empty-values"] = config.RuleConfig{Level: "disable"}
	} else if v, ok := parsed["empty-values"].(map[string]any); ok {
		cfg.Rules["empty-values"] = config.RuleConfig{
			Level:  "error",
			Config: v,
		}
	} else {
		cfg.Rules["empty-values"] = config.RuleConfig{Level: "error"}
	}

	l := linter.New(cfg)
	l.RegisterRule("empty-values", &Runner{})

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

func TestEmptyValues_Disabled(t *testing.T) {
	conf := "empty-values: disable\nbraces: disable\ncommas: disable\n"
	check(t, "---\nfoo:\n", conf)
	check(t, "---\n{a:}\n", conf)
}

func TestEmptyValues_BlockMappings(t *testing.T) {
	conf := "empty-values:\n  forbid-in-block-mappings: true\n"
	check(t,
		"---\nfoo:\nbar: 1\n",
		conf,
		types.Problem{Line: 2, Column: 5},
	)
}

func TestEmptyValues_FlowMappings(t *testing.T) {
	conf := "empty-values:\n  forbid-in-flow-mappings: true\nbraces: disable\ncommas: disable\n"
	check(t,
		"---\n{a:}\n",
		conf,
		types.Problem{Line: 2, Column: 4},
	)
}

func TestEmptyValues_BlockSequences(t *testing.T) {
	conf := "empty-values:\n  forbid-in-block-sequences: true\n"
	check(t,
		"---\n- \n- foo\n",
		conf,
		types.Problem{Line: 2, Column: 2},
	)
}
