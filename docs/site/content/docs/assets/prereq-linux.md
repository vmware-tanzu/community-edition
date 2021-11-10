### Linux Local Bootstrap Machine Prerequisites
||
|:--- |
|RAM: 6 GB|
|CPU: 2|
|[Docker](https://docs.docker.com/engine/install/) <BR> In Docker, you must create the docker group and add your user before you attempt to create a standalone or management cluster. Complete steps 1 to 4 in the [Manage Docker as a non-root user](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user) procedure in the Docker documentation.|
|[Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) |
|Latest version of Chrome, Firefox, Safari, Internet Explorer, or  Edge|
|System time is synchronized with a Network Time Protocol (NTP) server.|
|Ensure your Linux bootstrap machine is using cgroup v1, for more information, see **Check and set the cgroup** below.|

#### Check and set the cgroup 

1. Check the cgroup by running the following command:

    ```sh
    docker info | grep -i cgroup 
    ```

    You should see the following output:

    ```sh
    Cgroup Driver: cgroupfs
    Cgroup Version: 1
    ```

2. If your Linux distribution is configured to use cgroups v2, you will need to set the `systemd.unified_cgroup_hierarchy=0` kernel parameter to restore cgroups v1. See the instructions for setting kernel parameters for your Linux distribution, including:

    [Fedora 32+](https://fedoramagazine.org/docker-and-fedora-32/)  
    [Arch Linux](https://wiki.archlinux.org/title/Kernel_parameters)  
    [OpenSUSE](https://doc.opensuse.org/documentation/leap/reference/html/book-reference/cha-grub2.html)
