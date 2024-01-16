/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

var ErrUnsafeReadme = fmt.Errorf("readme contained unsafe html content")

// ValidateReadme validates the contents of a Markdown README.md file.
// If the contents are invalid, a string will be returned containing the
// diff of what the Markdown would look like as HTML if properly sanitized.
func ValidateReadme(readme string) (string, error) {
	unsafe := readmeToHTML(readme)
	safe := sanitizeHTML(unsafe)
	// After converting the Markdown to HTML,
	// make sure there is no diff after sanitizing it.
	// Unescape any encoded HTML tags for proper comparison.
	if diff := cmp.Diff(unscapeHTML(unsafe), unscapeHTML(safe)); diff != "" {
		return diff, ErrUnsafeReadme
	}
	return "", nil
}

func readmeToHTML(rawMarkdown string) string {
	s := string(blackfriday.Run([]byte(rawMarkdown)))
	// Fix issue where single tags get extra space on conversion (e.g. "<hr />")
	s = strings.ReplaceAll(s, " />", "/>")
	return s
}

var bluemondayPolicy = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	// Allow fenced code block classes
	p = p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	// Allow links without ref="nofollow" which are not set automatically on links on markdown conversion
	p = p.RequireNoFollowOnLinks(false)
	// Allow custom height and width on images
	p = p.AllowAttrs("width", "height").OnElements("img")
	// Allow HTML comments
	p.AllowComments()
	return p
}()

func sanitizeHTML(unsafeHTML string) string {
	return bluemondayPolicy.Sanitize(unsafeHTML)
}

func unscapeHTML(safeHTML string) string {
	return html.UnescapeString(safeHTML)
}
