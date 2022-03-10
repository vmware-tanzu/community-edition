#!/bin/bash
# Original file from the Open Policy Agent project: https://github.com/open-policy-agent/opa/blob/master/docs/website/scripts/load-docs.sh
# Edits based on the Harbor project: https://github.com/goharbor/website/blob/main/load-docs.sh

set -xe

ORIGINAL_COMMIT=$(git symbolic-ref -q --short HEAD || git name-rev --name-only HEAD)
# If no name can be found "git name-rev" returns
# "undefined", in which case we'll just use the
# current commit ID.
if [[ "${ORIGINAL_COMMIT}" == "undefined" ]]; then
    ORIGINAL_COMMIT=$(git rev-parse HEAD)
fi

ROOT_DIR=$(git rev-parse --show-toplevel)
GIT_VERSION=$(git --version)

# Look at the git tags and generate a list of releases
# that we want to show docs for.
if [[ -z ${OFFLINE} ]]; then
    git fetch ${REPOSITORY_URL:-https://github.com/vmware-tanzu/community-edition.git}
fi
ALL_RELEASES=$(git ls-remote https://github.com/vmware-tanzu/community-edition | grep release | awk -F/ '{print $3}' | sort -r -V)
RELEASES=()
PREV_MAJOR_VER="-1"
PREV_MINOR_VER="-1"

for release in ${ALL_RELEASES}; do
    CUR_SEM_VER=${release#"release-"}

    # ignore any release candidate versions, for now if they
    # are the "latest" they'll be documented under "edge"
    if [[ "${CUR_SEM_VER}" == *"rc"* ]]; then
      continue
    fi

    SEMVER_REGEX='[^0-9]*\([0-9]*\)[.]\([0-9]*\)'
    CUR_MAJOR_VER=$(echo ${CUR_SEM_VER} | sed -e "s#${SEMVER_REGEX}#\1#")
    CUR_MINOR_VER=$(echo ${CUR_SEM_VER} | sed -e "s#${SEMVER_REGEX}#\2#")
    #CUR_PATCH_VER=$(echo ${CUR_SEM_VER} | sed -e "s#${SEMVER_REGEX}#\3#")

    # The releases are sorted in order by semver from newest to oldest, and we only want
    # the latest point release for each minor version
    if [[ "${CUR_MAJOR_VER}" != "${PREV_MAJOR_VER}" || \
             ("${CUR_MAJOR_VER}" = "${PREV_MAJOR_VER}" && \
                "${CUR_MINOR_VER}" != "${PREV_MINOR_VER}") ]]; then
        RELEASES+=(${release})
    fi

    PREV_MAJOR_VER=${CUR_MAJOR_VER}
    PREV_MINOR_VER=${CUR_MINOR_VER}
done

echo "Git version: ${GIT_VERSION}"

echo "Saving current workspace state"
STASH_TOKEN=$(od -A n -t d -N 3 /dev/urandom |tr -d ' ')
git stash push --include-untracked -m "${STASH_TOKEN}"

function restore_tree {
    echo "Returning to commit ${ORIGINAL_COMMIT}"
    git checkout ${ORIGINAL_COMMIT}

    # Only pop from the stash if we had stashed something earlier
    if [[ -n "$(git stash list | head -1 | grep ${STASH_TOKEN} || echo '')" ]]; then
        git stash pop
    fi
}

function cleanup {
    EXIT_CODE=$?

    if [[ "${EXIT_CODE}" != "0" ]]; then
        # on errors attempt to restore the starting tree state
        restore_tree

        echo "Error loading docs"
        exit ${EXIT_CODE}
    fi

    echo "Docs loading complete"
}

trap cleanup EXIT

echo "Cleaning generated folder"
rm -rf ${ROOT_DIR}/docs/site/generated/*

# include the main branch as we want to be able to display the latest main
RELEASES+=("main")

for release in "${RELEASES[@]}"; do

    echo "Checking out release ${release}"

    # Don't error if the checkout fails
    set +e
    if [[ ${release} != "main" ]]; then
        git fetch https://github.com/vmware-tanzu/community-edition.git ${release}:${release}-local
        git checkout ${release}-local
    else
        git fetch https://github.com/vmware-tanzu/community-edition.git ${release}:main-dev
        git checkout main-dev
    fi
    errc=$?
    set -e

    if [[ ${release} != "main" ]]; then
        release=$(echo $release | awk -F- '{print $2}')
        version_docs_dir=${ROOT_DIR}/docs/site/generated/docs/v${release}
        else
        version_docs_dir=${ROOT_DIR}/docs/site/generated/docs/${release}
    fi
    mkdir -p ${version_docs_dir}

    echo "Copying doc content from tag ${release}"
    cp -r ${ROOT_DIR}/docs/site/content/docs/* ${version_docs_dir}/

done

# Move generated content to the right place
rm -fr content/docs/*
mkdir -p content/docs
mv generated/docs/* content/docs

# Go back to the original tree state
restore_tree
