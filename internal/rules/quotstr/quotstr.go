package quotedstrings

import (
	"fmt"
	"iter"
	"regexp"
	"strings"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	QuoteType              string `yaml:"quote-type"`
	RequiredOnlyForStrings bool   `yaml:"required"`
	CheckKeys              bool   `yaml:"check-keys"`
	AllowQuotedQuotes      bool   `yaml:"allow-quoted-quotes"`
	ExtraRequired          []string `yaml:"extra-required"`
	ExtraAllowed           []string `yaml:"extra-allowed"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "quoted-strings"
}

func (r *Rule) Type() string {
	return "token"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		QuoteType:              "any",
		RequiredOnlyForStrings: true,
		CheckKeys:              false,
		AllowQuotedQuotes:      false,
		ExtraRequired:          []string{},
		ExtraAllowed:           []string{},
	}
}

func getQuotedStringType(token *types.Token) string {
	if token.Style == "'" {
		return "single"
	} else if token.Style == "\"" {
		return "double"
	}
	return ""
}

func (r *Rule) Check(config Config, token *types.Token, prev, next, nextNext *types.Token, context map[string]any) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		quoteType := config.QuoteType

		if prev != nil && prev.Type == types.TagToken {
			return
		}

		if token.Type != types.ScalarToken {
			return
		}

		if prev != nil && prev.Type == types.KeyToken && !config.CheckKeys {
			return
		}

		valueQuotedType := getQuotedStringType(token)

		extraAllowedRegexes := []*regexp.Regexp{}
		for _, pattern := range config.ExtraAllowed {
			if re, err := regexp.Compile(pattern); err == nil {
				extraAllowedRegexes = append(extraAllowedRegexes, re)
			}
		}

		for _, re := range extraAllowedRegexes {
			if re.MatchString(token.Value) {
				return
			}
		}

		extraRequiredRegexes := []*regexp.Regexp{}
		for _, pattern := range config.ExtraRequired {
			if re, err := regexp.Compile(pattern); err == nil {
				extraRequiredRegexes = append(extraRequiredRegexes, re)
			}
		}

		extraRequired := false
		for _, re := range extraRequiredRegexes {
			if re.MatchString(token.Value) {
				extraRequired = true
				break
			}
		}

		if quoteType == "single" || quoteType == "double" {
			wantedQuoteType := quoteType
			if extraRequired && valueQuotedType == "" {
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Desc:   fmt.Sprintf("string value is not quoted with %s quotes", wantedQuoteType),
				}) {
					return
				}
			} else if valueQuotedType != "" && valueQuotedType != wantedQuoteType {
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Desc:   fmt.Sprintf("string value is not quoted with %s quotes", wantedQuoteType),
				}) {
					return
				}
			}
		}

		if (config.RequiredOnlyForStrings || extraRequired) && valueQuotedType == "" {
			isString := false

			if token.Plain {
				lowerVal := strings.ToLower(token.Value)
				if lowerVal == "null" || lowerVal == "true" || lowerVal == "false" ||
					lowerVal == ".inf" || lowerVal == "-.inf" || lowerVal == ".nan" {
					return
				}

				matched, _ := regexp.MatchString(`^[-+]?[0-9]`, token.Value)
				if matched {
					return
				}

				matched, _ = regexp.MatchString(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`, token.Value)
				if matched {
					return
				}

				isString = true
			}

			if isString || extraRequired {
				msg := "string value is not quoted"
				if quoteType == "single" || quoteType == "double" {
					msg = fmt.Sprintf("string value is not quoted with %s quotes", quoteType)
				}

				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Desc:   msg,
				}) {
					return
				}
			}
		}

		if valueQuotedType != "" && !config.AllowQuotedQuotes {
			hasQuoteChar := strings.Contains(token.Value, "'") || strings.Contains(token.Value, "\"")
			if !hasQuoteChar {
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Desc:   "quoted string without quote characters",
				}) {
					return
				}
			}
		}
	}
}
