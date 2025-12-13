package anchors

import (
	"fmt"
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type anchorInfo struct {
	Line   int
	Column int
	Used   bool
}

type Config struct {
	ForbidUndeclaredAliases bool `yaml:"forbid-undeclared-aliases"`
	ForbidDuplicatedAnchors bool `yaml:"forbid-duplicated-anchors"`
	ForbidUnusedAnchors     bool `yaml:"forbid-unused-anchors"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "anchors"
}

func (r *Rule) Type() string {
	return "token"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		ForbidUndeclaredAliases: true,
		ForbidDuplicatedAnchors: false,
		ForbidUnusedAnchors:     false,
	}
}

func (r *Rule) Check(config Config, token *types.Token, prev, next, nextNext *types.Token, context map[string]any) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if config.ForbidUndeclaredAliases || config.ForbidDuplicatedAnchors || config.ForbidUnusedAnchors {
			switch token.Type {
			case types.StreamStartToken, types.DocumentStartToken, types.DocumentEndToken:
				context["anchors"] = make(map[string]*anchorInfo)
			}
		}

		anchors, ok := context["anchors"].(map[string]*anchorInfo)
		if !ok {
			anchors = make(map[string]*anchorInfo)
			context["anchors"] = anchors
		}

		if config.ForbidUnusedAnchors && token.Type == types.AliasToken {
			if info, exists := anchors[token.Value]; exists {
				info.Used = true
			}
		}

		if config.ForbidUndeclaredAliases && token.Type == types.AliasToken {
			if _, exists := anchors[token.Value]; !exists {
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Desc:   fmt.Sprintf("found undeclared alias \"%s\"", token.Value),
				}) {
					return
				}
			}
		}

		if config.ForbidDuplicatedAnchors && token.Type == types.AnchorToken {
			if _, exists := anchors[token.Value]; exists {
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Desc:   fmt.Sprintf("found duplicated anchor \"%s\"", token.Value),
				}) {
					return
				}
			}
		}

		if config.ForbidUnusedAnchors {
			isEndOfBlock := false
			if next != nil {
				switch next.Type {
				case types.StreamEndToken, types.DocumentStartToken, types.DocumentEndToken:
					isEndOfBlock = true
				}
			}

			if isEndOfBlock {
				for anchor, info := range anchors {
					if !info.Used {
						if !yield(types.Problem{
							Line:   info.Line,
							Column: info.Column,
							Desc:   fmt.Sprintf("found unused anchor \"%s\"", anchor),
						}) {
							return
						}
					}
				}
			}
		}

		if config.ForbidUndeclaredAliases || config.ForbidDuplicatedAnchors || config.ForbidUnusedAnchors {
			if token.Type == types.AnchorToken {
				anchors[token.Value] = &anchorInfo{
					Line:   token.StartMark.Line,
					Column: token.StartMark.Column,
					Used:   false,
				}
			}
		}
	}
}
