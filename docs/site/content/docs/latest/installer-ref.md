# Installer Interface Options {#ui-options}

By default, `tanzu management-cluster create --ui` opens the installer interface locally, at [http://127.0.0.1:8080](http://127.0.0.1:8080) in your default browser. You can use the `--browser` and `--bind` options to control where the installer interface runs:

* `--browser` specifies the local browser to open the interface in. Supported values are `chrome`, `firefox`, `safari`, `ie`, `edge`, or `none`. Use `none` with `--bind` to run the interface on a different machine.
* `--bind` specifies the IP address and port to serve the interface from. For example, if another process is already using [http://127.0.0.1:8080](http://127.0.0.1:8080), use `--bind` to serve the interface from a different local port.

Example:

```sh
tanzu management-cluster create --ui --bind 192.168.1.87:5555 --browser none
```

**Warning**: The `tanzu management-cluster create` command takes time to complete.
While `tanzu management-cluster create` is running, do not run additional invocations of `tanzu management-cluster create` on the same bootstrap machine to deploy multiple management clusters, change context, or edit `~/.kube-tkg/config`.

**Warning**: Serving the installer interface from a non-default IP address and port could expose the `tanzu` CLI to a potential security risk while the interface is running. VMware recommends passing in to the `--bind` option an IP and port on a secure network.

* For information about the configurations of the different sizes of node instances, for example `t3.large`, or `t3.xlarge`, see [Amazon EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/).
* For information about when to create a Virtual Private Cloud (VPC) and when to reuse an existing VPC, see [Resource Usage in Your Amazon Web Services Account](aws.md#aws-resources).
