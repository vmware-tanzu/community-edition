# Glossary

The following section provides a glossary of the main components and concepts involved in a Tanzu Community
Edition deployment.

[A](#a) | [B](#b) | [C](#c) | [E](#e) | [I](#i) | [K](#k) | [M](#m) | [O](#o) |[P](#p) | [T](#t) | [U](#u) | [V](#v) | [W](#w) | [Y](#w) |

## A

---

### Add-ons

Same as packages (see below).

## B

---

### Bootstrap

The bootstrap (noun) machine is the laptop, host, or server on which you download and run the Tanzu CLI. This is where the initial bootstrapping (verb) of a management cluster occurs before it is pushed to the platform where it will run. You run `tanzu`, `kubectl` and other commands on the bootstrap machine.

Using the Tanzu CLI to deploy a cluster to a target platform is often referred to as bootstrapping (verb).

## C

---

### Cert-manager

{{% include "/docs/assets/cert-manager-desc.md" %}}

## E

---

### Extensions

Same as packages (see below).

## I

---

### imgpkg

{{% include "/docs/assets/imgpkg-desc.md" %}}

## K

---

### kapp-controller

[kapp-controller](https://carvel.dev/kapp-controller/) is a Carvel tool and is the Tanzu Community Edition package manager. In Tanzu clusters, kapp-controller is constantly watching for package repositories. When a cluster is told about this package repository (likely via the Tanzu package repository command), kapp-controller can pull down that repository and make all the packages available to the cluster.

### Kbld

{{% include "/docs/assets/kbld-desc.md" %}}

### Kind cluster

During the deployment of the management cluster, either from the installer interface or the CLI,
Tanzu Community Edition creates a temporary management cluster using a [Kind](https://kind.sigs.k8s.io/) cluster on the bootstrap machine. Then, Tanzu Community Edition uses it to provision the final management cluster to the platform of your choice, depending on whether you are deploying to vSphere, Amazon EC2, Azure, or Docker. After the deployment of the management cluster finishes successfully, the temporary `kind` cluster is deleted.

## M

---

### Managed Cluster

{{% include "/docs/assets/mgmt-desc.md" %}}

## O

---

### OCI Registry

{{% include "/docs/assets/oci-desc.md" %}}

## P

---

### Package

{{% include "/docs/assets/package-description.md" %}}

### Package Repository

{{% include "/docs/assets/package-repository.md" %}}

## T

---

### Tanzu CLI

Tanzu CLI provides commands that facilitate many of the operations that you can perform with your clusters.
However, for certain operations, you still need to use `kubectl`.

### Tanzu Community Edition installer

The Tanzu Community Edition installer (the installer) is a graphical wizard that you launch in your browser by running the ``tanzu management-cluster create --ui`` command. The installer runs locally in a browser on the bootstrap machine and provides a user interface to guide you through the process of deploying a management cluster.

### Target Platform (Infrastructure Provider)

The target platform is the cloud provider or local Docker where you will deploy your cluster. This is also
referred to as your infrastructure provider.
There are four available target platforms:

* AWS
* Microsoft Azure
* Docker
* vSphere

## U

---

### Unmanaged Cluster

{{% include "/docs/assets/unmanaged-desc.md" %}}

## V

---

### Vendir

{{% include "/docs/assets/vendir-desc.md" %}}

## W

---

### Workload Cluster

After you deploy the management cluster, you can deploy a workload cluster. The workload cluster is deployed by the management cluster. The workload cluster is used to run your application workloads. The workload cluster is deployed using the Tanzu CLI.

## Y

---

### ytt

{{% include "/docs/assets/ytt-desc.md" %}}
