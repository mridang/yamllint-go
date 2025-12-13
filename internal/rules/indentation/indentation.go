package indentation

import (
	"fmt"
	"iter"
	"strings"

	"github.com/mridang/yamllint-go/internal/common"
	"github.com/mridang/yamllint-go/internal/types"
)

const (
	ROOT = iota
	B_MAP
	F_MAP
	B_SEQ
	F_SEQ
	B_ENT
	KEY
	VAL
)

type Parent struct {
	Type             int
	Indent           int
	LineIndent       int
	ExplicitKey      bool
	ImplicitBlockSeq bool
}

type Config struct {
	Spaces                any  `yaml:"spaces"`
	IndentSequences       any  `yaml:"indent-sequences"`
	CheckMultiLineStrings bool `yaml:"check-multi-line-strings"`
}

type Rule struct{}

func (r *Rule) ID() string {
	return "indentation"
}

func (r *Rule) Type() string {
	return "token"
}

func (r *Rule) DefaultConfig() Config {
	return Config{
		Spaces:                "consistent",
		IndentSequences:       true,
		CheckMultiLineStrings: false,
	}
}

func getRealEndLine(token *types.Token) int {
	if token.Type == types.ScalarToken && token.Style != "" {
		return token.EndMark.Line
	}
	if token.StartMark.Line != token.EndMark.Line {
		return token.EndMark.Line
	}
	return token.StartMark.Line
}

func checkScalarIndentation(config Config, token *types.Token, context map[string]any) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if token.StartMark.Line == token.EndMark.Line {
			return
		}

		stack := context["stack"].([]*Parent)

		computeExpectedIndent := func(foundIndent int) int {
			detect := func(baseIndent int) int {
				spacesVal := context["spaces"]
				if _, isInt := spacesVal.(int); !isInt {
					context["spaces"] = foundIndent - baseIndent
				}
				return baseIndent + context["spaces"].(int)
			}

			if token.Plain {
				return token.StartMark.Column
			} else if token.Style == "\"" || token.Style == "'" {
				return token.StartMark.Column + 1
			} else if token.Style == ">" || token.Style == "|" {
				parent := stack[len(stack)-1]
				if parent.Type == B_ENT {
					return detect(token.StartMark.Column)
				} else if parent.Type == KEY {
					return detect(token.StartMark.Column)
				} else if parent.Type == VAL {
					curLine := context["cur_line"].(int)
					if token.StartMark.Line > curLine {
						return detect(parent.Indent)
					} else if len(stack) >= 2 && stack[len(stack)-2].ExplicitKey {
						return detect(token.StartMark.Column)
					} else {
						if len(stack) >= 2 {
							return detect(stack[len(stack)-2].Indent)
						}
					}
				} else {
					return detect(parent.Indent)
				}
			}
			return 0
		}

		var expectedIndent *int
		lineNo := token.StartMark.Line
		lineStart := token.StartMark.Index

		for {
			nextNewLine := strings.IndexByte(token.StartMark.Buffer[lineStart:token.EndMark.Index-1], '\n')
			if nextNewLine != -1 {
				lineStart = lineStart + nextNewLine + 1
			} else {
				lineStart = 0
			}

			if lineStart == 0 {
				break
			}
			lineNo++

			indent := 0
			for lineStart+indent < len(token.StartMark.Buffer) && token.StartMark.Buffer[lineStart+indent] == ' ' {
				indent++
			}
			if lineStart+indent < len(token.StartMark.Buffer) && token.StartMark.Buffer[lineStart+indent] == '\n' {
				continue
			}

			if expectedIndent == nil {
				val := computeExpectedIndent(indent)
				expectedIndent = &val
			}

			if indent != *expectedIndent {
				if !yield(types.Problem{
					Line:   lineNo,
					Column: indent,
					Desc:   fmt.Sprintf("wrong indentation: expected %d but found %d", *expectedIndent, indent),
				}) {
					return
				}
			}
		}
	}
}

