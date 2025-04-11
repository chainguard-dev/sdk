#!/usr/bin/env bash

# Copyright 2023 Chainguard, Inc.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

VERSION=v1.64.5
URL_BASE=https://raw.githubusercontent.com/golangci/golangci-lint
URL=$URL_BASE/$VERSION/install.sh
# If you update the version above you might need to update the checksum
# if it does not match. We say might because in the past the install script
# has been unchanged even if there is a new verion of golangci-lint.
# To obtain the checksum, download the install script and run the following
# command:
# > sha256sum <path-to-install-script>
INSTALL_CHECKSUM=9e99d38f3213411a1b6175e5b535c72e37c7ed42ccf251d331385a3f97b695e7

if [[ ! -f .golangci.yml ]]; then
    echo 'ERROR: missing .golangci.yml in repo root' >&2
    exit 1
fi

if ! command -v golangci-lint; then
    INSTALL_SCRIPT=$(mktemp -d)/install.sh
    curl -sfL $URL >$INSTALL_SCRIPT
    if echo "${INSTALL_CHECKSUM} $INSTALL_SCRIPT" | sha256sum --check; then
        chmod 755 $INSTALL_SCRIPT
        $INSTALL_SCRIPT -b /tmp $VERSION
        export PATH=$PATH:/tmp
        pwd
    else
        echo 'ERROR: install script sha256 checksum invalid' >&2
        exit 1
    fi
fi

golangci-lint version

error=0
while read -r i; do
  echo "Checking golangci-lint for $i"
  pushd "$i"
  golangci-lint run ./... || error=1
  popd
done <<< "$(find . -name go.mod -exec dirname {} \;)"

exit $error
