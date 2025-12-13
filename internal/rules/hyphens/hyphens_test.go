package hyphens

import (
	"sort"
	"testing"

	"github.com/mridang/yamllint-go/internal/config"
	"github.com/mridang/yamllint-go/internal/linter"
	"github.com/mridang/yamllint-go/internal/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

/*
Port of tests/rules/test_hyphens.py from yamllint
This MUST match Python behaviour 1:1.
*/

func check(
	t *testing.T,
	yamlContent string,
	confYAML string,
	problems ...types.Problem,
) {
	t.Helper()

	// Parse rule config YAML
	var parsedConf map[string]any
	err := yaml.Unmarshal([]byte(confYAML), &parsedConf)
	require.NoError(t, err)

	cfg := &config.Config{
		Rules: map[string]config.RuleConfig{},
	}

	// Handle disable / enable / config
	if v, ok := parsedConf["hyphens"]; ok {
		switch vv := v.(type) {
		case string:
			if vv == "disable" {
				cfg.Rules["hyphens"] = config.RuleConfig{Level: "disable"}
			}
		case map[string]any:
			cfg.Rules["hyphens"] = config.RuleConfig{
				Level:  "error",
				Config: vv,
			}
		}
	}

	l := linter.New(cfg)
	l.RegisterRule("hyphens", &Runner{})

	actual, err := l.Run([]byte(yamlContent))
	require.NoError(t, err)

	// Sort for stable comparison
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
		require.Contains(t, actual[i].Desc, problems[i].Desc, "Problem %d: description mismatch", i)
	}
}

func TestHyphens_Disabled(t *testing.T) {
	conf := `hyphens: disable`

	check(t,
		"---\n- elem1\n- elem2\n",
		conf,
	)
	check(t,
		"---\n- elem1\n-  elem2\n",
		conf,
	)
	check(t,
		"---\n-  elem1\n-  elem2\n",
		conf,
	)
	check(t,
		"---\n-  elem1\n- elem2\n",
		conf,
	)
	check(t,
		"---\nobject:\n  - elem1\n  -  elem2\n",
		conf,
	)
	check(t,
		"---\nobject:\n  -  elem1\n  -  elem2\n",
		conf,
	)
	check(t,
		"---\nobject:\n  subobject:\n    - elem1\n    -  elem2\n",
		conf,
	)
	check(t,
		"---\nobject:\n  subobject:\n    -  elem1\n    -  elem2\n",
		conf,
	)
}

func TestHyphens_Enabled(t *testing.T) {
	conf := `hyphens: {max-spaces-after: 1}`

	check(t,
		"---\n- elem1\n- elem2\n",
		conf,
	)
	check(t,
		"---\n- elem1\n-  elem2\n",
		conf,
		types.Problem{Line: 3, Column: 3, Desc: "too many spaces after hyphen"},
	)
	check(t,
		"---\n-  elem1\n-  elem2\n",
		conf,
		types.Problem{Line: 2, Column: 3, Desc: "too many spaces after hyphen"},
		types.Problem{Line: 3, Column: 3, Desc: "too many spaces after hyphen"},
	)
	check(t,
		"---\n-  elem1\n- elem2\n",
		conf,
		types.Problem{Line: 2, Column: 3, Desc: "too many spaces after hyphen"},
	)
	check(t,
		"---\nobject:\n  - elem1\n  -  elem2\n",
		conf,
		types.Problem{Line: 4, Column: 5, Desc: "too many spaces after hyphen"},
	)
	check(t,
		"---\nobject:\n  -  elem1\n  -  elem2\n",
		conf,
		types.Problem{Line: 3, Column: 5, Desc: "too many spaces after hyphen"},
		types.Problem{Line: 4, Column: 5, Desc: "too many spaces after hyphen"},
	)
	check(t,
		"---\nobject:\n  subobject:\n    - elem1\n    -  elem2\n",
		conf,
		types.Problem{Line: 5, Column: 7, Desc: "too many spaces after hyphen"},
	)
	check(t,
		"---\nobject:\n  subobject:\n    -  elem1\n    -  elem2\n",
		conf,
		types.Problem{Line: 4, Column: 7, Desc: "too many spaces after hyphen"},
		types.Problem{Line: 5, Column: 7, Desc: "too many spaces after hyphen"},
	)
}

func TestHyphens_Max3(t *testing.T) {
	conf := `hyphens: {max-spaces-after: 3}`

	check(t,
		"---\n-   elem1\n-   elem2\n",
		conf,
	)
	check(t,
		"---\n-    elem1\n-   elem2\n",
		conf,
		types.Problem{Line: 2, Column: 5, Desc: "too many spaces after hyphen"},
	)
	check(t,
		"---\na:\n  b:\n    -   elem1\n    -   elem2\n",
		conf,
	)
	check(t,
		"---\na:\n  b:\n    -    elem1\n    -    elem2\n",
		conf,
		types.Problem{Line: 4, Column: 9, Desc: "too many spaces after hyphen"},
		types.Problem{Line: 5, Column: 9, Desc: "too many spaces after hyphen"},
	)
}
