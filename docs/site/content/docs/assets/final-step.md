### Finalize the Deployment

1. Click **Review Configuration** to see the details of the management cluster that you have configured. When you click **Review Configuration**, Tanzu populates the cluster configuration file, which is located in the `~/.config/tanzu/tkg/clusterconfigs` subdirectory, with the settings that you specified in the interface. You can optionally copy the cluster configuration file without completing the deployment. You can copy the cluster configuration file to another bootstrap machine and deploy the management cluster from that machine. For example, you might do this so that you can deploy the management cluster from a bootstrap machine that does not have a Web browser.

1. (Optional) Under **CLI Command Equivalent**, click the **Copy** button to copy the CLI command for the configuration that you specified.

   Copying the CLI command allows you to reuse the command at the command line to deploy management clusters with the configuration that you specified in the interface. This can be useful if you want to automate management cluster deployment.

1. (Optional) Click **Edit Configuration** to return to the installer wizard to modify your configuration.
1. Click **Deploy Management Cluster**.

Deployment of the management cluster can take several minutes. The first run of `tanzu management-cluster create` takes longer than subsequent runs because it has to pull the required Docker images into the image store on your bootstrap machine. Subsequent runs do not require this step, so are faster. You can follow the progress of the deployment of the management cluster in the installer interface or in the terminal in which you ran `tanzu management-cluster create --ui`. If the machine on which you run `tanzu management-cluster create` shuts down or restarts before the local operations finish, the deployment will fail. If you inadvertently close the browser or browser tab in which the deployment is running before it finishes, the deployment continues in the terminal.
