# Metrics

The Tanzu Community Edition project automatically reports metrics
regarding cluster creation, builds, and more.

_*Note*_: At the time of this writing, metrics can only be visualized via the internal VMware WaveFront instance.
In the future, it is intended that the WaveFront dashboard used to monitor Tanzu Community Edition performance
will be made public and viewable by the community.

Metrics are _only_ reported on GitHub Action runs and are injested directly by WaveFront.
This gives the maintainers and community reproducible data points for potential performance regressions.

## Reports on Pull Request runs

### PR Sources

- `tce.github.<PR-number>`: The GitHub PR number that the metric data-point originated from

### PR Metrics

- `tce.pr.management-cluster.create-time`: The time it takes to create a management cluster reported in seconds
- `tce.pr.management-cluster.delete-time`: The time it takes to delete a management cluster reported in seconds

### PR Metadata

- `sha`: The git sha the data point originated from
- `version`: The git tag version the data point originated from
