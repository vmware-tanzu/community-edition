# Tanzu Community Edition Release Process

## Summary

This assumes that you are familiar with this [Release Process](https://github.com/vmware-tanzu/community-edition/blob/main/docs/developer/release-process.md) document.

This assumes that you are a `Maintainer` for Tanzu Community Edition (TCE) and that you have some authority (hence the access) to create, manage, etc the release process. This document covers key points in the release process.

## Prerequisites

This document assumes that you follow the typical workflows in dealing with GitHub. Meaning:

- you *always* are working on *your own fork* of the repo and not directly on the clone of the repo itself
- you have set your fork's `upstream` to point to the TCE repo
- your git is configured such thar the TCE repo (or globally in the example here) is set to `git config --global pull.rebase true`
- you have the [GitHub CLI](https://github.com/cli/cli) installed locally

## Approaching a Release

As you start to wind down your development cycle, your going to need to do 2 major things:

- Branch off from main and create a release branch
- Create subsequent `RC` and `GA` releases

### Branching from `main`

Your release is beginning to stablize and you don't want to include any more "unstable" code into the release. You will want to create a release branch (ie `release-0.12` for the `v0.12.0` release) to limit the amount of change in the release and also allow individual contributors to begin creating PRs to merge into the next release. A great time to do this would be when you are creating one (or ideally the first release) candidate (`RC`). It's kind of a signal that we are stabilizing and close to `done`.

Thankfully, there is automation that exists which will handle this entire process for you. To create your `RC`, branch and set all appropriate tags, update all automation marker files, and etc you just need to execute the following GitHub action:

```bash
gh workflow run --repo https://github.com/vmware-tanzu/community-edition/ \
  --ref main \
  release-tag-and-branch.yaml \
  --raw-field release_version=<input version of RC>
```

It's pretty simple! If this is your first `RC`, then that command would look like:

```bash
gh workflow run --repo https://github.com/vmware-tanzu/community-edition/ \
  --ref main \
  release-tag-and-branch.yaml \
  --raw-field release_version=v0.12.0-rc.1
```

This command then would tag `main` with `v0.12.0-rc.1`, create a `release-0.12` branch, creates `dev` tag in case there are any PRs later on, and sets all the marker files needed for the final `GA` (and beyond) release.

The `main` branch is then setup such that the `dev` tag is created, and sets all the marker files so `main` is the development branch for the `v0.13.0` release.

### Creating Additional RCs and the Final GA Release

Ideally, create subsequent `RC` builds should always be done on the release branch for that version (ie `v0.13.0-rc.2` is only ever created on the `release-0.13` branch in incrementing order). The `GA` release should **ALWAYS** be created from the release branch.

To create a `RC` or `GA` release, you first want to make sure your release branch on your fork is up to date. You can do this by running the following commands:

```bash
git checkout <release branch>
git pull --rebase upstream <release branch>
git push
git fetch upstream --tags
git push origin --tags
```

For example, the `release-0.12` that would look like this:

```bash
git checkout release-0.12
git pull --rebase upstream release-0.12
git push
git fetch upstream --tags
git push origin --tags
```

Then you want to make sure you are dealing with the very tip (ie the latest) of the release branch. Some of this will be repetative to what you just ran, but that's ok! You can do this by fetching everything upstream, rebasing, and resetting the branch by running:

```bash
git fetch upstream
git checkout release-0.12
git reset --hard upstream/release-0.12
```

Finally, to kick off a release, just tag the repo with your release version. If we were cutting `RC2` it might look like this:

```bash
git tag -m v0.12.0-rc.2 v0.12.0-rc.2
git push upstream v0.12.0-rc.2
```

And if we were ready to cut `GA` after cutting `RC2`, you would start at the beginning of this section *Creating Additional RCs and the Final GA Release* and go through all of those previous steps and then instead of tagging for `v0.12.0-rc.2` you would tag for the `GA` version `v0.12.0` like this:

```bash
git tag -m v0.12.0 v0.12.0
git push upstream v0.12.0
```

That's it! For any release, grab the release notes from the provided place and update them.

**IMPORTANT:** If this is the `GA` release, update the release notes but **DO NOT** publish the release until release day. Resist the urge to push the publish button. Good luck!
