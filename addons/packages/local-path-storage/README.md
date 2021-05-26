# Local Path Storage

This package provides local path node storage for use on Docker

## Configuration

| Value                                                 | Required/Optional | Description                                                                                                                                                                                                                                                                               |
|-------------------------------------------------------|-------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `namespace`                                           | Required          | The namespace to deploy the local-path-storage pods                                                                                                                                                                                                                                       |

## Usage Examples

A StorageClass is required in order to use PVCs and store data (necessary for services
like Prometheus). The local-path-storage provider enables local CAPD clusters to store data locally.
