# Connect to Cluster Nodes with SSH

You can use SSH to connect to nodes in management clusters and workload clusters. The SSH key pair that you created when you deployed the management cluster must be available on the machine on which you run the SSH command.

The SSH key used to access clusters are associated with the following Linux users:

- vSphere management cluster and workload cluster nodes running on both Photon OS and Ubuntu: `capv`
- Amazon EC2 bastion nodes: `ubuntu`
- Amazon EC2 management cluster and workload cluster nodes running on Ubuntu: `ubuntu`
- Amazon EC2 management cluster and workload cluster nodes running on Amazon Linux: `ec2-user`
- Azure management cluster and workload cluster nodes (always Ubuntu): `capi`

To connect to a node by using SSH, run one of the following commands from a machine containing the SSH key:

- vSphere nodes: `ssh capv@<em>node_address`
- Amazon EC2 bastion nodes, management cluster, and workload nodes on Ubuntu: `ssh ubuntu@node_address`
- Amazon EC2 management cluster and workload nodes running on Amazon Linux: `ssh ec2-user@node_address`
- Azure nodes: `ssh capi@<em>node_address`

Each cluster host contains the public key. The hosts will not accept an SSH password (`PasswordAuthentication` is set to `no`).
