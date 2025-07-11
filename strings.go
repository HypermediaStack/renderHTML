package renderHTML

import (
	"fmt"
	"strings"
)

var htmlEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&#34;",
	"'", "&#39;",
)

// #region EscapeString

type escapeStringEntity struct {
	content string
}

// String returns HTML text of the current element.
func (p *escapeStringEntity) String() string {
	return htmlEscaper.Replace(p.content)
}

// EscapeString escapes special characters like "<" to become "&lt;". It
// escapes only five such characters: <, >, &, ' and ".
// UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
// always true.
func EscapeString(format string, args ...any) *escapeStringEntity {
	return &escapeStringEntity{fmt.Sprintf(format, args...)}
}

// #region RawString

type rawStringEntity struct {
	content string
}

// String returns HTML text of the current element.
func (p *rawStringEntity) String() string {
	return p.content
}

func rawString(format string, args ...any) *rawStringEntity {
	return &rawStringEntity{fmt.Sprintf(format, args...)}
}

// RawString creates an element that is used to include WYSIWYG elements.
// Especially to include HTML formatted text strings.
func RawString(format string, args ...any) *rawStringEntity {
	return rawString(format, args...)
}
