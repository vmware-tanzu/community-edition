| Architecture   | CPU | RAM  | Required software |
|:---------------|:----|:-----|:------------------|
| x86_64 / amd64 | 1   | 2 GB | [Docker engine](https://docs.docker.com/engine/install), (Optional) [minikube](https://minikube.sigs.k8s.io/docs/start/)|

Tanzu Community Edition supports two cluster providers for unmanaged clusters: Kind and minikube

- Kind is the default cluster provider and is included as default with the unmanaged cluster binary, you just need to install Docker.
- minikube is an alternative cluster provider, if you plan to use minikube as your cluster provider, you must first install minikube and a minikube supported container or virtual machine manager such as Docker.
