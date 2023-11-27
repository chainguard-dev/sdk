/*
Copyright 2021 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uidp

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// NOTE: We are using UID over UUID for 2 reasons:
// 1. UUIDs are very large.
// 2. UUIDs are not very random.
// Ref: https://neilmadden.blog/2018/08/30/moving-away-from-uuids/

// UID will is used for the primary key for items that must be globally unique.
//   - A UID is 20 bytes of random bytes, URL safe hex encoded.
type UID string

// SUID is be used to form the primary key for items that must be unique within some scoping (non-global).
//   - A SUID is 8 random bytes, URL safe hex encoded.
type SUID string

// UIDP is be used to denote the fully-qualified path for scoped keys.
//   - A UIDP will consist of '/' delimited SUID segments with a UID root, following POSIX directory semantics.
//   - The "basename" SUID is our key within the scoping of the "dirname" UIDP.
type UIDP string

func NewUID() UID {
	token := make([]byte, 20)
	_, _ = rand.Read(token)
	return UID(hex.EncodeToString(token))
}

func NewSUID() SUID {
	token := make([]byte, 8)
	_, _ = rand.Read(token)
	return SUID(hex.EncodeToString(token))
}

func NewUIDP(path UIDP) UIDP {
	if len(path) == 0 {
		return UIDP(NewUID())
	}
	return UIDP(fmt.Sprintf("%s/%s", path, NewSUID()))
}

func (u UIDP) NewChild() UIDP {
	return NewUIDP(u)
}

func (u UID) String() string {
	return string(u)
}

func (u SUID) String() string {
	return string(u)
}

func (u UIDP) String() string {
	return string(u)
}
