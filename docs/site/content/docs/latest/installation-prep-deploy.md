Before you deploy a cluster, ensure you have the following prerequisites in place:

## 1. Docker
Complete the installation steps for your operating system: [](https://docs.docker.com/engine/install/)

## 2. Kubectl
Run the following command to install Kubectl, for more information, see the (Install Tools topic in the Kubernetes documentation)[https://kubernetes.io/docs/tasks/tools/]

linux

curl -LO https://dl.k8s.io/release/v1.20.1/bin/linux/amd64/kubectl
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

macOS

curl -LO https://dl.k8s.io/release/v1.20.1/bin/darwin/amd64/kubectl
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

## 3. Security
<!--- Can break this out into separate topics for each provider -->

### AWS