# Create a Package and Package Repository Example

This example creates a package and a package repository.

Tanzu Community Edition provides the necessary components to create and run application platforms by leveraging a declarative packaging system that uses Kubernetes primitives to discover, install, and manage versioned packages in a running cluster.  

While most common use cases will reuse an existing package, this example explores how cluster operators can create their own package. It provides the steps to create a package using an OCI image, that can be discovered, installed, and managed using the package manager built into Tanzu Community Edition.

* A full description of the package creation process is available [here](package-creation-step-by-step).  
* A description of the Tanzu Community Edition packages architecture is available [here](architecture/#package-management).  
* A description of how to work with existing packages and package repositories is available [here](package-management).  
* Descriptions of a package and a package repository are available [here](glossary/#p)

## Prerequisites

{{% include "/docs/latest/assets/prereq-pkg-creation.md" %}}

* This examples uses GitHub Container Registry.

* Ensure you can access the following application which is used in this package creation example:

   ```sh
   ghcr.io/vladimirvivien/services/timeapp@sha256:898ad0977be31b2c16952a825e9505eb72328ec797351c8f634ace8cb02e5c5c
   ```

   It is a simple [Go](https://golang.org/) application which runs a webserver that returns the current time.  
   The application is compiled and published using [ko](https://github.com/google/ko) as a container image on GitHub Container Registry.  
   The application can be configured using the following two environment variables:  
   `TIME_FORMAT` - specifies the time layout using the Go’s time package layout idioms  
   `PORT` - specifies the port value on which the server will listen for requests  

Complete the following steps:

## 1. Create a Package

Define the package bundle, this will reference the published container image shown in the [Prerequisites](packages-create-example/#prerequisites). There are four steps required here:

* Create a directory structure.
* Create a [ytt](glossary/#ytt)-annotated Kubernetes configuration file.
* Use [kbld](glossary/#kbld) to resolve the referenced container image.
* Use [imgpkg](glossary/#imgpkg) to unify the Kubernetes configuration file and its referenced container image into a versioned OCI image and push it to an image registry.

1. Create a directory structure to store the artifacts needed for the package bundle:

    ```sh
    mkdir -p timeapp-package/config timeapp-package/.imgpkg
    ```

    Following the kapp-controller convention for package bundle format, the structure for the artifacts that make up the package will look like the following:

    ```sh
    timeapp-package/  
    └── .imgpkg/  
        └── images.yml  
    └── config/  
        └── package.yml  
        └── values.yml  
    ```

2. Create a [ytt](glossary/#ytt)-annotated Kubernetes configuration, `package.yaml` file for the application as follows:

    ```sh
    #@ load("@ytt:data", "data")

    #@ def labels():
    time-app: ""
    #@ end

    ---
    apiVersion: v1
    kind: Service
    metadata:
      namespace: default
      name: time-app
    spec:
      ports:
      - port: #@ data.values.**svc_port**
        targetPort: #@ data.values.**app_port**
      selector: #@ labels()
    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      namespace: default
      name: time-app
    spec:
      selector:
        matchLabels: #@ labels()
      template:
        metadata:
          labels: #@ labels()
        spec:
          containers:
          - name: time-app
            image: ghcr.io/vladimirvivien/services/timeapp@sha256:898ad0977be31b2c16952a825e9505eb72328ec797351c8f634ace8cb02e5c5c
            env:
            - name: TIME_FORMAT
              value: #@ data.values.**time_format**
    ```

    This file contains three value placeholders of interest:  
    `svc_port` for the service’s port value  
    `app_port` for the service’s targetPort  
    `time_format` for the application’s TIME_FORMAT environment variable.

    The actual values used in the ytt template, are stored in a file named `values.yaml`, which defines the default values that will be injected in the Kubernetes configuration above:

    ```sh
    #@data/values
    ---
    svc_port: 80
    app_port: 80
    time_format: ANSIC
    ```

3. Now that the content for the package bundle is complete, it is almost ready to be compiled and published as an OCI container image. Use [kbld](glossary/#kbld) to resolve the referenced container image before publishing the package as an OCI image. This step will generate a lock file that ensures build fidelity whenever the bundle is built. Run the following command:

    ```sh
    kbld -f timeapp-package/config/ --imgpkg-lock-output timeapp-package/.imgpkg/images.yml
    ```

4. At this time, the package bundle is ready to be compiled and published to an OCI image registry. To do this, use another Carvel tool, [imgpkg](glossary/#imgpkg), to unify the Kubernetes configuration file and its referenced container image into a versioned OCI image and push it to an image registry. Run the following command:

   ```sh
   imgpkg push -b ghcr.io/vladimirvivien/tanzu-packages/timeapp-pkg:1.0.0 -f timeapp-package/
   ```

    This command will push the OCI image `timeapp-pkg:1.0.0` to the `ghcr.io` image registry as a versioned package bundle ready to be installed and consumed by a Tanzu Community Edition cluster.

    At this point, the application binary is published as a Tanzu package bundle, and is accessible via GitHub’s package repository. However, it is not ready to be discovered and consumed by your cluster yet. For that to happen, it must be associated with a package repository.

## 2. Create a Package Repository

Before a Tanzu package can be deployed in a cluster, it must be made discoverable via a package repository. In most cases, a new package is added to an existing package repository so that users can easily discover and install it. As an illustrative exercise, however, this example walks through the creation and deployment of a Tanzu package repository. There are two steps required:

* Create a directory structure.
* Create the package repository content.

1. Create a directory structure to store the artifacts needed for the package repository. Run the following command to create the directory structure for the package repository:

    ```sh
    mkdir -p pkg-repo/.imgpkg pkg-repo/packages/timeapp-pkg.vladimirvivien.github.com
    ```

    This  command creates a directory named `packages`. This is where metadata configuration content for packages associated with the repository will be stored. Following the kapp-controller naming convention for package repositories, the configuration files for the package associated with this repository are placed in a directory with a fully-qualified name of `timeapp-pkg.vladimimirvivien.github.com`.

   The directory structure for the artifacts in the package repository should look like the following:

    ```sh
    pkg-repo/
    └── .imgpkg/
        └── images.yml
    └── packages/
        └── timeapp-pkg.vladimirvivien.github.com
            └── metadata.yml
            └── v1.0.0.yml
    ```

2. Create the package repository content. The content of the package repository consists of configuration YAMLs that represent a published package. Create Kubernetes custom resources to document the available packages so that the kapp-controller can surface that information for package discovery.

    a. The first CR is of type `PackageMetadata`, is a namespaced resource that is used to provide overall descriptive information about its associated published package. The information in the PackageMetadata CR is intended to be high-level and is consistent across many published versions of the package.  The following source snippet shows the YAML, saved in `metadata.yaml`, and describes the package:

    ```sh
    apiVersion: data.packaging.carvel.dev/v1alpha1
    kind: PackageMetadata
    metadata:
      # This will be the name of our package
      name: timeapp-pkg.vladimirvivien.github.com
    spec:
      displayName: "Time App"
      longDescription: "Simple application that returns time"
      shortDescription: "Simple time service"
    ```

   b. The next CR, of type Package, is a namespaced resource that is used to represent a published version of a package. In addition to descriptive information, this CR contains a reference to the OCI image of the package and configuration information that is used by the Tanzu package manager (kapp-controller) to install the package on the cluster.

   Following the naming convention for kapp-controller YAMLs, the custom resource is saved in file `v1.0.0.yaml` (below) to indicate that it is for a specific version. Note that the name attribute references the same qualified name in the previous Metadata CR.

    ```sh
    apiVersion: data.packaging.carvel.dev/v1alpha1
    kind: Package
    metadata:
      name: timeapp-pkg.vladimirvivien.github.com.1.0.0
    spec:
      refName: timeapp-pkg.vladimirvivien.github.com
      version: 1.0.0
      releaseNotes: |
        First release of the time service
      valuesSchema:
        openAPIv3:
          title: Values schema
          examples:
            - svc_port: 80
              app_port: 80
              time_format: ANSIC
          properties:
            svc_port:
              type: integer
              description: Port number for the service.
              default: 8080
              examples:
                - 8080
            app_port:
              type: integer
              description: Target port for the application.
              default: 8080
              examples:
                - 8080
            time_format:
              type: string
              description: Go language time package format (name or explicit format).
              default: ANSIC
              examples:
                - stranger
      template:
        spec:
          fetch:
            - imgpkgBundle:
                image: ghcr.io/vladimirvivien/tanzu-packages/timeapp-pkg@sha256:3b22b78f32f72642daea927a98588887d5fc74a9281599a7caf56925b7d63f26
          template:
            - ytt:
                paths:
                  - "config/"
            - kbld:
                paths:
                  - "-"
                  - ".imgpkg/images.yml"
          deploy:
            - kapp: {}
    ```

For each new published version of the package, a new custom resource YAML file can be added to the package repository, for example, `v1.0.1.yml`, `v1.0.2.yml`. This makes it possible for a package repository to surface multiple versions of the same software and for the package manager to facilitate version upgrade or downgrade.

## 3. Publish the Package Repository Image

Now that the content for the package repository is complete, it is ready to be compiled and published as an OCI container image. There are two steps required:

* Use kbld to resolve the referenced container image.
* Use imgpkg to compile and publish the package repository as an image to an OCI registry.

1. Use kbld to resolve the referenced container image before publishing the package as an OCI image.  Kbld will resolve the referenced OCI image of the package, based on its digest and generate a lock file that ensures the same image version is used in future builds:

    ```sh
    kbld -f pkg-repo/packages/ --imgpkg-lock-output pkg-repo/.imgpkg/images.yml
    ```

2. Use imgpkg to compile and publish the package repository as an image to an OCI registry. Imgpkg will unify the Kubernetes CR configuration files and the referenced container image into a package repository bundle (which itself is an OCI image):

    ```sh
    imgpkg push -b ghcr.io/vladimirvivien/tanzu-packages-repo:1.0.0 -f pkg-repo
    ```

    This command will push `tanzu-packages-repo:1.0.0` to the `ghcr.io` image registry as a versioned OCI image that represents a package bundle repository that can be installed by the Tanzu Community Edition internal package manager. You can now install the package onto a running cluster.

## 4. Add the Package Repository to Tanzu Community Edition

Before a package can be accessed in Tanzu Community Edition, the package manager must be made aware of its availability. There are two steps required:

* Use the `tanzu package repository add` command to register a new package.
* Use the `tanzu package repository list` to validate the repository was added.

1. Use the following command to register a new package repository with the cluster’s package manager:

    ```sh
    tanzu package repository add my-app-repo \
      --url ghcr.io/vladimirvivien/tanzu-packages-repo:1.0.0 \
      --namespace app-ns \
      --create-namespace
    ```

    This command adds the `ghcr.io/vladimirvivien/tanzu-packages-repo:1.0.0` image as a package repository to the `my-app-repo` cluster. Notice that the package manager will associate namespace `app-ns` with this repository and future interactions will need to specify the namespace.

    Internally, the package manager deploys the PackageRepository custom resource (defined above) allowing its packages to be discovered.

2. Validate the installation of the package repository using the following command:

    ```sh
    tanzu package repository list --all-namespaces
    ```

    This output shows that the package repository has been successfully installed and is ready to be used to install applications on the cluster:

    ```sh
    NAME            REPOSITORY                                               STATUS                NAMESPACE
    my-app-repo     ghcr.io/vladimirvivien/tanzu-packages-repo:1.0.0         Reconcile succeeded   app-ns
    tanzu-core      projects.registry.vmware.com/.../v1.21.2_vmware.1-tkg.1  Reconcile succeeded   tkg-system
    ```

3. (Optional) Because Tanzu packages are implemented using native Kubernetes primitives, you can also use the following kubectl command to validate the installation of the package repository:

    ```sh
    kubectl get packagerepositories -A
    ```

    The output should look similar to the following:

    ```sh
    NAMESPACE    NAME         AGE     DESCRIPTION
    app-ns       my-app-repo  4m42s   Reconcile succeeded
    tkg-system   tanzu-core   18m     Reconcile succeeded
    ```

## 5. Install the Package

After the repository has been successfully installed, you can install a  package from it.

1. Use the following command to inspect the repository and get a list of available packages:

    ```sh
    tanzu package available list --namespace app-ns
    ```

    The output should look similar to:

    ```sh
    NAME                                   DISPLAY-NAME  SHORT-DESCRIPTION
    timeapp-pkg.vladimirvivien.github.com  Time App      Simple time service
    ```

2. Use the following command to install version 1.0.0 of the `timeapp` package in `app-ns` namespace

    ```sh
    tanzu package install timeapp \
      --package-name timeapp-pkg.vladimirvivien.github.com \
      --version 1.0.0 \
      --namespace app-ns
    ```

3. Verify the package was installed correctly with the following command:

    ```sh
    tanzu package installed list --namespace app-ns
    ```

    The output should look similar to the following:

    ```txt
    NAME     PACKAGE-NAME                           PACKAGE-VERSION  STATUS
    timeapp  timeapp-pkg.vladimirvivien.github.com  1.0.0            Reconcile succeeded
    ```

4. (Optional) Internally, the package manager created a PackageInstall custom resource to represent the installation of the package on the cluster. To further understand the process, use the following kubectl command on the cluster:

    ```sh
    kubectl get pkgi --namespace app-ns
    ```

    The output should look similar to this:

    ```txt
    NAME      PACKAGE NAME                            PACKAGE VERSION   DESCRIPTION           AGE
    timeapp   timeapp-pkg.vladimirvivien.github.com   1.0.0             Reconcile succeeded   4m4s
    ```

## 6. Validate the Application Deployment

After the package manager successfully deploys the content of the package, the cluster should contain a Kubernetes deployment, along with its associated pods (as defined in the package YAML earlier).

1. Validate the Kubernetes deployment and its pods with the following commands:

    ```sh
    kubectl get deployments
    ```

    The output should look similar to:

    ```txt
    NAME       READY   UP-TO-DATE   AVAILABLE   AGE
    time-app   1/1     1            1           10m
    ```

    The cluster now has a deployment as was specified in the package.

2. Additionally, the following command shows the pods associated with the package are deployed successfully on the cluster:

    ```sh
    kubectl get pods -o wide
    ```

    The output should look similar to:

    ```txt
    NAME                        READY   STATUS    RESTARTS   AGE
    time-app-6c768bb7fb-xgvx8   1/1     Running   0          60m
    ```

## 7. Debug errors

Creating your own package involves multiple steps and you can easily encounter errors. Because Tanzu packages are expressed as pure Kubernetes primitives, you can use kubectl commands to debug.

1. The following shows an error message stating that the PackageRegistry resource is failing during deployment:

    ```sh
    kubectl get pkgr -A

    NAMESPACE    NAME             AGE   DESCRIPTION
    my-apps      my-custom-repo   26s   Reconcile failed: Fetching resources: Error...
    tkg-system   tanzu-core       42h   Reconcile succeeded

    ```

2. Another simple way to debug Tanzu package errors is to use the describe subcommand in kubectl to retrieve information about the failed resource:

    ```sh
    kubectl describe pkgr/my-app-repo --namespace app-ns

    Name:         my-app-repo
    Namespace:    app-ns
    Labels:       <none>
    Annotations:  <none>
    API Version:  packaging.carvel.dev/v1alpha1
    Kind:         PackageRepository
    ...

    Status:
      Conditions:
        Message:                       Deploying: Error (see .status.usefulErrorMessage for details)
        Status:                        True
        Type:                          ReconcileFailed
      Consecutive Reconcile Failures:  80
      Deploy:
        Error:       Deploying: Error (see .status.usefulErrorMessage for details)
        Exit Code:   1
        Finished:    true
        Started At:  2021-10-21T12:19:58Z
        Stderr:      kapp: Error: Applying create package/timeapp-pkg.vladimirvivien.github.com.1.0.0 (data.packaging.carvel.dev/v1alpha1) namespace: app-ns:
      Creating resource package/timeapp-pkg.vladimirvivien.github.com.1.0.0 (data.packaging.carvel.dev/v1alpha1) namespace: app-ns:  "timeapp-pkg.vladimirvivien.github.com.1.0.0" is invalid: metadata.name:
        Invalid value: "timeapp-pkg.vladimirvivien.github.com.1.0.0": must be <spec.refName> + '.' + <spec.version> (reason: Invalid)
        Stdout:  Target cluster 'https://100.64.0.1:443'

    ```
