/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login

import (
	"errors"
	"testing"
)

func TestOpenBrowserErrorAs(t *testing.T) {
	tests := map[string]struct {
		err  error
		want bool
	}{
		"nil": {
			err:  nil,
			want: false,
		},
		"success": {
			err: &OpenBrowserError{
				errors.New("unit test"),
			},
			want: true,
		},
		"failure": {
			err:  errors.New("unit test"),
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var want *OpenBrowserError
			got := errors.As(test.err, &want)
			if got != test.want {
				t.Errorf("As() expected %t, got %t", test.want, got)
			}
			if got {
				t.Log(want.Error())
			}
		})
	}
}
