### Linux Local Bootstrap Machine Prerequisites
||
|:--- |
|Arch: x86; ARM is currently unsupported|
|RAM: 6 GB|
|CPU: 2|
|[Docker](https://docs.docker.com/engine/install/) <BR> Add your non-root user account to the docker user group. Create the group if it does not already exist. This lets the Tanzu CLI access the Docker socket, which is owned by the root user. For more information, see steps 1 to 4 in the [Manage Docker as a non-root user](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user) procedure in the Docker documentation.|
|[Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) |
|Latest version of Chrome, Firefox, Safari, Internet Explorer, or  Edge|
|System time is synchronized with a Network Time Protocol (NTP) server.|
|Ensure your bootstrap machine is using [cgroup v1](https://man7.org/linux/man-pages/man7/cgroups.7.html). For more information, see [Check and set the cgroup](../support-matrix/#check-and-set-the-cgroup).|
