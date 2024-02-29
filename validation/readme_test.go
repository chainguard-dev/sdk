/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"
)

const (
	safeReadme = `
<!--testing:start-->
# safe

<img src="logo.png" width="36px" height="36px"/>

[click here](./other.md)
<!--testing:end-->

---

Here is some dangerous code inside a code block:


` + "```go" + `
Hello <STYLE>.XSS{background-image:url("javascript:alert('XSS')");}</STYLE><A CLASS=XSS></A>World
` + "```"

	unsafeReadme1 = `
<!--testing:start-->
# unsafe

<img src="logo.png" width="36px" height="36px"/>

[click here](./other.md)
<!--testing:end-->

---

Here is some dangerous code outside a code block:

Hello <STYLE>.XSS{background-image:url("javascript:alert('XSS')");}</STYLE><A CLASS=XSS></A>World
`

	unsafeReadme2 = `
# unsafe

TAKE THAT!!!!!!!
<script>alert("XSS")</script>

`
)

func TestValidateReadme(t *testing.T) {
	tests := map[string]struct {
		Input  string
		Expect bool
	}{
		"empty":    {"", true},
		"safe":     {safeReadme, true},
		"unsafe 1": {unsafeReadme1, false},
		"unsafe 2": {unsafeReadme2, false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			diff, got := ValidateReadme(tt.Input)
			if (got == nil) != tt.Expect {
				t.Errorf("Expected ValidateReadme(`%s`) to return (err == nil) == %v, but got %v. diff: %s", tt.Input, tt.Expect, got, diff)
			}
		})
	}
}
