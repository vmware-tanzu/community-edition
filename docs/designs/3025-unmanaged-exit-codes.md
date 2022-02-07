# Exit codes for `unmanaged-cluster` bootstrapping

`unmanaged-cluster` bootstrapping should provide meaningful exit codes
to capture success and failure.

This will enable better interroperability with other tools
and CI/CD systems.

* [Proposal for exit codes in unmanaged-cluster](https://github.com/vmware-tanzu/community-edition/issues/3025)

## Summary

Currently, during `unmanaged-cluster` bootstrapping,
whenever the command line tool exits control flow,
it returns `0` back to the parent process.
This means that `unmanaged-cluster` cannot provide meaningful
feedback on error to a parent process
without that parent parsing string output from `unmanaged-cluster`.

Instead, we should provide meaningful codes
that represent major steps during the bootstrapping process.
This does not mean that every possible failure case needs an exit code.
But rather, we should provide a subset of codes that communicate
what step in the bootstrapping process failure occurred.

## Motivation

Without this change, our users won't be able to easily
create automation around bootstrapping `unmanaged-cluster`s.
Our users integrate tools like bash and python scripts or CI/CD systems
into their bootstrapping workflows.

### Goals

* Provide a set of exit codes that meaningfully coordinate with failures during boostrapping
* Capture exit codes for the most common failure states

### Non-Goals/Future Work

* Exit codes that capture every possible failure case
* Recommendation on how to utilize the exit codes

## Proposal

The following is the proposed set of exit codes
and their corresponding meaning

```text
0  - Success!
1  - Provided configurations could not be validated
2  - Could not create local cluster directories
3  - Unable to get TKR BOM
4  - Could not render config
5  - TKR BOM not parseable
6  - Could not resolve kapp controller bundle
7  - Unable to create new cluster
8  - Unable to use existing cluster (if provided)
9  - Could not install kapp controller to cluster
10 - Could not install core package repo to cluster
11 - Could not install additional package repo
12 - Could not install CNI package
13 - Failed to merge kubeconfig and set context
```

These high level codes will exist within the `tanzu` package in the `unmanaged-cluster` go module.
This will enable users of the `tanzu` package to get these by default
and not have to managed exit codes themselves.
It would be an anti-pattern to place the exit codes in deeper packages (like `image` or `kapp`)
since those packages may be used by API consumers outside of the higher level user flow that exists in the `tanzu` package.

### User Stories

#### Story 1

As a user of TCE `unmanaged-cluster`,
when writing an automation script to boostrap clusters,
I want to be able to catch certain failure scenarios in my script
In order to take the appropriate action automatically.

#### Story 2

As a user of TCE `unmanaged-cluster`,
when deploying clusters via CI/CD,
I want to be able to display a slack message given certain failure states
In order to have immediate visibility into the problem.

## Compatibility

N/a - This change should not break previous behavior
and only represents an enhacement.

## Alternatives

No good alternatives exist
unless we expect users to parse output from `stdout` and `stderr`.

## Additional Details

### Test Plan

Since this is ultimately about implementation for each code,
these exit codes should be covered in unit tests within the code
and do not require more extensive e2e testing or validation.

### Graduation Criteria

This is targeting a v0.11.0 release
and will be in step with `unmanaged-cluster` being in alpha state.
