package braces

import (
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	Forbid               any `yaml:"forbid"`
	MinSpacesInside      int `yaml:"min-spaces-inside"`
	MaxSpacesInside      int `yaml:"max-spaces-inside"`
	MinSpacesInsideEmpty int `yaml:"min-spaces-inside-empty"`
	MaxSpacesInsideEmpty int `yaml:"max-spaces-inside-empty"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "braces"
}

func (r *Rule) Type() string {
	return "token"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		Forbid:               false,
		MinSpacesInside:      0,
		MaxSpacesInside:      0,
		MinSpacesInsideEmpty: -1,
		MaxSpacesInsideEmpty: -1,
	}
}

func (r *Rule) Check(config Config, token *types.Token, prev, next, nextNext *types.Token, context map[string]any) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if token.Type == types.FlowMappingStartToken {
			if config.Forbid == true {
				if next != nil {
					if !yield(types.Problem{
						Line:   next.StartMark.Line,
						Column: next.StartMark.Column,
						Desc:   "forbidden flow mapping",
					}) {
						return
					}
				}
			} else if str, ok := config.Forbid.(string); ok && str == "non-empty" {
				if next != nil && next.Type != types.FlowMappingEndToken {
					if !yield(types.Problem{
						Line:   next.StartMark.Line,
						Column: next.StartMark.Column,
						Desc:   "forbidden non-empty flow mapping",
					}) {
						return
					}
				}
			}

			if next != nil {
				spaces := next.StartMark.Index - token.EndMark.Index

				if next.Type == types.FlowMappingEndToken {
					minEmpty := config.MinSpacesInsideEmpty
					if minEmpty == -1 {
						minEmpty = config.MinSpacesInside
					}

					maxEmpty := config.MaxSpacesInsideEmpty
					if maxEmpty == -1 {
						maxEmpty = config.MaxSpacesInside
					}

					if minEmpty != -1 && spaces < minEmpty {
						if !yield(types.Problem{
							Line:   token.StartMark.Line,
							Column: token.StartMark.Column,
							Desc:   "too few spaces inside empty braces",
						}) {
							return
						}
					}
					if maxEmpty != -1 && spaces > maxEmpty {
						if !yield(types.Problem{
							Line:   token.StartMark.Line,
							Column: token.StartMark.Column,
							Desc:   "too many spaces inside empty braces",
						}) {
							return
						}
					}
				} else {
					if config.MinSpacesInside != -1 && spaces < config.MinSpacesInside {
						if !yield(types.Problem{
							Line:   token.StartMark.Line,
							Column: token.StartMark.Column,
							Desc:   "too few spaces inside braces",
						}) {
							return
						}
					}
					if config.MaxSpacesInside != -1 && spaces > config.MaxSpacesInside {
						if token.EndMark.Line == next.StartMark.Line {
							if !yield(types.Problem{
								Line:   token.StartMark.Line,
								Column: token.StartMark.Column,
								Desc:   "too many spaces inside braces",
							}) {
								return
							}
						}
					}
				}
			}
		} else if token.Type == types.FlowMappingEndToken {
			if prev != nil && prev.Type != types.FlowMappingStartToken {
				spaces := token.StartMark.Index - prev.EndMark.Index

				if config.MinSpacesInside != -1 && spaces < config.MinSpacesInside {
					if !yield(types.Problem{
						Line:   token.StartMark.Line,
						Column: token.StartMark.Column,
						Desc:   "too few spaces inside braces",
					}) {
						return
					}
				}
				if config.MaxSpacesInside != -1 && spaces > config.MaxSpacesInside {
					if prev.EndMark.Line == token.StartMark.Line {
						if !yield(types.Problem{
							Line:   token.StartMark.Line,
							Column: token.StartMark.Column,
							Desc:   "too many spaces inside braces",
						}) {
							return
						}
					}
				}
			}
		}
	}
}
