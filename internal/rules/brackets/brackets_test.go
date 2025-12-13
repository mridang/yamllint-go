package brackets

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

	if bracketsConf, ok := parsedConf["brackets"].(string); ok && bracketsConf == "disable" {
		cfg.Rules["brackets"] = config.RuleConfig{Level: "disable"}
	} else if bracketsConf, ok := parsedConf["brackets"].(map[string]any); ok {
		cfg.Rules["brackets"] = config.RuleConfig{
			Level:  "error",
			Config: bracketsConf,
		}
	}

	l := linter.New(cfg)
	l.RegisterRule("brackets", &Runner{})

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
	conf := "brackets: disable"
	check(t, `---
array1: []
array2: [ ]
array3: [   a, b]
array4: [a, b, c ]
array5: [a, b, c ]
array6: [  a, b, c ]
array7: [   a, b, c ]
`, conf)
}

func TestForbid(t *testing.T) {
	conf := `brackets:
  forbid: false
`
	check(t, "---\narray: []\n", conf)
	check(t, "---\narray: [a, b]\n", conf)
	check(t, "---\narray: [\n  a,\n  b\n]\n", conf)

	conf = `brackets:
  forbid: true
`
	check(t, "---\narray:\n  - a\n  - b\n", conf)
	check(t, "---\narray: []\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "forbidden"},
	)
	check(t, "---\narray: [a, b]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "forbidden"},
	)
	check(t, "---\narray: [\n  a,\n  b\n]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "forbidden"},
	)

	conf = `brackets:
  forbid: non-empty
`
	check(t, "---\narray:\n  - a\n  - b\n", conf)
	check(t, "---\narray: []\n", conf)
	check(t, "---\narray: [\n\n]\n", conf)
	check(t, "---\narray: [\n# a comment\n]\n", conf)
	check(t, "---\narray: [a, b]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "forbidden"},
	)
	check(t, "---\narray: [\n  a,\n  b\n]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "forbidden"},
	)
}

func TestMinSpaces(t *testing.T) {
	conf := `brackets:
  max-spaces-inside: -1
  min-spaces-inside: 0
  max-spaces-inside-empty: -1
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: []\n", conf)

	conf = `brackets:
  max-spaces-inside: -1
  min-spaces-inside: 1
  max-spaces-inside-empty: -1
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: []\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
	)
	check(t, "---\narray: [ ]\n", conf)
	check(t, "---\narray: [a, b]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
		types.Problem{Line: 2, Column: 13, Desc: "too few spaces"},
	)
	check(t, "---\narray: [ a, b ]\n", conf)
	check(t, "---\narray: [\n  a,\n  b\n]\n", conf)

	conf = `brackets:
  max-spaces-inside: -1
  min-spaces-inside: 3
  max-spaces-inside-empty: -1
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: [ a, b ]\n", conf,
		types.Problem{Line: 2, Column: 10, Desc: "too few spaces"},
		types.Problem{Line: 2, Column: 15, Desc: "too few spaces"},
	)
	check(t, "---\narray: [   a, b   ]\n", conf)
}

func TestMaxSpaces(t *testing.T) {
	conf := `brackets:
  max-spaces-inside: 0
  min-spaces-inside: -1
  max-spaces-inside-empty: -1
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: []\n", conf)
	check(t, "---\narray: [ ]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too many spaces"},
	)
	check(t, "---\narray: [a, b]\n", conf)
	check(t, "---\narray: [ a, b ]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too many spaces"},
		types.Problem{Line: 2, Column: 14, Desc: "too many spaces"},
	)
	check(t, "---\narray: [   a, b   ]\n", conf,
		types.Problem{Line: 2, Column: 11, Desc: "too many spaces"},
		types.Problem{Line: 2, Column: 18, Desc: "too many spaces"},
	)
	check(t, "---\narray: [\n  a,\n  b\n]\n", conf)

	conf = `brackets:
  max-spaces-inside: 3
  min-spaces-inside: -1
  max-spaces-inside-empty: -1
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: [   a, b   ]\n", conf)
	check(t, "---\narray: [    a, b     ]\n", conf,
		types.Problem{Line: 2, Column: 12, Desc: "too many spaces"},
		types.Problem{Line: 2, Column: 21, Desc: "too many spaces"},
	)
}

func TestMinAndMaxSpaces(t *testing.T) {
	conf := `brackets:
  max-spaces-inside: 0
  min-spaces-inside: 0
  max-spaces-inside-empty: -1
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: []\n", conf)
	check(t, "---\narray: [ ]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too many spaces"},
	)
	check(t, "---\narray: [   a, b]\n", conf,
		types.Problem{Line: 2, Column: 11, Desc: "too many spaces"},
	)

	conf = `brackets:
  max-spaces-inside: 1
  min-spaces-inside: 1
  max-spaces-inside-empty: -1
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: [a, b, c ]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
	)

	conf = `brackets:
  max-spaces-inside: 2
  min-spaces-inside: 0
  max-spaces-inside-empty: -1
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: [a, b, c ]\n", conf)
	check(t, "---\narray: [  a, b, c ]\n", conf)
	check(t, "---\narray: [   a, b, c ]\n", conf,
		types.Problem{Line: 2, Column: 11, Desc: "too many spaces"},
	)
}

