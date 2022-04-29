---
title: "Comparing Managed and Unmanaged Clusters on VMware Tanzu Community Edition"
slug: managed-unmanaged-clusters-tanzu-community-edition
date: 2022-04-06
author: Steven Pousty
image: /img/K8s-Cluster-img1.png
excerpt: "VMware Tanzu Community Edition offers two different Kubernetes cluster configuration options: managed and unmanaged. In this blog I will walk you through a comparison of managed and unmanaged clusters, some benefits and drawbacks, and then some use cases where each type works best."
tags: ['Steven Pousty']
---
 
Greetings, friends! VMware Tanzu Community Edition offers two different Kubernetes cluster configuration options: managed and unmanaged. In this blog I will walk you through a comparison of managed and unmanaged clusters, some benefits and drawbacks, and then some use cases where each type works best.

## Back to basics — a Kubernetes cluster

To ensure we are on the same page, I want to cover the high-level basics of what makes up a Kubernetes cluster. But first, it’s important to note that in this case, the word cluster does not have to mean multiple computers or even virtual machines; it can be multiple running containers.

One major division in a Kubernetes cluster is control plane versus nodes (formerly called master nodes versus worker nodes). From the [Kubernetes documentation](https://kubernetes.io/docs/concepts/overview/components/):

Nodes:
“A Kubernetes cluster consists of a set of worker machines, called [nodes](https://kubernetes.io/docs/concepts/architecture/nodes/), that run containerized applications…. The worker node(s) host the [Pods](https://kubernetes.io/docs/concepts/workloads/pods/) that are the components of the application workload.”

Control plane:
“The [control plane](https://kubernetes.io/docs/reference/glossary/?all=true#term-control-plane) manages the worker nodes and the Pods in the cluster.” Typically, the control plane will be spread across multiple “machines.”

While the control plane and the Node can be on the same machine, typically, a Kubernetes cluster consists of one or more control plane machines and one or more nodes.

!["Kubernetes cluster showing one or more control plane machines and one or more nodes"](/img/K8s-Cluster-img1.png)

## Unmanaged clusters

An [unmanaged cluster](https://tanzucommunityedition.io/docs/v0.11/getting-started-unmanaged/) is just a plain Kubernetes cluster spun up using your favorite tools like [kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/) or vRealize Automation. In the case of Tanzu Community Edition, you can use a [single command](https://tanzucommunityedition.io/docs/v0.11/getting-started-unmanaged/) with the Tanzu CLI, and it will spin up a Kubernetes-in-Docker (*kind*) cluster on your local machine. In a kind cluster, all the pieces discussed above, including Nodes, are actually containers running in your local container engine (such as Docker Desktop).

!["Tanzu Community Edition Unmanaged Cluster Control Plane and Nodes in Container"](/img/Cluster-Container-img2.png)

Let’s take a look at some of the benefits and drawbacks as well as the use cases for unmanaged clusters in Tanzu Community Edition.

### Benefits and drawbacks

Some of the nice features of unmanaged clusters are:

1. Really fast start-up and tear down since it is just a single container. Kind clusters start a root container and then bring up all the pieces of the cluster inside that root container.
2. Relatively small resource usage because we are only bringing up a single “cluster” and containers. Given the efficient nature of containers we will only use the memory needed to run their services.
3. All the pieces are containers running in the root container, making clean up as easy as deleting a single container.
4. Once the container images are downloaded, all work can be done with slow or low internet connectivity.
5. You pay no costs to a cloud provider.

To get these features, here are some of the tradeoffs you will have to make:

1. These clusters will have to be managed in a one-off manner with manual steps or automation tools.
2. Currently, the clusters can only run on your local machine. If you do not have a container engine on your machine or do not have enough resources, this won’t work for you.
3. Because adding nodes or more instances of control plane services to your cluster only adds containers to your root container, you can’t get true scaling or high availability (HA). You can simulate and test things like failover, but it is only within that root container.
4. Unmanaged clusters are fundamentally different in several ways from managed clusters, and therefore will not be an exact representation of managed development, staging, or production clusters.

### Use cases

With these benefits and drawbacks in mind, here are some of the typical use cases for an unmanaged cluster:

1. Given the quick spin up and tear down, unmanaged clusters are a perfect platform for learning the basics of Kubernetes and Tanzu. It’s a quick way to get started and if you mess up your installation, it’s not a lot of overhead to just delete and recreate the cluster.
2. Given the speed and ease of both cluster setup and cleanup, doing local experimentation in an unmanaged cluster before trying things in a cluster running in a cloud provider is another great fit.
3. The same goes for local development, where a quick turnaround of seeing your code is important. Since it is all on your local machine, attaching debuggers and working with files will be faster too.
4. If you want to test if your application containers will run in Kubernetes and how you might write the YAML for health probes, an unmanaged cluster is a great location.
5. Continuous Integration and Continuous Delivery (CI/CD) pipelines needing disposable Tanzu/Kubernetes instances in non-production validation parts of the pipeline are well served by unmanaged clusters. This use case was one of the primary motivations for the unmanaged cluster architecture.

## Managed clusters

[Managed clusters](https://tanzucommunityedition.io/docs/v0.11/getting-started/) take all the goodness we get from declarative infrastructure, Kubernetes, [Custom Resource Definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/), and the Kubernetes [Cluster API](https://cluster-api.sigs.k8s.io/) to allow you to use Kubernetes clusters to manage the clusters. The basic idea is to create one Kubernetes cluster (a management cluster) and then use that cluster to create and manage other Kubernetes clusters where we do all work (workload clusters).

The management cluster uses the Kubernetes Cluster API to manage the full lifecycle of other Kubernetes clusters. The workload cluster is where you run all your applications, databases, and other containerized goodness. Typically, the management cluster is on separate machines/infrastructure from any of the workload clusters as its only job should be managing its clusters.

A management cluster with multiple workload clusters is the typical pattern recommended for the majority of use cases.

!["Managed cluster control planes with workloads"](/img/mgmt-cluster-controlplane-img3.png)

### Benefits and drawbacks

There are a lot of benefits we get from using this type of architecture:

1. The control plane and worker nodes (for both the management and workload clusters) can actually exist on multiple infrastructure nodes. Because of this, you can achieve true HA and scalability.
2. As opposed to a shell script, or other [imperative-based infrastructure](https://openupthecloud.com/declarative-vs-imperative-infra/) creation techniques, cluster creation can be done declaratively. You state what resources you want available in your Workload cluster, and the management cluster takes care of making that a reality.
3. There are simple, easy commands to create workload clusters.
4. Management clusters simplify workload cluster management tasks, such as upgrading or scaling nodes, into single commands. The management cluster understands all the pieces and resources that need to be handled.
5. Because you use your management cluster to manage all your workload clusters, you now have a centralized place where you can query and manage clusters.
6. The management cluster stores the credentials needed to create workload clusters. If you give end users access to that management cluster, they can spin up workload clusters without the cloud provider credentials being revealed.
7. Unlike creating the management cluster or unmanaged clusters, delegated users do not need Docker on their machine to create workload clusters. Instead, all they need is [kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/).

With these benefits, there are also some drawbacks:

1. The management cluster is yet another Kubernetes cluster to manage, support, and incur costs for the infrastructure. Granted, management clusters are usually relatively small, but the overhead still exists.
2. Given the importance of the management cluster, it is not recommended to run workloads on this management cluster. You will also typically run multiple control plane nodes for your management clusters to ensure high availability.
3. It will take longer to provision a workload cluster, especially when starting without a management cluster, than it will to provision an unmanaged Tanzu Community Edition cluster. Even if a management cluster already exists, it will still take longer because you need to provision infrastructure in a cloud provider, which is slower than spinning up containers on your local machine.

### Use cases

With this set of tradeoffs, there are important use cases that benefit from this architecture:

1. Production Tanzu clusters should always use a managed clusters architecture.
2. If you need your development clusters to match production clusters, then the management cluster can ensure more similarity between the environments.
3. Running multiple production, staging, and development Kubernetes clusters (basically, a fleet of workload clusters) are one of the prime use cases since you centralize all the cluster information and methods to query, as well as update, the information.
4. If you allow people to create their own workload clusters but tie the costs to a single set of cloud provider credentials, this pattern can be useful. It is beneficial that you don’t have to expose the sensitive provider credentials to the user, which would allow them to do whatever they want in the cloud provider.

## Take home

We covered a lot of ground in this post, but I think we have set you up with a clearer understanding of two major tools you can add to your Tanzu and modern app toolbelt. Depending on your needs, you can now go with the quick and lightweight unmanaged Tanzu clusters, which are great for work on your local machine. Or you can go with the production-ready, management-focused style from a managed clusters architecture, which works anywhere from your local machine to your favorite cloud provider.

Have fun playing with these new tools from Tanzu Community Edition, and be sure to [stop by](https://tanzucommunityedition.io/community/) and let us know what you think!
