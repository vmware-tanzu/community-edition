1. Download the release for [Linux](https://github.com/vmware-tanzu/tce/releases/download/v0.5.0/tce-linux-amd64-v0.5.0.tar.gz) via web browser.

1. _[Alternative]_ Download the release using CLI

    ```sh
    curl -H "Authorization: token ${GH_ACCESS_TOKEN}" \
        -H "Accept: application/vnd.github.v3.raw" \
        -L https://api.github.com/repos/vmware-tanzu/tce/contents/hack/get-tce-release.sh | \
        bash -s RELEASE_VERSION DISTRIBUTION
    ```

    > Alternatively, you may download a release using the provided remote script.
    > - The TCE release version and release distribution _must_ be passed as arguments to `bash` in order to download a releases. So, for example, to download v0.6.0 for Linux, simple provide `bash -s v0.6.0 linux` as arguments.
    > - This script requires `curl`, `grep`, `sed`, `tr`, and `jq` in order to work
    > - The release will be downloaded to the local directory as `tce-linux-amd64-v0.6.0.tar.gz`
    > - *_Note:_* This _currently_ requires the use of a GitHub personal access token.
      Follow [the GitHub documentation](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token) to aquire and use a personal access token.

1. Unpack the release.

    ```sh
    tar xzvf ~/Downloads/tce-linux-amd64-v0.5.0.tar.gz
    ```

1. Run the install script (make sure to use the appropriate directory for your platform).

    ```sh
    cd tce-linux-amd64-v0.5.0
    ./install.sh
    ```

    > This installs the `Tanzu` CLI and puts all the plugins in their proper location.
    > The first time you run the `tanzu` command the installed plugins and plugin repositories are initialized. This action might take a minute.

1. If you wish to run commands against any of the Kubernetes clusters that are created, you will need to download and install the latest version of `kubectl`.

    ```sh
    curl -LO https://dl.k8s.io/release/v1.20.1/bin/linux/amd64/kubectl
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    ```
