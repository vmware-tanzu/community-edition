# Uninstall the Tanzu CLI

Complete the following steps to uninstall the Tanzu CLI and any associated configurations.

**Warning**: Uninstalling the Tanzu CLI will delete access to any managed or unmanaged clusters that were deployed with the Tanzu CLI.

## Linux

1. Run the following command to uninstall the Tanzu CLI:

   ```sh
   ~/.local/share/tce/uninstall.sh
   ```

   The `~/.local/share/tce` folder is deleted.

## MacOS

1. Run the following command to uninstall the Tanzu CLI:

   ```sh
   ~/Library/Application\ Support/tce/uninstall.sh
   ```

   The `~/Library/Application Support/tce` folder is deleted.

## Windows

1. Open a Command Prompt as an administrator and run the following command to uninstall the Tanzu CLI:

   ```sh
   run uninstall.bat
   ```

   The `%LocalAppData%\tce` folder is deleted.

For information about how to delete a management and workload cluster, see:

[Delete Management Cluster](delete-mgmt)  
[Delete Workload Cluster](delete-cluster)
