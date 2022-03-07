#!/bin/bash
# Original file from the Open Policy Agent project: https://github.com/open-policy-agent/opa/blob/master/docs/website/scripts/load-docs.sh

set -xe

ORIGINAL_COMMIT=$(git symbolic-ref -q --short HEAD || git name-rev --name-only HEAD)
# If no name can be found "git name-rev" returns
# "undefined", in which case we'll just use the
# current commit ID.
if [[ "${ORIGINAL_COMMIT}" == "undefined" ]]; then
    ORIGINAL_COMMIT=$(git rev-parse HEAD)
fi

ROOT_DIR=$(git rev-parse --show-toplevel)
#RELEASES_YAML_FILE=${ROOT_DIR}/data/releases.yaml
GIT_VERSION=$(git --version)

# Look at the git tags and generate a list of releases
# that we want to show docs for.
if [[ -z ${OFFLINE} ]]; then
    git fetch ${REPOSITORY_URL:-https://github.com/vmware-tanzu/community-edition.git}
fi
ALL_RELEASES=$(git ls-remote https://github.com/vmware-tanzu/community-edition.git/ | grep release | awk -F/ '{print $3}' | sort -r -V)
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

    SEMVER_REGEX='[^0-9]*\([0-9]*\)[.]\([0-9]*\)[.]\([0-9]*\)\([0-9A-Za-z-]*\)'
    CUR_MAJOR_VER=$(echo ${CUR_SEM_VER} | sed -e "s#${SEMVER_REGEX}#\1#")
    CUR_MINOR_VER=$(echo ${CUR_SEM_VER} | sed -e "s#${SEMVER_REGEX}#\2#")
    CUR_PATCH_VER=$(echo ${CUR_SEM_VER} | sed -e "s#${SEMVER_REGEX}#\3#")

    # ignore versions from before we used this static site generator
    #if [[ (${CUR_MAJOR_VER} -lt 1) || \
    #        (${CUR_MAJOR_VER} -le 1 && ${CUR_MINOR_VER} -lt 10) ]]; then
    #    continue
    #fi

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
STASH_TOKEN=$(shuf --random-source=/dev/urandom -i 100000-999999 -n 1)
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
rm -rf ${ROOT_DIR}/generated/*

#echo "Removing data/releases.yaml file"
#rm -f ${RELEASES_YAML_FILE}

#mkdir -p $(dirname ${RELEASES_YAML_FILE})

#echo 'Adding "latest" version to releases.yaml'
#echo "- latest" > ${RELEASES_YAML_FILE}

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
        echo "here goes main"
        git fetch https://github.com/vmware-tanzu/community-edition.git ${release}:main-dev
        git checkout main-dev
    fi
    errc=$?
    set -e

    # only add the version to the releases.yaml data file
    # if we were able to check out the version, otherwise skip it..
#    if [[ "${errc}" == "0" ]]; then
#        echo "Adding ${release} to releases.yaml"
#        echo "- ${release}" >> ${RELEASES_YAML_FILE}
#    else
#        echo "WARNING: Failed to check out version ${version}!!"
#    fi
    if [[ ${release} != "main" ]]; then
        release=$(echo $release | awk -F- '{print $2}')
        version_docs_dir=${ROOT_DIR}/docs/site/generated/docs/${release}
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

# Create the "edge" version from current working tree
#echo 'Adding "edge" to releases.yaml'
#echo "- edge" >> ${RELEASES_YAML_FILE}

# Link instead of copy so we don't need to re-generate each time.
# Use a relative link so it works in a container more easily.
#ln -s ${ROOT_DIR}/docs ${ROOT_DIR}/content/docs/edge

# Create a "latest" version from the latest semver found
#ln -s ${ROOT_DIR}/docs/website/generated/docs/${RELEASES[0]} ${ROOT_DIR}/docs/website/generated/docs/latest