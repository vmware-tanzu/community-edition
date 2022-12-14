# Oracle CPI & CSI

> This package provides cloud provider interface using oracle-cpi.

For more information, see the [GitHub page](https://github.com/oracle/oci-cloud-controller-manager) of Oracle CPI.

## Configuration

The following configuration values can be set to customize the Oracle CPI & CSI installation.

### Oracle CPI & CSI Configuration

| Value                  | Required/Optional | Description                                                                         |
|------------------------|-------------------|-------------------------------------------------------------------------------------|
| `compartment`          | Required          | compartment configures Compartment within which the cluster resides.                |
| `vcn`                  | Required          | vcn configures the Virtual Cloud Network (VCN) within which the cluster resides.    |
| `loadBalancer.subnet1` | Required          | subnet1 configures one of two subnets to which load balancers will be added..       |
| `loadBalancer.subnet2` | Required          | subnet2 configures the second of two subnets to which load balancers will be added. |
