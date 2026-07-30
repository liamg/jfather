package jfather

// Option configures optional, non-standard parsing behaviour. All options
// are off by default, so Unmarshal stays strict JSON unless a caller opts
// in.
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
