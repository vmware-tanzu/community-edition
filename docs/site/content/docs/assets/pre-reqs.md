## Prerequisites

Tanzu Community Edition currently works on **macOS** and **Linux** AMD64 (also known as
x64) environments. Windows and other architectures may be supported in the
future.

The Docker runtime is required on the deployment machine, regardless of your
final deployment environment. Before proceeding, please ensure [Docker has
been installed](https://docs.docker.com/engine/install/) and is running.

Kubectl is required, download and install the latest version of `kubectl`.

    ```sh
    curl -LO https://dl.k8s.io/release/v1.20.1/bin/darwin/amd64/kubectl
    sudo install -o root -g wheel -m 0755 kubectl /usr/local/bin/kubectl
    ```
