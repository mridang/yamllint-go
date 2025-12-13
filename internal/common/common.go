package common

import "github.com/mridang/yamllint-go/internal/types"

// SpacesAfter checks spaces immediately after a token, stopping at the next
// non-space token on the same line.
func SpacesAfter(
	token, _, next *types.Token,
	min, max int,
	minDesc, maxDesc string,
) *types.Problem {
	if next == nil {
		return nil
	}
	if token.EndMark.Line != next.StartMark.Line {
		return nil
	}
	if min == -1 && max == -1 {
		return nil
	}

	spaces := next.StartMark.Column - token.EndMark.Column

	if min != -1 && spaces < min {
		return &types.Problem{
			Line:   token.StartMark.Line + 1,
			Column: token.EndMark.Column + 1,
			Desc:   minDesc,
		}
	}
	if max != -1 && spaces > max {
		return &types.Problem{
			Line:   token.StartMark.Line + 1,
			Column: token.EndMark.Column + 1,
			Desc:   maxDesc,
		}
	}
	return nil
}

// SpacesBefore checks spaces immediately before a token, stopping at the previous
// non-space token on the same line.
func SpacesBefore(
	token, prev, _ *types.Token,
	min, max int,
	minDesc, maxDesc string,
) *types.Problem {
	if prev == nil {
		return nil
	}
	if token.StartMark.Line != prev.EndMark.Line {
		return nil
	}
	if min == -1 && max == -1 {
		return nil
	}

	spaces := token.StartMark.Column - prev.EndMark.Column

	if min != -1 && spaces < min {
		return &types.Problem{
			Line:   token.StartMark.Line + 1,
			Column: token.StartMark.Column + 1,
			Desc:   minDesc,
		}
	}
	if max != -1 && spaces > max {
		return &types.Problem{
			Line:   token.StartMark.Line + 1,
			Column: token.StartMark.Column + 1,
			Desc:   maxDesc,
		}
	}
	return nil
}
