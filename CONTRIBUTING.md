# Tanzu Community Edition

First of all, thank you for investing your time in contributing to Tanzu Community Edition (TCE).
These guidelines will help you get started. Please note that we require a
[Contributor License Agreement](#contributor-license-agreement) signed in order
to accept most contributions.

## Building from source

This section describes how to build TCE from source. Build and test is driven by our [Makefile](Makefile).
Commands meant to be used directly by developers feature help text. You can see this by running `make help`.

### Fetch the source

```shell
git clone https://github.com/vmware-tanzu/community-edition
```

### Building the CLI and all plugins from source

TCE consists of the `tanzu` CLI and multiple CLI plugins that facilitate functionality from cluster management to authentication.
The CLI and some of its plugins live in different repositories. To build the CLI and all plugins, including those
hosted in the TCE repository, run the following.

```shell
make install-all-tanzu-cli-plugins
```

After build and install, you'll see an output similar to the following.

```text
[COMPLETE] built and installed plugins at /home/josh/.local/share/tanzu-cli/. These plugins will be automatically detected by tanzu CLI.

[COMPLETE] built and installed tanzu CLI at /home/josh/bin/tanzu. Move this binary to a location in your path!
```

As seen in the message above, you can now move `tanzu` from the location it was installed into a location in your path (such as `/usr/local/bin`).
Plugins are automatically installed in the correctly location, so when calling `tanzu`, the plugins functionality is picked up automatically.

### Building only TCE-specific plugins from source

If you already have `tanzu` CLI installed and wish to only compile and install TCE-specific plugins, run the following.

```shell
make install-tce-cli-plugins
```

After build and install, you'll see an output similar to the following.

```text
[COMPLETE] built and installed TCE-specific plugins at /home/josh/.local/share/tanzu-cli/. These plugins will be automatically detected by your tanzu CLI.
```

Now that the TCE-specifc plugins are installed on your system, you will see their command when running `tanzu`.

### Running TCE-specific plugin tests

To run tests on TCE-specific CLI plugins, run the following.

```shell
make test-plugins
```

### Linting

Several linters are in place to ensure conformant Go code and documentation is written.

To run all linters for the entire project, run:

```shell
make lint
```

To lint just Markdown files (including documentation markdown found under the `docs/` directory), run:

```shell
make mdlint
```

## Contribution workflow

This section describes the process for contributing a bug fix or new feature.

### Before you submit a pull request

This project operates according to the _talk, then code_ rule.
If you plan to submit a pull request for anything more than a typo or obvious bug fix, first you _should_ [raise an issue][new-issue] to discuss your proposal, before submitting any code.

Depending on the size of the feature you may be expected to first write a design proposal. Follow the [Proposal Process](https://github.com/vmware-tanzu/community-edition/blob/main/GOVERNANCE.md#proposal-process) documented in TCE's Governance.

### Commit message and PR guidelines

- Have a short subject on the first line and a body.
- Use the imperative mood (ie "If applied, this commit will (subject)" should make sense).
- Put a summary of the main area affected by the commit at the start,
  with a colon as delimiter. For example 'docs:', 'extensions/(extensionname):', 'design:' or something similar.
- Do not merge commits that don't relate to the affected issue (e.g. "Updating from PR comments", etc). Should
  the need to cherrypick a commit or rollback arise, it should be clear what a specific commit's purpose is.
- If the main branch has moved on, you'll need to rebase before we can merge,
  so merging upstream main or rebasing from upstream before opening your
  PR will probably save you some time.

Pull requests *must* include a `Fixes #NNNN` or `Updates #NNNN` comment.
Remember that `Fixes` will close the associated issue, and `Updates` will link the PR to it.

#### Sample commit message

```text
extensions/extenzi: Add quux functions

To implement the quux functions from #xxyyz, we need to
flottivate the crosp, then ensure that the orping is
appred.

Fixes #xxyyz

Signed-off-by: Your Name <you@youremail.com>
```

### Merging commits

Maintainers should prefer to merge pull requests with the [Squash and merge](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-request-merges#squash-and-merge-your-pull-request-commits) option.
This option is preferred for a number of reasons.
First, it causes GitHub to insert the pull request number in the commit subject which makes it easier to track which PR changes landed in.
Second, it gives maintainers an opportunity to edit the commit message to conform to TCE standards and general [good practice](https://chris.beams.io/posts/git-commit/).
Finally, a one-to-one correspondence between pull requests and commits makes it easier to manage reverting changes and increases the reliability of bisecting the tree (since CI runs at a pull request granularity).

At a maintainer's discretion, pull requests with multiple commits can be merged with the [Create a merge commit](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-request-merges) option.
Merging pull requests with multiple commits can make sense in cases where a change involves code generation or mechanical changes that can be cleanly separated from semantic changes.
The maintainer should review commit messages for each commit and make sure that each commit builds and passes tests.

### Pre commit CI

Before a change is submitted it should pass all the pre commit CI jobs.
If there are unrelated test failures the change can be merged so long as a reference to an issue that tracks the test failures is provided.

### Multiple Go Modules

The TCE project is made up of many different, separate Go modules.
Each module has its own `go.mod` and `go.sum` file located in the "root" of that individual, digestible piece of software.
This allows TCE to not have a single dependency graph, but rather, multiple, independent dependency graphs.
Experimental plugins or packages can pull in different versions of the same library without creating conflicts and collisions.
For this reason, contributors are discouraged from creating interlinking dependencies between the various Go modules.
If you need a "shared" library, please open an issue to discuss with the community.

This has several development implications:

- When working with a piece of code, to allow your editor to have auto-complete and analysis capabilities, open the directory in your editor that contains the `go.mod` file.
  Your editor then sees this as the "root" of the code you're working on so that various tools like the gopls language server will work correctly.
  - Example: If I wanted to work on `cli/cmd/plugin/my-plugin/command.go`, in order to enable Go editor features, I'd open `cli/cmd/plugin/my-plugin` as the top level directory in my editor.
- When adding automation or testing, ensure that your scripts have entered the right directory to execute the right command. For example, because there is no `go.mod` at the top level directory, `go` commands won't work. You must first enter the appropriate directory.

For more information, see the `cli/cmd/README.md` file.

#### Nested Makefiles

It is expected that each individual go module in the TCE repo have its own Makefile.
This enables individual package and plugin authors to have full control over their development operations
without having to modify the top level Makefile.

However, to support discoverability and maintain high level operations,
it _is_ expected that each Makefile provide the following targets:

- `make`: Displays a help message with all poosible make targets
- `make test`: invokes unit tests
- `make e2e-test`: invokes an E2E testing suite
- `make lint`: invokes linting protocols for the individual module. For example, in a Go project, it should call Golangci-lint.
- `make get-deps`: gets the necessary dependencies for running, testing, and building.  Typically is `go mod download` in Go modules
- `make build`: builds the individual piece of software

Some of these targets may be irrelevant to you and your project.
The top level TCE Makefile still expects these targets to be present,
but it's ok to simply print a message stating the target is being skipped or is not applicable.

Beyond the expected targets listed above, package authors are encouraged to create targets that are useful
and relevenat to their development needs.

Users can call:

```shell
make makefile
```

to generate a makefile to stdout that can be used in your project.
This is a good starting point for new packages and plugins integrating directly into the TCE repository.

## Contributor License Agreement

All contributors to this project must have a signed Contributor License
Agreement (**"CLA"**) on file with us. The CLA grants us the permissions we
need to use and redistribute your contributions as part of the project; you or
your employer retain the copyright to your contribution. Before a PR can pass
all required checks, our CLA action will prompt you to accept the agreement.
Head over to [https://cla.vmware.com/](https://cla.vmware.com/) to see your
current agreement(s) on file or to sign a new one.

We generally only need you (or your employer) to sign our CLA once and once
signed, you should be able to submit contributions to any VMware project.

Note: if you would like to submit an "_obvious fix_" for something like a typo,
formatting issue or spelling mistake, you may not need to sign the CLA.
