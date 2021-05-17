{{< tabs name="tab_with_hugo" >}}
{{% tab name="Hugo in a container" %}}

{{< note >}}
The commands below use Docker as default container engine. Set the `CONTAINER_ENGINE` environment variable to override this behaviour.
{{< /note >}}

1.  Build the image locally:

      ```bash
      # Use docker (default)
      make container-image

      ### OR ###

      # Use podman
      CONTAINER_ENGINE=podman make container-image
      ```

2. After building the `kubernetes-hugo` image locally, build and serve the site:

      ```bash
      # Use docker (default)
      make container-serve

      ### OR ###

      # Use podman
      CONTAINER_ENGINE=podman make container-serve
      ```

3.  In a web browser, navigate to `https://localhost:1313`. Hugo watches the
    changes and rebuilds the site as needed.

4.  To stop the local Hugo instance, go back to the terminal and type `Ctrl+C`,
    or close the terminal window.

{{% /tab %}}
{{% tab name="Hugo on the command line" %}}

Alternately, install and use the `hugo` command on your computer:

1.  Install the [Hugo](https://gohugo.io/getting-started/installing/) version specified in [`website/netlify.toml`](https://raw.githubusercontent.com/kubernetes/website/master/netlify.toml).

2.  If you have not updated your website repository, the `website/themes/docsy` directory is empty.
    The site cannot build without a local copy of the theme. To update the website theme, run:

    ```bash
    git submodule update --init --recursive --depth 1
    ```

3.  In a terminal, go to your Kubernetes website repository and start the Hugo server:

      ```bash
      cd <path_to_your_repo>/website
      hugo server --buildFuture
      ```

4.  In a web browser, navigate to `https://localhost:1313`. Hugo watches the
    changes and rebuilds the site as needed.

5.  To stop the local Hugo instance, go back to the terminal and type `Ctrl+C`,
    or close the terminal window.

{{% /tab %}}
{{< /tabs >}}