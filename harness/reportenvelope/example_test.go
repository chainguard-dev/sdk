/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package reportenvelope_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"

	"chainguard.dev/sdk/harness/reportenvelope"
)

func Example() {
	// The report DEK and its wrapped forms come from WrapReportKey. The
	// wrapped values here stand in for RSA-OAEP ciphertext.
	dek := make([]byte, 32)
	if _, err := rand.Read(dek); err != nil {
		panic(err)
	}
	keys := reportenvelope.KeyRegion{
		Epoch: 1,
		Slots: []reportenvelope.Slot{
			{Kind: reportenvelope.OrgKEK, KeyVersion: 1, Wrapped: []byte("org-kek wrap")},
			{Kind: reportenvelope.CustomerSubmitKey, Wrapped: []byte("customer wrap")},
		},
	}

	var sealed bytes.Buffer
	w, err := reportenvelope.Seal(&sealed, dek, keys)
	if err != nil {
		panic(err)
	}
	if _, err := io.WriteString(w, "the rendered report"); err != nil {
		panic(err)
	}
	if err := w.Close(); err != nil {
		panic(err)
	}

	// A reader picks its slot, unwraps the DEK from it, and opens the body.
	r := bytes.NewReader(sealed.Bytes())
	h, err := reportenvelope.ReadHeader(r)
	if err != nil {
		panic(err)
	}
	slot, ok := h.Keys.Slot(reportenvelope.CustomerSubmitKey)
	if !ok {
		panic("no customer-submit-key slot")
	}
	fmt.Printf("unwrap %q\n", slot.Wrapped)
	body, err := h.Open(r, dek)
	if err != nil {
		panic(err)
	}
	report, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(report))
	// Output:
	// unwrap "customer wrap"
	// the rendered report
}
