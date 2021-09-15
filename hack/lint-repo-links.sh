#!/usr/bin/env bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# This script is used to crawl a github repository, find any internal VMware links,
# including in issues, prs, comments, and git history, and report then them.
# Not all links reported are required to be removed.
# 
# The links that are searched for include:
# - docs.google.com
# - drive.google.com
# - eng.vmware.com
# - vmware.slack.com
# - miro.com
#
# It uses the gh CLI tool to find issues and PRs
# Ensure you have the gh CLI and have authenticated using `gh auth login`

set -eu

BOLD='\033[1m'
CLEAR='\033[0m'

function warning {
  YELLOW='\033[0;33m'
  local message=${1}

  echo -e "${YELLOW}${BOLD}${message}${CLEAR}"
}

function error {
  RED='\033[0;31m'
  local message=${1}

  echo -e "${RED}${BOLD}${message}${CLEAR}"
}

function success {
  GREEN='\033[0;32m'
  local message=${1}

  echo -e "${GREEN}${BOLD}${message}${CLEAR}"
}

function info {
  BLUE='\033[0;34m'
  local message=${1}

  echo -e "${BLUE}${BOLD}${message}${CLEAR}"
}

function check_dependencies {
  # this script requires git & the gh CLI tool
  warning "--- Checking dependencies ---"

  if ! git --version > /dev/null 2>&1; then
    error "git binary test failed. Git required"
    exit 1
  fi

  if ! gh --version > /dev/null 2>&1; then
    error "gh binary test failed. GitHub CLI tool required"
    exit 1
  fi

  if ! jq --version > /dev/null 2>&1; then
    error "jq binary test failed. jq json required"
    exit 1
  fi

  if ! rg --version > /dev/null 2>&1; then
    error "rg binary test failed. Ripgrep utility required"
    exit 1
  fi

    # Inform user authenticated to github is required
    info "🛑 - Note: ensure you've logged into to the gh CLI using \`gh auth login\` or have GITHUB_TOKEN set in your environment"

  success "--- Dependencies verified --- \n"
}

success "-----------------------------------------------"
success "--- 🚀 GitHub repo link checker starting ✨ ---"
success "-----------------------------------------------"

check_dependencies

warning "--- Checking git commit history for internal links ---\n"

for commit in $(git rev-list main); do
    commit_message=$(git log --format=%B -n 1 "$commit")

    if echo "$commit_message" | rg -i -w "docs.google.com|drive.google.com|eng.vmware.com|vmware.slack.com|miro.com"; then
        error "Found commit containing an internal link: $commit"
        info "$(git log --format=%B -n 1 "$commit")"
        echo
    fi
done

# set non-interactive mode for github CLI to prevent weird outputs
gh config set prompt disabled

# lint github issue titles using jq queries
warning "Checking Github issues for titles containing jira links:"
gh issue list -L 2000 --state all --json title --jq '.[] | select(.title| test("jira"; "i"))' | jq

warning "\nChecking Github issues for titles containing vmware-tanzu-private links:"
gh issue list -L 2000 --state all --json title --jq '.[] | select(.title| test("vmware-tanzu-private"; "i"))' | jq

warning "\nChecking Github issues for titles containing slack links:"
gh issue list -L 2000 --state all --json title --jq '.[] | select(.title| test("slack"; "i"))' | jq

warning "\nChecking Github issues for titles containing miro links:"
gh issue list -L 2000 --state all --json title --jq '.[] | select(.title| test("miro"; "i"))' | jq

info "\nNote: The following checks may take over 1 hour as each issue and PR must be analized\n"
info "\nWarning: Running this script too frequently will result in you being rate-limited by GitHub. If this occurs, go get a coffee and try again in an hour!"

warning "\n----- Checking Github issues ---\n"

# get valid issue numbers (since pull requests and "issues" are the same to the GitHub API)
for i in $(gh issue list -L 1740 --state all --json number --jq '.[] | .number'); do
    printf "."

    # lint github issue contents body
    if gh issue view "$i" --json body --jq '.[]' | rg -i -w "docs.google.com|drive.google.com|eng.vmware.com|vmware.slack.com|miro.com"; then
        error "Found internal link in issue body #$i\n"
    fi

    # Lint the issue comments
    if gh issue view "$i" --json comments --jq '.[] | map(.body)' | rg -i -w "docs.google.com|drive.google.com|eng.vmware.com|vmware.slack.com|miro.com"; then
        error "Found internal link in issue comment #$i\n"
    fi
done

# lint github PR contents

warning "\n----- Checking Github PR ---"

# get valid PR numbers (since pull requests and "issues" are the same to the GitHub API)
for i in $(gh pr list -L 1740 --state all --json number --jq '.[] | .number'); do
    printf "."

    # lint github issue contents body
    if gh pr view "$i" --json body --jq '.[]' | rg -i -w "docs.google.com|drive.google.com|eng.vmware.com|vmware.slack.com|miro.com"; then
        error "\nFound internal link in pr body #$i\n"
    fi

    # Lint the pr comments
    if gh pr view "$i" --json comments --jq '.[] | map(.body)' | rg -i -w "docs.google.com|drive.google.com|eng.vmware.com|vmware.slack.com|miro.com"; then
        error "\nFound internal link in pr comment #$i\n"
    fi
done

success "\n--------------"
success "--- Done!! ---"
success "--------------"

