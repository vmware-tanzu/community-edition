### <a id="ui-options"></a> Installer Interface Options

TESTBy default, `tanzu management-cluster create --ui` opens the installer interface locally, at http://127.0.0.1:8080 in your default browser. You can use the `--browser` and `--bind` options to control where the installer interface runs:


- `--browser` specifies the local browser to open the interface in. Supported values are `chrome`, `firefox`, `safari`, `ie`, `edge`, or `none`. Use `none` with `--bind` to run the interface on a different machine.
- `--bind` specifies the IP address and port to serve the interface from. For example, if another process is already using http://127.0.0.1:8080, use `--bind` to serve the interface from a different local port.

Example:
        ```
        tanzu management-cluster create --ui --bind 192.168.1.87:5555 --browser none
        ```  




<p class="note warning"><strong>Warning</strong>: The <code>tanzu management-cluster create</code> command takes time to complete.
While <code>tanzu management-cluster create</code> is running, do not run additional invocations of <code>tanzu management-cluster create</code> on the same bootstrap machine to deploy multiple management clusters, change context, or edit <code>~/.kube-tkg/config</code>.</p>





  <p class="note warning"><strong>Warning</strong>: Serving the installer interface from a non-default IP address and port could expose the <code>tanzu</code> CLI to a potential security risk while the interface is running. VMware recommends passing in to the <code>--bind</code> option an IP and port on a secure network.</p>

  - For information about the configurations of the different sizes of node instances, for example `t3.large`, or `t3.xlarge`, see [Amazon EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/).
- For information about when to create a Virtual Private Cloud (VPC) and when to reuse an existing VPC, see [Resource Usage in Your Amazon Web Services Account](aws.md#aws-resources).