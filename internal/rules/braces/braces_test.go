package braces

import (
	"fmt"
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

	if bracesConf, ok := parsedConf["braces"].(string); ok && bracesConf == "disable" {
		cfg.Rules["braces"] = config.RuleConfig{Level: "disable"}
	} else if bracesConf, ok := parsedConf["braces"].(map[string]any); ok {
		cfg.Rules["braces"] = config.RuleConfig{
			Level:  "error",
			Config: bracesConf,
		}
	}

	l := linter.New(cfg)
	l.RegisterRule("braces", &Runner{})

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

	// DEBUG output
	if len(problems) != len(actualProblems) {
		fmt.Printf("\n=== TEST FAILED: %s ===\n", t.Name())
		fmt.Printf("YAML:\n%s\n", yamlContent)
		fmt.Printf("Config:\n%s\n", confYAML)
		fmt.Printf("Expected %d problems, got %d\n\n", len(problems), len(actualProblems))

		fmt.Printf("EXPECTED:\n")
		for i, p := range problems {
			fmt.Printf("  [%d] Line=%d Col=%d Desc=%q\n", i, p.Line, p.Column, p.Desc)
		}

		fmt.Printf("\nACTUAL:\n")
		for i, p := range actualProblems {
			fmt.Printf("  [%d] Line=%d Col=%d Desc=%q\n", i, p.Line, p.Column, p.Desc)
		}
		fmt.Println()
	}

	require.Equal(t, len(problems), len(actualProblems), "Wrong number of problems")
	for i := range problems {
		require.Equal(t, problems[i].Line, actualProblems[i].Line, "Problem %d: line mismatch", i)
		require.Equal(t, problems[i].Column, actualProblems[i].Column, "Problem %d: column mismatch", i)
		require.Contains(t, actualProblems[i].Desc, problems[i].Desc, "Problem %d: description mismatch", i)
	}
}

func TestForbid(t *testing.T) {
	conf := `braces:
  forbid: false
`
	check(t, "---\ndict: {}\n", conf)
	check(t, "---\ndict: {a}\n", conf)
	check(t, "---\ndict: {a: 1}\n", conf)
	check(t, "---\ndict: {\n  a: 1\n}\n", conf)

	conf = `braces:
  forbid: true
`
	check(t, "---\ndict:\n  a: 1\n", conf)
	check(t, "---\ndict: {}\n", conf,
		types.Problem{Line: 2, Column: 8, Desc: "forbidden"},
	)
}

func TestMinSpacesEmpty(t *testing.T) {
	conf := `braces:
  max-spaces-inside: -1
  min-spaces-inside: -1
  max-spaces-inside-empty: 0
  min-spaces-inside-empty: 0
`
	check(t, "---\narray: {}\n", conf)
}
