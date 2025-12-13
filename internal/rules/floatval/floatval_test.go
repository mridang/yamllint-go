// file: internal/rules/floatval/floatval_test.go
package floatvalues

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

	if v, ok := parsed["float-values"].(string); ok && v == "disable" {
		cfg.Rules["float-values"] = config.RuleConfig{Level: "disable"}
	} else if v, ok := parsed["float-values"].(map[string]any); ok {
		cfg.Rules["float-values"] = config.RuleConfig{
			Level:  "error",
			Config: v,
		}
	} else {
		cfg.Rules["float-values"] = config.RuleConfig{Level: "error"}
	}

	l := linter.New(cfg)
	l.RegisterRule("float-values", &Runner{})

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

func TestFloatValues_Disabled(t *testing.T) {
	check(t,
		"---\n- 0.0\n- .NaN\n- .INF\n- .1\n- 10e-6\n",
		"float-values: disable\n",
	)
}

func TestFloatValues_NumeralBeforeDecimal(t *testing.T) {
	conf := "" +
		"float-values:\n" +
		"  require-numeral-before-decimal: true\n" +
		"  forbid-scientific-notation: false\n" +
		"  forbid-nan: false\n" +
		"  forbid-inf: false\n"

	check(t,
		"---\n"+
			"- 0.0\n"+
			"- .1\n"+
			"- '.1'\n"+
			"- string.1\n"+
			"- .1string\n"+
			"- !custom_tag .2\n"+
			"- &angle1 0.0\n"+
			"- *angle1\n"+
			"- &angle2 .3\n"+
			"- *angle2\n",
		conf,
		types.Problem{Line: 3, Column: 3},
		types.Problem{Line: 10, Column: 11},
	)
}

func TestFloatValues_Scientific(t *testing.T) {
	conf := "" +
		"float-values:\n" +
		"  require-numeral-before-decimal: false\n" +
		"  forbid-scientific-notation: true\n" +
		"  forbid-nan: false\n" +
		"  forbid-inf: false\n"

	check(t,
		"---\n"+
			"- 10e6\n"+
			"- 10e-6\n"+
			"- 0.00001\n"+
			"- '10e-6'\n"+
			"- string10e-6\n"+
			"- 10e-6string\n"+
			"- !custom_tag 10e-6\n"+
			"- &angle1 0.000001\n"+
			"- *angle1\n"+
			"- &angle2 10e-6\n"+
			"- *angle2\n"+
			"- &angle3 10e6\n"+
			"- *angle3\n",
		conf,
		types.Problem{Line: 2, Column: 3},
		types.Problem{Line: 3, Column: 3},
		types.Problem{Line: 11, Column: 11},
		types.Problem{Line: 13, Column: 11},
	)
}

func TestFloatValues_NaN(t *testing.T) {
	conf := "" +
		"float-values:\n" +
		"  require-numeral-before-decimal: false\n" +
		"  forbid-scientific-notation: false\n" +
		"  forbid-nan: true\n" +
		"  forbid-inf: false\n"

	check(t,
		"---\n"+
			"- .NaN\n"+
			"- &a .nan\n"+
			"- *a\n",
		conf,
		types.Problem{Line: 2, Column: 3},
		types.Problem{Line: 3, Column: 6},
	)
}

func TestFloatValues_Inf(t *testing.T) {
	conf := "" +
		"float-values:\n" +
		"  require-numeral-before-decimal: false\n" +
		"  forbid-scientific-notation: false\n" +
		"  forbid-nan: false\n" +
		"  forbid-inf: true\n"

	check(t,
		"---\n"+
			"- .inf\n"+
			"- -.INF\n"+
			"- &a .inf\n"+
			"- *a\n",
		conf,
		types.Problem{Line: 2, Column: 3},
		types.Problem{Line: 3, Column: 3},
		types.Problem{Line: 4, Column: 6},
	)
}
