### Mac Local Bootstrap Machine Prerequisites

||
|:--- |
|Arch: x86; ARM (M1) currently unsupported |
|RAM: 6 GB |
|CPU: 2|
|[Docker Desktop for Mac; Version <= 4.2.0](https://docs.docker.com/desktop/mac/release-notes/#docker-desktop-420)|
|[Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/) |
|Latest version of Chrome, Firefox, Safari, Internet Explorer, or  Edge|

#### Check the cgroup version 

1. Check the cgroup by running the following command:

    ```sh
    docker info | grep -i cgroup 
    ```

    You should see the following output:

    ```sh
    Cgroup Driver: cgroupfs
    Cgroup Version: 1
    ```

2. If you see cgroup version 2, you are running an incompatible version of
   Docker Desktop. To resolve this, we recommend running [Docker Desktop
   4.2.0](https://docs.docker.com/desktop/mac/release-notes/#docker-desktop-420).

    > In a future release, we'll support cgroupsv2 which will resolve this issue.
    > Please follow [issue
    > 2798](https://github.com/vmware-tanzu/community-edition/issues/2798) for
    > progress.

