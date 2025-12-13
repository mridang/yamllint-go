package octalvalues

import (
	"sort"
	"testing"

	"github.com/mridang/yamllint-go/internal/config"
	"github.com/mridang/yamllint-go/internal/linter"
	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func check(t *testing.T, yamlContent string, confYAML string, problems ...types.Problem) {
	t.Helper()

	var parsedConf map[string]any
	err := yaml.Unmarshal([]byte(confYAML), &parsedConf)
	require.NoError(t, err)

	cfg := &config.Config{
		Rules: map[string]config.RuleConfig{},
	}

	if ruleConf, ok := parsedConf["octal-values"].(string); ok && ruleConf == "disable" {
		cfg.Rules["octal-values"] = config.RuleConfig{Level: "disable"}
	} else if ruleConf, ok := parsedConf["octal-values"].(map[string]any); ok {
		cfg.Rules["octal-values"] = config.RuleConfig{
			Level:  "error",
			Config: ruleConf,
		}
	}

	l := linter.New(cfg)
	l.RegisterRule("octal-values", &Runner{})

	actualProblems, err := l.Run([]byte(yamlContent))
	require.NoError(t, err)

	sort.Slice(actualProblems, func(i, j int) bool {
		if actualProblems[i].Line != actualProblems[j].Line {
			return actualProblems[i].Line < actualProblems[j].Line
		}
		return actualProblems[i].Column < actualProblems[j].Column
	})

	sort.Slice(problems, func(i, j int) bool {
		if problems[i].Line != problems[j].Line {
			return problems[i].Line < problems[j].Line
		}
		return problems[i].Column < problems[j].Column
	})

	require.Equal(t, len(problems), len(actualProblems), "Wrong number of problems")
	for i := range problems {
		require.Equal(t, problems[i].Line, actualProblems[i].Line, "Problem %d: line mismatch", i)
		require.Equal(t, problems[i].Column, actualProblems[i].Column, "Problem %d: column mismatch", i)
		require.Contains(t, actualProblems[i].Desc, problems[i].Desc, "Problem %d: description mismatch", i)
	}
}

func TestDisabled(t *testing.T) {
	conf := `
octal-values: disable
new-line-at-end-of-file: disable
document-start: disable
`
	check(t, "user-city: 010", conf)
	check(t, "user-city: 0o10", conf)
}

func TestImplicitOctalValues(t *testing.T) {
	conf := `
octal-values:
  forbid-implicit-octal: true
  forbid-explicit-octal: false
new-line-at-end-of-file: disable
document-start: disable
`
	check(t, "after-tag: !custom_tag 010", conf)
	check(t, "user-city: 010", conf,
		types.Problem{Line: 1, Column: 15, Desc: "forbidden implicit octal"},
	)
	check(t, "user-city: abc", conf)
	check(t, "user-city: 010,0571", conf)
	check(t, "user-city: '010'", conf)
	check(t, `user-city: "010"`, conf)
	check(t, "user-city:\n  - 010", conf,
		types.Problem{Line: 2, Column: 8, Desc: "forbidden implicit octal"},
	)
	check(t, "user-city: [010]", conf,
		types.Problem{Line: 1, Column: 16, Desc: "forbidden implicit octal"},
	)
	check(t, "user-city: {beijing: 010}", conf,
		types.Problem{Line: 1, Column: 25, Desc: "forbidden implicit octal"},
	)
	check(t, "explicit-octal: 0o10", conf)
	check(t, "not-number: 0abc", conf)
	check(t, "zero: 0", conf)
	check(t, "hex-value: 0x10", conf)
	check(t, "number-values:\n  - 0.10\n  - .01\n  - 0e3\n", conf)
	check(t, "with-decimal-digits: 012345678", conf)
	check(t, "with-decimal-digits: 012345679", conf)
}

func TestExplicitOctalValues(t *testing.T) {
	conf := `
octal-values:
  forbid-implicit-octal: false
  forbid-explicit-octal: true
new-line-at-end-of-file: disable
document-start: disable
`
	check(t, "user-city: 0o10", conf,
		types.Problem{Line: 1, Column: 16, Desc: "forbidden explicit octal"},
	)
	check(t, "user-city: abc", conf)
	check(t, "user-city: 0o10,0571", conf)
	check(t, "user-city: '0o10'", conf)
	check(t, "user-city:\n  - 0o10", conf,
		types.Problem{Line: 2, Column: 9, Desc: "forbidden explicit octal"},
	)
	check(t, "user-city: [0o10]", conf,
		types.Problem{Line: 1, Column: 17, Desc: "forbidden explicit octal"},
	)
	check(t, "user-city: {beijing: 0o10}", conf,
		types.Problem{Line: 1, Column: 26, Desc: "forbidden explicit octal"},
	)
	check(t, "implicit-octal: 010", conf)
	check(t, "not-number: 0oabc", conf)
	check(t, "zero: 0", conf)
	check(t, "hex-value: 0x10", conf)
	check(t, "number-values:\n  - 0.10\n  - .01\n  - 0e3\n", conf)
	check(t, `user-city: "010"`, conf)
	check(t, "with-decimal-digits: 0o012345678", conf)
	check(t, "with-decimal-digits: 0o012345679", conf)
}
