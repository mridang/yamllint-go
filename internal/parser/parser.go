package parser

import (
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
	"github.com/mridang/yamllint-go/internal/types"
)

// Parse converts raw YAML content into the linter's internal types.
func Parse(content string) ([]*types.Token, []*types.Comment, []*types.Line, error) {
	lines := parseLines(content)
	rawTokens := lexer.Tokenize(content)

	var linterTokens []*types.Token
	var linterComments []*types.Comment

	// Explicit StreamStartToken (0-based line/col)
	linterTokens = append(linterTokens, &types.Token{
		Type:      types.StreamStartToken,
		StartMark: types.Mark{Index: 0, Line: 0, Column: 0, Buffer: content},
		EndMark:   types.Mark{Index: 0, Line: 0, Column: 0, Buffer: content},
		Encoding:  "utf-8",
	})

	for i := 0; i < len(rawTokens); i++ {
		t := rawTokens[i]

		// 1. Handle Comments
		if t.Type == token.CommentType {
			comment := mapComment(t, content)
			if len(linterTokens) > 0 {
				comment.TokenBefore = linterTokens[len(linterTokens)-1]
			}
			linterComments = append(linterComments, comment)
			continue
		}

		// 2. Ignore Spaces/Invalid
		if t.Type == token.SpaceType || t.Type == token.InvalidType {
			continue
		}

		// 3. Merge Logic for Anchors/Aliases/Tags
		if isIndicator(t.Type) && i+1 < len(rawTokens) {
			next := rawTokens[i+1]
			if isMergeableValue(next.Type) {
				linterTokens = append(linterTokens, mapMergedToken(t, next, content))
				i++ // consumed next
				continue
			}
		}

		// 4. Standard Token Mapping
		linterTokens = append(linterTokens, mapToken(t, content))
	}

	// Explicit StreamEndToken
	endIdx := len(content)
	lastLine := 0
	if len(lines) > 0 {
		lastLine = len(lines) - 1 // 0-based
	}
	linterTokens = append(linterTokens, &types.Token{
		Type:      types.StreamEndToken,
		StartMark: types.Mark{Index: endIdx, Line: lastLine, Column: 0, Buffer: content},
		EndMark:   types.Mark{Index: endIdx, Line: lastLine, Column: 0, Buffer: content},
	})

	// Link comments to TokenAfter
	tokenIdx := 0
	for _, comment := range linterComments {
		for tokenIdx < len(linterTokens) {
			if linterTokens[tokenIdx].StartMark.Index > comment.Pointer {
				comment.TokenAfter = linterTokens[tokenIdx]
				break
			}
			tokenIdx++
		}
	}

	for j := 1; j < len(linterComments); j++ {
		linterComments[j].CommentBefore = linterComments[j-1]
	}

	return linterTokens, linterComments, lines, nil
}

func isIndicator(typ token.Type) bool {
	return typ == token.AnchorType || typ == token.AliasType || typ == token.TagType
}

func isMergeableValue(typ token.Type) bool {
	return typ == token.StringType ||
		typ == token.IntegerType ||
		typ == token.FloatType ||
		typ == token.BoolType
}

func mapMergedToken(indicator, value *token.Token, buffer string) *types.Token {
	// go-yaml Position.Line and Position.Column are 1-based. Convert to 0-based.
	startLine := indicator.Position.Line - 1
	startCol := indicator.Position.Column - 1
	startIdx := indicator.Position.Offset

	endIdx := value.Position.Offset + len(value.Origin)
	endCol := startCol + (endIdx - startIdx)

	return &types.Token{
		Type: mapTokenType(indicator.Type),
		StartMark: types.Mark{
			Index:  startIdx,
			Line:   startLine,
			Column: startCol,
			Buffer: buffer,
		},
		EndMark: types.Mark{
			Index:  endIdx,
			Line:   startLine,
			Column: endCol,
			Buffer: buffer,
		},
		Value: value.Value,
		Style: mapStyle(value.Type),
		Plain: isPlain(value.Type),
	}
}

func mapToken(t *token.Token, buffer string) *types.Token {
	// go-yaml Position.Line and Position.Column are 1-based. Convert to 0-based.
	startLine := t.Position.Line - 1
	startCol := t.Position.Column - 1
	startIdx := t.Position.Offset

	// compute end from raw offsets (exclusive)
	endIdx := startIdx + len(t.Origin)
	endCol := startCol + len(t.Origin)

	val := t.Value
	if t.Type == token.AnchorType || t.Type == token.AliasType {
		if len(val) > 1 && (val[0] == '&' || val[0] == '*') {
			val = val[1:]
		}
	}

	return &types.Token{
		Type: mapTokenType(t.Type),
		StartMark: types.Mark{
			Index:  startIdx,
			Line:   startLine,
			Column: startCol,
			Buffer: buffer,
		},
		EndMark: types.Mark{
			Index:  endIdx,
			Line:   startLine,
			Column: endCol,
			Buffer: buffer,
		},
		Value: val,
		Style: mapStyle(t.Type),
		Plain: isPlain(t.Type),
	}
}

func parseLines(content string) []*types.Line {
	var lines []*types.Line
	start := 0
	lineNo := 1
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			lines = append(lines, &types.Line{LineNo: lineNo, Buffer: content, Start: start, End: i})
			start = i + 1
			lineNo++
		}
	}
	if start <= len(content) {
		lines = append(lines, &types.Line{LineNo: lineNo, Buffer: content, Start: start, End: len(content)})
	}
	return lines
}

func mapComment(t *token.Token, buffer string) *types.Comment {
	// keep aligned with our 0-based token marks
	return &types.Comment{
		LineNo:   t.Position.Line - 1,
		ColumnNo: t.Position.Column - 1,
		Buffer:   buffer,
		Pointer:  t.Position.Offset,
	}
}

func isPlain(typ token.Type) bool {
	switch typ {
	case token.SingleQuoteType, token.DoubleQuoteType, token.LiteralType, token.FoldedType:
		return false
	default:
		return true
	}
}

func mapStyle(typ token.Type) string {
	switch typ {
	case token.SingleQuoteType:
		return "'"
	case token.DoubleQuoteType:
		return "\""
	case token.LiteralType:
		return "|"
	case token.FoldedType:
		return ">"
	default:
		return ""
	}
}

func mapTokenType(typ token.Type) types.TokenType {
	switch typ {
	case token.DocumentHeaderType:
		return types.DocumentStartToken
	case token.DocumentEndType:
		return types.DocumentEndToken
	case token.SequenceEntryType:
		return types.BlockEntryToken
	case token.MappingKeyType:
		return types.KeyToken
	case token.MappingValueType:
		return types.ValueToken
	case token.SequenceStartType:
		return types.FlowSequenceStartToken
	case token.SequenceEndType:
		return types.FlowSequenceEndToken
	case token.MappingStartType:
		return types.FlowMappingStartToken
	case token.MappingEndType:
		return types.FlowMappingEndToken
	case token.CollectEntryType:
		return types.FlowEntryToken
	case token.AliasType:
		return types.AliasToken
	case token.AnchorType:
		return types.AnchorToken
	case token.TagType:
		return types.TagToken
	case token.DirectiveType:
		return types.DirectiveToken
	case token.LiteralType, token.FoldedType, token.SingleQuoteType, token.DoubleQuoteType,
		token.StringType, token.IntegerType, token.FloatType, token.BoolType,
		token.NullType, token.InfinityType, token.NanType,
		token.BinaryIntegerType, token.OctetIntegerType, token.HexIntegerType:
		return types.ScalarToken
	default:
		return types.ScalarToken
	}
}
