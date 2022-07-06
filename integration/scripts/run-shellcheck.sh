#!/usr/bin/env bash
#
# Copyright (c) 2022 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Simple script to find all bash scripts and run shellcheck on them
#

set -eu

is_bash() {
    [[ $1 == *.sh ]] && return 0
    [[ $1 == */bash-completion/* ]] && return 0
    [[ $(file -b --mime-type "$1") == text/x-shellscript ]] && return 0
    return 1
}

while IFS= read -r -d $'' file; do
    if is_bash "$file"; then
        shellcheck -x -W0 -s bash "$file" || continue
    fi
done < <(find . -type f \! -path "./.git/*" -print0)
