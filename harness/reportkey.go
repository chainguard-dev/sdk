/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package harness

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
)

// ReportKeyBits is the RSA modulus size required for report recipient keys.
const ReportKeyBits = 3072

// ParseReportPublicKey decodes and checks the DER SubjectPublicKeyInfo a client
// sends as CreateHunt's report_public_key. The key must have a ReportKeyBits
// modulus and support RSA-OAEP-SHA256 encryption. Validation errors identify the
// failed check so intake can reject the key before recording an immutable key
// that cannot seal a report.
func ParseReportPublicKey(der []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("report public key is not a DER SubjectPublicKeyInfo: %w", err)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("report public key must be RSA, got %T", pub)
	}
	if bits := rsaPub.N.BitLen(); bits != ReportKeyBits {
		return nil, fmt.Errorf("report public key must be %d bits, got %d", ReportKeyBits, bits)
	}
	// Parsing does not check whether RSA-OAEP accepts the modulus and exponent.
	if _, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, nil, nil); err != nil {
		return nil, fmt.Errorf("report public key cannot encrypt with RSA-OAEP-SHA256: %w", err)
	}
	return rsaPub, nil
}
