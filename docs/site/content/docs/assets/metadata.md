In the optional **Metadata** section, provide descriptive information about the cluster:

- **Location**: The geographical location in which the clusters run.
- **Description**: A description of this cluster. The description has a maximum length of 63 characters and must start and end with a letter. It can contain only lower case letters, numbers, and hyphens, with no spaces.
- **Labels**: Key/value pairs to help users identify clusters, for example `release : beta`, `environment : staging`, or `environment : production`. For more information, see [Labels and Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/).<br />
You can click **Add** to apply multiple labels to the clusters.

Any metadata that you specify here applies to the management cluster, standalone clusters, and workload clusters, and can be accessed by using the cluster management tool of your choice.
