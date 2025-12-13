package types

// TokenType mimics the PyYAML class types
type TokenType int

const (
	// Stream Tokens
	StreamStartToken TokenType = iota
	StreamEndToken

	// Document Tokens
	DocumentStartToken
	DocumentEndToken

	// Block Tokens
	BlockSequenceStartToken
	BlockMappingStartToken
	BlockEndToken

	// Flow Tokens
	FlowSequenceStartToken // '['
	FlowSequenceEndToken   // ']'
	FlowMappingStartToken  // '{'
	FlowMappingEndToken    // '}'

	// Data Tokens
	KeyToken        // '?'
	ValueToken      // ':'
	BlockEntryToken // '-'
	FlowEntryToken  // ','

	// Complex Tokens
	AliasToken     // *anchor
	AnchorToken    // &anchor
	TagToken       // !tag
	ScalarToken    // "value", 'value', value
	DirectiveToken // %YAML 1.2
)

// Mark corresponds to the PyYAML Mark object.
// It represents a specific point in the file.
type Mark struct {
	Index  int    // The byte offset in the buffer (Python: pointer)
	Line   int    // The line number (0-based or 1-based depending on parser, usually 0 in logic)
	Column int    // The column number
	Buffer string // The ENTIRE file content
}

// Token corresponds to the objects inside parser.Token.
type Token struct {
	Type      TokenType
	StartMark Mark
	EndMark   Mark

	// Fields specific to Token subclasses
	Value    string // Used by: Alias, Anchor, Tag, Scalar, Directive
	Name     string // Used by: Directive
	Encoding string // Used by: StreamStart
	Plain    bool   // Used by: Scalar
	Style    string // Used by: Scalar
}

// Problem corresponds to the 'LintProblem' class
type Problem struct {
	Line   int
	Column int
	Desc   string
	Rule   string // Optional: to track which rule generated it
}

// Line corresponds to the 'Line' class in yamllint
type Line struct {
	LineNo int
	Buffer string // The full file content buffer
	Start  int    // Start index of the line in the buffer
	End    int    // End index of the line in the buffer
}

// Content returns the substring of the buffer corresponding to this line.
func (l *Line) Content() string {
	if l.Start > l.End || l.Start < 0 || l.End > len(l.Buffer) {
		return ""
	}
	return l.Buffer[l.Start:l.End]
}

// Comment corresponds to the 'Comment' class in yamllint
type Comment struct {
	LineNo        int
	ColumnNo      int
	Buffer        string
	Pointer       int
	TokenBefore   *Token
	TokenAfter    *Token
	CommentBefore *Comment
}

// String returns the content of the comment up to the next newline.
func (c *Comment) String() string {
	if c.Pointer >= len(c.Buffer) {
		return ""
	}
	rest := c.Buffer[c.Pointer:]
	// Find next newline
	for i, char := range rest {
		if char == '\n' || char == 0 {
			return rest[:i]
		}
	}
	return rest
}

// IsInline determines if the comment is inline with code or standalone.
func (c *Comment) IsInline() bool {
	if c.TokenBefore == nil || c.TokenBefore.Type == StreamStartToken {
		return false
	}
	if c.LineNo != c.TokenBefore.EndMark.Line+1 {
		return false
	}
	// Check if the character immediately before the token end is NOT a newline
	idx := c.TokenBefore.EndMark.Index - 1
	if idx >= 0 && idx < len(c.Buffer) {
		if c.Buffer[idx] == '\n' {
			return false
		}
	}
	return true
}
