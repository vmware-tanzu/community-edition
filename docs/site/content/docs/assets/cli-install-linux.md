1. Download the release for [Linux](https://github.com/vmware-tanzu/community-edition/releases/download/v0.7.0/tce-linux-amd64-v0.7.0.tar.gz) via web browser.

1. _[Alternative]_ Download the release using CLI. Alternatively, you may download a release using the provided remote script.

    ```sh
    curl -H "Authorization: token ${GITHUB_TOKEN}" \
        -H "Accept: application/vnd.github.v3.raw" \
        -L https://api.github.com/repos/vmware-tanzu/community-edition/contents/hack/get-tce-release.sh | \
        bash -s <RELEASE-VERSION-DISTRIBUTION>
    ```

    > - Where ``<RELEASE-VERSION-DISTRIBUTION>`` is the Tanzu Community Edition release version and distribution. This is a required argument for `bash` to download a releases. For example, to download v0.7.0 for Linux, provide:  <br>`bash -s v0.7.0 linux`
    > - This script requires `curl`, `grep`, `sed`, `tr`, and `jq` in order to work
    > - The release will be downloaded to the local directory as `tce-linux-amd64-v0.7.0.tar.gz`
    > - *_Note:_* This _currently_ requires the use of a GitHub personal access token.
      Follow [the GitHub documentation](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token) to aquire and use a personal access token.

1. Unpack the release.

    ```sh
    tar xzvf ~/<DOWNLOAD-DIR>/tce-linux-amd64-v0.7.0.tar.gz
    ```

1. Run the install script (make sure to use the appropriate directory for your platform).

    ```sh
    cd tce-linux-amd64-v0.7.0
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
