## Installation Procedure

1. You must download and install the latest version of `kubectl`. For more information, see [Install and Set Up kubectl on Linux](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) in the Kubernetes documentation.

1. You must download and install the latest version of `docker`. For more information, see [Install Docker Engine](https://docs.docker.com/engine/install/) in the Docker documentation.

### Option 1: Homebrew

1. Make sure you have the [Homebrew package manager installed](https://brew.sh/)

1. Run the following in your terminal:

    ```sh
    brew install vmware-tanzu/tanzu/tanzu-community-edition
    ```

1. Run the post install configuration script. Note the output of the `brew install` step for the correct location of the configure script:

    ```sh
    {HOMEBREW-INSTALL-LOCATION}/configure-tce.sh
    ```

    > This puts all the Tanzu plugins in the correct location.
    > The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.

### Option 2: Curl GitHub release

1. Download the release for [Linux](https://github.com/vmware-tanzu/community-edition/releases/download/{{< release_latest >}}/tce-linux-amd64-{{< release_latest >}}.tar.gz) via web browser.

1. _[Alternative]_ Download the release using the CLI. You may download a release using the provided remote script piped into bash.

    ```sh
    curl -H "Accept: application/vnd.github.v3.raw" \
        -L https://api.github.com/repos/vmware-tanzu/community-edition/contents/hack/get-tce-release.sh | \
        bash -s <RELEASE-VERSION> <RELEASE-OS-DISTRIBUTION>
    ```

    > - Where ``<RELEASE-VERSION>`` is the Tanzu Community Edition release version. This is a required argument.
    > - Where ``<RELEASE-OS-DISTRIBUTION>`` is the Tanzu Community Edition release version and distribution. This is a required argument.
    > - For example, to download {{< release_latest >}} for Linux, provide:  <br>`bash -s {{< release_latest >}} linux`
    > - This script requires `curl`, `grep`, `sed`, `tr`, and `jq` in order to work
    > - The release will be downloaded to the local directory as `tce-linux-amd64-{{< release_latest >}}.tar.gz`
    > - *_Note:_* A GitHub personal access token may be provided to the script as the `GITHUB_TOKEN` environment variable. This bypasses GitHub API rate limiting but is _not_ required. Follow [the GitHub documentation](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token) to acquire and use a personal access token.

1. Unpack the release.

    ```sh
    tar xzvf ~/<DOWNLOAD-DIR>/tce-linux-amd64-{{< release_latest >}}.tar.gz
    ```

1. Run the install script (make sure to use the appropriate directory for your platform).

    ```sh
    cd tce-linux-amd64-{{< release_latest >}}
    ./install.sh
    ```

    > This installs the `Tanzu` CLI and puts all the plugins in the correct location.
    > The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.

1. You must download and install the latest version of `kubectl`.

    ```sh
    curl -LO https://dl.k8s.io/release/v1.20.1/bin/linux/amd64/kubectl
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    ```

    For more information, see [Install and Set Up kubectl on Linux](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) in the Kubernetes documentation.