func TestMinSpacesEmpty(t *testing.T) {
	conf := `brackets:
  max-spaces-inside: -1
  min-spaces-inside: -1
  max-spaces-inside-empty: 0
  min-spaces-inside-empty: 0
`
	check(t, "---\narray: []\n", conf)

	conf = `brackets:
  max-spaces-inside: -1
  min-spaces-inside: -1
  max-spaces-inside-empty: -1
  min-spaces-inside-empty: 1
`
	check(t, "---\narray: []\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
	)
	check(t, "---\narray: [ ]\n", conf)

	conf = `brackets:
  max-spaces-inside: -1
  min-spaces-inside: -1
  max-spaces-inside-empty: -1
  min-spaces-inside-empty: 3
`
	check(t, "---\narray: []\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
	)
	check(t, "---\narray: [   ]\n", conf)
}

func TestMaxSpacesEmpty(t *testing.T) {
	conf := `brackets:
  max-spaces-inside: -1
  min-spaces-inside: -1
  max-spaces-inside-empty: 0
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: []\n", conf)
	check(t, "---\narray: [ ]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too many spaces"},
	)

	conf = `brackets:
  max-spaces-inside: -1
  min-spaces-inside: -1
  max-spaces-inside-empty: 1
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: []\n", conf)
	check(t, "---\narray: [ ]\n", conf)
	check(t, "---\narray: [  ]\n", conf,
		types.Problem{Line: 2, Column: 10, Desc: "too many spaces"},
	)

	conf = `brackets:
  max-spaces-inside: -1
  min-spaces-inside: -1
  max-spaces-inside-empty: 3
  min-spaces-inside-empty: -1
`
	check(t, "---\narray: []\n", conf)
	check(t, "---\narray: [   ]\n", conf)
	check(t, "---\narray: [    ]\n", conf,
		types.Problem{Line: 2, Column: 12, Desc: "too many spaces"},
	)
}

func TestMinAndMaxSpacesEmpty(t *testing.T) {
	conf := `brackets:
  max-spaces-inside: -1
  min-spaces-inside: -1
  max-spaces-inside-empty: 2
  min-spaces-inside-empty: 1
`
	check(t, "---\narray: []\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
	)
	check(t, "---\narray: [ ]\n", conf)
	check(t, "---\narray: [  ]\n", conf)
	check(t, "---\narray: [   ]\n", conf,
		types.Problem{Line: 2, Column: 11, Desc: "too many spaces"},
	)
}

func TestMixedEmptyNonempty(t *testing.T) {
	conf := `brackets:
  max-spaces-inside: -1
  min-spaces-inside: 1
  max-spaces-inside-empty: 0
  min-spaces-inside-empty: 0
`
	check(t, "---\narray: [ a, b ]\n", conf)
	check(t, "---\narray: [a, b]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
		types.Problem{Line: 2, Column: 13, Desc: "too few spaces"},
	)
	check(t, "---\narray: []\n", conf)
	check(t, "---\narray: [ ]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too many spaces"},
	)

	conf = `brackets:
  max-spaces-inside: 0
  min-spaces-inside: -1
  max-spaces-inside-empty: 1
  min-spaces-inside-empty: 1
`
	check(t, "---\narray: [ a, b ]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too many spaces"},
		types.Problem{Line: 2, Column: 14, Desc: "too many spaces"},
	)
	check(t, "---\narray: [a, b]\n", conf)
	check(t, "---\narray: []\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
	)
	check(t, "---\narray: [ ]\n", conf)

	conf = `brackets:
  max-spaces-inside: 2
  min-spaces-inside: 1
  max-spaces-inside-empty: 1
  min-spaces-inside-empty: 1
`
	check(t, "---\narray: [ a, b  ]\n", conf)
	check(t, "---\narray: [a, b   ]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
		types.Problem{Line: 2, Column: 15, Desc: "too many spaces"},
	)
	check(t, "---\narray: []\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
	)
	check(t, "---\narray: [ ]\n", conf)
	check(t, "---\narray: [   ]\n", conf,
		types.Problem{Line: 2, Column: 11, Desc: "too many spaces"},
	)

	conf = `brackets:
  max-spaces-inside: 1
  min-spaces-inside: 1
  max-spaces-inside-empty: 1
  min-spaces-inside-empty: 1
`
	check(t, "---\narray: [ a, b ]\n", conf)
	check(t, "---\narray: [a, b]\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
		types.Problem{Line: 2, Column: 13, Desc: "too few spaces"},
	)
	check(t, "---\narray: []\n", conf,
		types.Problem{Line: 2, Column: 9, Desc: "too few spaces"},
	)
	check(t, "---\narray: [ ]\n", conf)
}
