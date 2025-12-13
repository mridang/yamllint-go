package comments

import (
	"fmt"
	"iter"

	"github.com/mridang/yamllint-go/internal/types"
)

type Config struct {
	RequireStartingSpace bool `yaml:"require-starting-space"`
	IgnoreShebangs       bool `yaml:"ignore-shebangs"`
	MinSpacesFromContent int  `yaml:"min-spaces-from-content"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "comments"
}

func (r *Rule) Type() string {
	return "comment"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		RequireStartingSpace: true,
		IgnoreShebangs:       true,
		MinSpacesFromContent: 2,
	}
}

func (r *Rule) Check(config Config, comment *types.Comment) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if config.MinSpacesFromContent != -1 && comment.IsInline() {
			if comment.TokenBefore != nil {
				gap := comment.Pointer - comment.TokenBefore.EndMark.Index
				if gap < config.MinSpacesFromContent {
					if !yield(types.Problem{
						Line:   comment.LineNo,
						Column: comment.ColumnNo,
						Desc: fmt.Sprintf("too few spaces before comment: expected %d",
							config.MinSpacesFromContent),
					}) {
						return
					}
				}
			}
		}

		if config.RequireStartingSpace {
			textStart := comment.Pointer + 1
			bufferLen := len(comment.Buffer)

			for textStart < bufferLen && comment.Buffer[textStart] == '#' {
				textStart++
			}

			if textStart < bufferLen {
				if config.IgnoreShebangs &&
					comment.LineNo == 1 &&
					comment.ColumnNo == 1 &&
					comment.Buffer[textStart] == '!' {
					return
				}

				char := comment.Buffer[textStart]
				if char != ' ' && char != '\n' && char != '\r' && char != 0 {
					column := comment.ColumnNo + (textStart - comment.Pointer)
					if !yield(types.Problem{
						Line:   comment.LineNo,
						Column: column,
						Desc:   "missing starting space in comment",
					}) {
						return
					}
				}
			}
		}
	}
}
