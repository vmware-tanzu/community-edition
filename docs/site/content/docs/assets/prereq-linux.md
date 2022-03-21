### Linux Local Bootstrap Machine Prerequisites
||
|:--- |
|Arch: x86; ARM is currently unsupported|
|RAM: 6 GB|
|CPU: 2|
|[Docker](https://docs.docker.com/engine/install/) <BR> In Docker, you must create the docker group and add your user before you attempt to create a management cluster. Complete steps 1 to 4 in the [Manage Docker as a non-root user](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user) procedure in the Docker documentation.|
|[Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) |
|Latest version of Chrome, Firefox, Safari, Internet Explorer, or  Edge|
|System time is synchronized with a Network Time Protocol (NTP) server.|
|Ensure your bootstrap machine is using [cgroup v1](https://man7.org/linux/man-pages/man7/cgroups.7.html). For more information, see [Check and set the cgroup](../support-matrix/#check-and-set-the-cgroup).|
