package jfather

// Option configures optional parsing behaviour. The dialect options
// (AllowComments, AllowTrailingCommas, AllowUnescapedControlChars) are all
// off by default, so Unmarshal stays strict JSON unless a caller opts in.
type Option func(*parser)

// AllowComments permits JavaScript-style comments — // to end of line and
// /* ... */ block comments — anywhere whitespace is allowed. Standard
// JSON forbids comments, but several dialects permit them (JSONC,
// tsconfig, and ARM/Bicep deployment templates).
func AllowComments() Option {
	return func(p *parser) { p.allowComments = true }
}

// AllowTrailingCommas permits a single trailing comma before a closing
// ']' or '}'. Standard JSON forbids it.
func AllowTrailingCommas() Option {
	return func(p *parser) { p.allowTrailingCommas = true }
}

// MaxDepth overrides the maximum nesting depth of arrays and objects,
// which defaults to DefaultMaxDepth. Input nested deeper than the limit is
// rejected with an error rather than parsed, protecting the parser — and
// callers that recursively walk the resulting tree — from stack exhaustion
// on maliciously deep input. Values less than 1 are ignored, leaving the
// default in place; the limit cannot be disabled.
func MaxDepth(depth int) Option {
	return func(p *parser) {
		if depth >= 1 {
			p.maxDepth = depth
		}
	}
}

// AllowUnescapedControlChars permits literal control characters (U+0000
// through U+001F), such as raw newlines and tabs, to appear inside string
// values. Standard JSON requires these to be escaped, but some producers
// emit them literally — for example Azure ARM/Bicep templates routinely
// break a long expression across several lines inside a single string
// value. This mirrors Python's json.loads(..., strict=False).
func AllowUnescapedControlChars() Option {
	return func(p *parser) { p.allowUnescapedControls = true }
}
