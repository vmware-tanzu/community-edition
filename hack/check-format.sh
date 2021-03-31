#!/bin/bash

# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

# Change directories to the parent directory of the one in which this
# script is located.
cd "$(dirname "${BASH_SOURCE[0]}")/.."

# Ensure the temp out file is removed when this program exits.
out="$(mktemp)"
on_exit() {
  [ -z "${out}" ] || [ ! -e "${out}" ] || rm -f "${out}"
}
trap on_exit EXIT

# Run goformat on all the sources.
flags="-e -s -w"
[ -z "${PROW_JOB_ID-}" ] || flags="-d ${flags}"
eval "gofmt ${flags} ./cli/cmd/ ./cli/pkg/" | tee "${out}"

# Check to see if there any suggestions.
goformat_exit_code=0; test -z "$(head -n 1 "${out}")" || goformat_exit_code=1

# Truncate the out file.
rm -f "${out}" && touch "${out}"

# Run goimports on all the sources.
go get golang.org/x/tools/cmd/goimports
cmd=$(go list -f \{\{\.Target\}\} golang.org/x/tools/cmd/goimports)
flags="-e -w"
[ -z "${PROW_JOB_ID-}" ] || flags="-d ${flags}"
eval "${cmd} ${flags} ./cli/cmd/ ./cli/pkg/" | tee "${out}"

# Check to see if there any suggestions.
goimports_exit_code=0; test -z "$(head -n 1 "${out}")" || goimports_exit_code=1

# If running on Prow, exit with a non-zero code if either of the tests failed.
if [ -n "${PROW_JOB_ID-}" ]; then
  [ "${goformat_exit_code}" -eq "0" ] ||  exit "${goformat_exit_code}"
  [ "${goimports_exit_code}" -eq "0" ] || exit "${goimports_exit_code}"
fi
