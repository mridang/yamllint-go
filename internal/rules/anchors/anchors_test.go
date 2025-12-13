package anchors

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

// check runs the linter on yaml with the given config and compares problems
// This mirrors the Python RuleTestCase.check() method
func check(t *testing.T, yamlContent string, confYAML string, problems ...types.Problem) {
	t.Helper()

	// Parse config YAML string
	var parsedConf map[string]any
	err := yaml.Unmarshal([]byte(confYAML), &parsedConf)
	require.NoError(t, err)

	// Build config
	cfg := &config.Config{
		Rules: map[string]config.RuleConfig{},
	}

	// Handle "anchors: disable" format
	if anchorsConf, ok := parsedConf["anchors"].(string); ok && anchorsConf == "disable" {
		cfg.Rules["anchors"] = config.RuleConfig{Level: "disable"}
	} else if anchorsConf, ok := parsedConf["anchors"].(map[string]any); ok {
		cfg.Rules["anchors"] = config.RuleConfig{
			Level:  "error",
			Config: anchorsConf,
		}
	}

	// Create linter and register anchors rule
	l := linter.New(cfg)
	l.RegisterRule("anchors", &Runner{})

	// Run linter
	actualProblems, err := l.Run([]byte(yamlContent))
	require.NoError(t, err)

	// Sort both (Python does this)
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

	// DEBUG: Print actual vs expected
	if len(problems) != len(actualProblems) {
		fmt.Printf("\n=== PROBLEM COUNT MISMATCH ===\n")
		fmt.Printf("Expected %d problems, got %d\n", len(problems), len(actualProblems))
	}

	fmt.Printf("\n=== EXPECTED vs ACTUAL ===\n")
	for i := 0; i < max(len(problems), len(actualProblems)); i++ {
		if i < len(problems) {
			fmt.Printf("Expected[%d]: Line=%d Col=%d Desc=%q\n", i, problems[i].Line, problems[i].Column, problems[i].Desc)
		}
		if i < len(actualProblems) {
			fmt.Printf("  Actual[%d]: Line=%d Col=%d Desc=%q\n", i, actualProblems[i].Line, actualProblems[i].Column, actualProblems[i].Desc)
		}
		if i < len(problems) && i < len(actualProblems) {
			if problems[i].Line != actualProblems[i].Line {
				fmt.Printf("  ^^^ LINE MISMATCH (off by %d)\n", actualProblems[i].Line-problems[i].Line)
			}
		}
		fmt.Println()
	}

	// Compare
	require.Equal(t, len(problems), len(actualProblems), "Wrong number of problems")
	for i := range problems {
		require.Equal(t, problems[i].Line, actualProblems[i].Line, "Problem %d: line mismatch", i)
		require.Equal(t, problems[i].Column, actualProblems[i].Column, "Problem %d: column mismatch", i)
		require.Equal(t, problems[i].Desc, actualProblems[i].Desc, "Problem %d: description mismatch", i)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestForbidUndeclaredAliases(t *testing.T) {
	conf := `anchors:
  forbid-undeclared-aliases: true
  forbid-duplicated-anchors: false
  forbid-unused-anchors: false
`

	// Second check - undeclared aliases, should fail
	check(t, `---
- &i 42
---
- &b true
- &b true
- &b true
- &s hello
- *b
- *i
- *f_m
- *f_m
- *f_m
- *f_s
- &f_s [1, 2]
...
---
block mapping: &b_m
  key: value
---
block mapping 1: &b_m_bis
  key: value
block mapping 2: &b_m_bis
  key: value
extended:
  <<: *b_m
  foo: bar
---
{a: 1, &x b: 2, c: &x 3, *x : 4, e: *y}
...
`, conf,
		types.Problem{Line: 9, Column: 3, Desc: `found undeclared alias "i"`},
		types.Problem{Line: 10, Column: 3, Desc: `found undeclared alias "f_m"`},
		types.Problem{Line: 11, Column: 3, Desc: `found undeclared alias "f_m"`},
		types.Problem{Line: 12, Column: 3, Desc: `found undeclared alias "f_m"`},
		types.Problem{Line: 13, Column: 3, Desc: `found undeclared alias "f_s"`},
		types.Problem{Line: 25, Column: 7, Desc: `found undeclared alias "b_m"`},
		types.Problem{Line: 28, Column: 37, Desc: `found undeclared alias "y"`},
	)
}
