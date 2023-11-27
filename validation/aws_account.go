/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"regexp"
)

const awsAccountPattern = `^[0-9]{12}$`

var (
	awsAccountPatternCompiled = regexp.MustCompile(awsAccountPattern)

	// ErrInvalidAWSAccount describes invalid AWS Account IDs by sharing the
	// regular expression they must match
	ErrInvalidAWSAccount = fmt.Errorf("AWS account ID must match %q", awsAccountPattern)
)

// ValidateAWSAccount checks an AWS account id is valid. AWS Accounts must be a
// 12 digit number.
func ValidateAWSAccount(id string) error {
	if awsAccountPatternCompiled.MatchString(id) {
		return nil
	}
	return ErrInvalidAWSAccount
}
