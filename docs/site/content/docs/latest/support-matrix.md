# Support Matrix
The following topic provides:
- A support matrix summary of the target platforms that are supported for each operating system.
- The supported operating systems and hardware/software required for your local bootstrap machine before you install Tanzu Community Edition. See [Local Client Machine](support-matrix/#local-client-machine) below.
- The supported target platforms you can bootstrap a cluster to. See [Target Platforms](support-matrix/#target-platforms) below.

## Support Matrix Summary
{{% include "/docs/assets/support-matrix.md" %}}
## Local Client Machine

Before you install Tanzu Community Edition, **one** of the following operating system and hardware/software configurations is required on your local machine.

{{% include "/docs/assets/prereq-linux.md" %}}

{{% include "/docs/assets/prereq-mac.md" %}}

{{% include "/docs/assets/prereq-windows.md" %}}




# Target Platforms

After you install Tanzu Community Edition on your local machine, you can use the Tanzu CLI to deploy a cluster to **one** of the following target platforms:


|Amazon EC2  |
|:------------------------ |
|Note: We do not support Photon on Amazon EC2|

|Microsoft Azure  |
|:------------------------ |
||

|Local Docker  |
|:------------------------|
|The following additional configuration is needed for the Docker engine on your local client machine: 6 GB of RAM and 4 CPUs (with no other containers running).<br> Check your Docker configuration as follows:<br>Linux: Run docker system info<br> Mac: Select Preferences > Resources > Advanced|
|15 GB of local machine disk storage for images |
|You cannot bootstrap a cluster to Docker from a Windows bootstrap machine, only Linux and Mac are supported at this time for Docker cluster deployments.|



|vSphere |
|:------------------------ |
||


