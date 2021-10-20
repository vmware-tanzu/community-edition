# Tanzu Community Edition Documentation

This directory contains documentation for this repo. Please refer to the
[style guide](site/content/docs/latest/contribute/style-guide.md) when writing
or editing documents.

## Requirements

This site uses [Hugo](https://github.com/gohugoio/hugo) for rendering. It is
recommended you run hugo locally to validate your changes render properly.

You may also run the `make mdlint` target from the root of this repo to run
markdown linting to catch any syntax and formatting errors.

### Local Hugo Rendering

Hugo is available for many platforms. It can be installed using:

* Linux: Most native package managers
* macOS: `brew install hugo`
* Windows: `choco install hugo-extended -confirm`

Once installed, you may run the following from the `docs/site` directory
to access a rendered view of the documentation:

```bash
hugo server --disableFastRender
```

Access the site at [http://localhost:1313](http://localhost:1313). Press
`Ctrl-C` when done viewing.