func (r *Rule) Check(config Config, token *types.Token, prev, next, nextNext *types.Token, context map[string]any) iter.Seq[types.Problem] {
	return func(yield func(types.Problem) bool) {
		if _, ok := context["stack"]; !ok {
			context["stack"] = []*Parent{{Type: ROOT, Indent: 0}}
			context["cur_line"] = -1
			context["spaces"] = config.Spaces
			context["indent-sequences"] = config.IndentSequences
		}

		isVisible := token.Type != types.StreamStartToken &&
			token.Type != types.StreamEndToken &&
			token.Type != types.BlockEndToken &&
			!(token.Type == types.ScalarToken && token.Value == "")

		curLine, _ := context["cur_line"].(int)
		firstInLine := isVisible && (token.StartMark.Line > curLine)

		stack := context["stack"].([]*Parent)

		detectIndent := func(baseIndent int, nextTok *types.Token) int {
			spacesVal := context["spaces"]
			if _, isInt := spacesVal.(int); !isInt {
				context["spaces"] = nextTok.StartMark.Column - baseIndent
			}
			return baseIndent + context["spaces"].(int)
		}

		if firstInLine {
			foundIndentation := token.StartMark.Column
			expected := stack[len(stack)-1].Indent

			if token.Type == types.FlowMappingEndToken || token.Type == types.FlowSequenceEndToken {
				expected = stack[len(stack)-1].LineIndent
			} else if stack[len(stack)-1].Type == KEY &&
				stack[len(stack)-1].ExplicitKey &&
				token.Type != types.ValueToken {
				expected = detectIndent(expected, token)
			}

			if foundIndentation != expected {
				var msg string
				if expected < 0 {
					msg = fmt.Sprintf("wrong indentation: expected at least %d", foundIndentation+1)
				} else {
					msg = fmt.Sprintf("wrong indentation: expected %d but found %d", expected, foundIndentation)
				}
				if !yield(types.Problem{
					Line:   token.StartMark.Line,
					Column: foundIndentation,
					Desc:   msg,
				}) {
					return
				}
			}
		}

		if token.Type == types.ScalarToken && config.CheckMultiLineStrings {
			for p := range checkScalarIndentation(config, token, context) {
				if !yield(p) {
					return
				}
			}
		}

		if isVisible {
			context["cur_line"] = getRealEndLine(token)
			if firstInLine {
				context["cur_line_indent"] = token.StartMark.Column
			}
		}

		curLineIndent, _ := context["cur_line_indent"].(int)

		if token.Type == types.BlockMappingStartToken {
			indent := token.StartMark.Column
			context["stack"] = append(stack, &Parent{Type: B_MAP, Indent: indent})

		} else if token.Type == types.FlowMappingStartToken {
			var indent int
			if next.StartMark.Line == token.StartMark.Line {
				indent = next.StartMark.Column
			} else {
				indent = detectIndent(curLineIndent, next)
			}
			context["stack"] = append(stack, &Parent{Type: F_MAP, Indent: indent, LineIndent: curLineIndent})

		} else if token.Type == types.BlockSequenceStartToken {
			indent := token.StartMark.Column
			context["stack"] = append(stack, &Parent{Type: B_SEQ, Indent: indent})

		} else if token.Type == types.BlockEntryToken &&
			!(next != nil && (next.Type == types.BlockEntryToken || next.Type == types.BlockEndToken)) {

			stack = context["stack"].([]*Parent)
			if stack[len(stack)-1].Type != B_SEQ {
				context["stack"] = append(stack, &Parent{Type: B_SEQ, Indent: token.StartMark.Column, ImplicitBlockSeq: true})
				stack = context["stack"].([]*Parent)
			}

			var indent int
			if next.StartMark.Line == token.EndMark.Line {
				indent = next.StartMark.Column
			} else if next.StartMark.Column == token.StartMark.Column {
				indent = next.StartMark.Column
			} else {
				indent = detectIndent(token.StartMark.Column, next)
			}
			context["stack"] = append(stack, &Parent{Type: B_ENT, Indent: indent})

		} else if token.Type == types.FlowSequenceStartToken {
			var indent int
			if next.StartMark.Line == token.StartMark.Line {
				indent = next.StartMark.Column
			} else {
				indent = detectIndent(curLineIndent, next)
			}
			context["stack"] = append(stack, &Parent{Type: F_SEQ, Indent: indent, LineIndent: curLineIndent})

		} else if token.Type == types.KeyToken {
			stack = context["stack"].([]*Parent)
			indent := stack[len(stack)-1].Indent
			p := &Parent{Type: KEY, Indent: indent, ExplicitKey: common.IsExplicitKey(token)}
			context["stack"] = append(stack, p)

		} else if token.Type == types.ValueToken {
			stack = context["stack"].([]*Parent)

			if next.Type == types.AnchorToken || next.Type == types.TagToken {
				if next.StartMark.Line == prev.StartMark.Line &&
					nextNext != nil &&
					next.StartMark.Line < nextNext.StartMark.Line {
					next = nextNext
				}
			}

			if next.Type != types.BlockEndToken &&
				next.Type != types.FlowMappingEndToken &&
				next.Type != types.FlowSequenceEndToken &&
				next.Type != types.KeyToken {

				var indent int
				if stack[len(stack)-1].ExplicitKey {
					indent = detectIndent(stack[len(stack)-1].Indent, next)
				} else if next.StartMark.Line == prev.StartMark.Line {
					indent = next.StartMark.Column
				} else if next.Type == types.BlockSequenceStartToken || next.Type == types.BlockEntryToken {
					indentSequences := context["indent-sequences"]

					if isBool, ok := indentSequences.(bool); ok && !isBool {
						indent = stack[len(stack)-1].Indent
					} else if isBool, ok := indentSequences.(bool); ok && isBool {
						spacesVal := context["spaces"]
						isSpacesConsistent := false
						if str, ok := spacesVal.(string); ok && str == "consistent" {
							isSpacesConsistent = true
						}

						if isSpacesConsistent && (next.StartMark.Column-stack[len(stack)-1].Indent == 0) {
							indent = -1
						} else {
							indent = detectIndent(stack[len(stack)-1].Indent, next)
						}
					} else {
						if next.StartMark.Column == stack[len(stack)-1].Indent {
							if str, ok := indentSequences.(string); ok && str == "consistent" {
								context["indent-sequences"] = false
							}
							indent = stack[len(stack)-1].Indent
						} else {
							if str, ok := indentSequences.(string); ok && str == "consistent" {
								context["indent-sequences"] = true
							}
							indent = detectIndent(stack[len(stack)-1].Indent, next)
						}
					}
				} else {
					indent = detectIndent(stack[len(stack)-1].Indent, next)
				}
				context["stack"] = append(stack, &Parent{Type: VAL, Indent: indent})
			}
		}

		consumedCurrentToken := false
		for {
			stack = context["stack"].([]*Parent)
			if len(stack) == 0 {
				break
			}
			top := stack[len(stack)-1]

			if top.Type == F_SEQ && token.Type == types.FlowSequenceEndToken && !consumedCurrentToken {
				context["stack"] = stack[:len(stack)-1]
				consumedCurrentToken = true
			} else if top.Type == F_MAP && token.Type == types.FlowMappingEndToken && !consumedCurrentToken {
				context["stack"] = stack[:len(stack)-1]
				consumedCurrentToken = true
			} else if (top.Type == B_MAP || top.Type == B_SEQ) &&
				token.Type == types.BlockEndToken &&
				!top.ImplicitBlockSeq &&
				!consumedCurrentToken {
				context["stack"] = stack[:len(stack)-1]
				consumedCurrentToken = true
			} else if top.Type == B_ENT &&
				token.Type != types.BlockEntryToken &&
				len(stack) >= 2 && stack[len(stack)-2].ImplicitBlockSeq &&
				token.Type != types.AnchorToken && token.Type != types.TagToken &&
				(next != nil && next.Type != types.BlockEntryToken) {
				context["stack"] = stack[:len(stack)-2]
			} else if top.Type == B_ENT &&
				(next != nil && (next.Type == types.BlockEntryToken || next.Type == types.BlockEndToken)) {
				context["stack"] = stack[:len(stack)-1]
			} else if top.Type == VAL &&
				token.Type != types.ValueToken &&
				token.Type != types.AnchorToken && token.Type != types.TagToken {
				context["stack"] = stack[:len(stack)-2]
			} else if top.Type == KEY &&
				(next != nil && (next.Type == types.BlockEndToken ||
					next.Type == types.FlowMappingEndToken ||
					next.Type == types.FlowSequenceEndToken ||
					next.Type == types.KeyToken)) {
				context["stack"] = stack[:len(stack)-1]
			} else {
				break
			}
		}
	}
}
