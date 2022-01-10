# vSphere Management + Workload Cluster E2E Test

This document is to talk about

- Subject Under Test (SUT)
- E2E Test strategy and the "Why?" behind it
- Handling failure scenarios

## Subject Under Test (SUT)

The subject under test for this particular E2E test is -

- `tanzu management-cluster create` - the command used to create management cluster on vSphere platform
- `tanzu cluster create` - the command used to create workload cluster on vSphere platform
- `tanzu management-cluster delete` - the command used to delete management cluster on vSphere platform
- `tanzu cluster delete` - the command used to delete workload cluster on vSphere platform

## E2E Test strategy and the "Why?" behind it

### Prerequisites to run the E2E test

Currently the vSphere E2E test requires some pre-requisites. The obvious thing is that we need a vSphere environment. We use VMC on AWS for this for our E2E tests. Steps to be done after provisioning a VMC environment -

- Get the connection details for the vCenter server. This would be the host name / IP and the credentials - the username and password for the vCenter. This will be used for `VSPHERE_SERVER`, `VSPHERE_USERNAME` and `VSPHERE_PASSWORD` environment variables

- Get the path of the Software Defined Data Center (SDDC) that you want to use. This will be used for `VSPHERE_DATACENTER` environment variable

- Create a SSH key pair using a tool like `ssh-keygen` and get the public key. This will be used for `VSPHERE_SSH_AUTHORIZED_KEY` and injected as `authorized_keys` in all the nodes of the cluster that is deployed. We create an RSA 4096 bits SSH key pair for our E2E tests.

- Get the path of the vSphere datastore to deploy the cluster. This will be used for `VSPHERE_DATASTORE`

- Get the path of the VM folder in which you want to place cluster VMs. This will be used for `VSPHERE_FOLDER`

- Get the network portgroup to assign each VM node to. This will be used for `VSPHERE_NETWORK`

- The network will have a DHCP server, change DHCP server settings such that some IPs (IP range) will be reserved to be used as static IPs and not to be used dynamically by the DHCP server. Two IPs from the static IPs list will be used as control plane endpoint IPs for `VSPHERE_MANAGEMENT_CLUSTER_ENDPOINT` and `VSPHERE_WORKLOAD_CLUSTER_ENDPOINT`

- Get the Name of the resource pool in which you want to place this cluster. This will be used for `VSPHERE_RESOURCE_POOL`

- In your vCenter, deploy a base OS image template corresponding to the TCE version used. This image is what will be used to deploy the cluster nodes. Our E2E tests use Photon OS, Version 3 and amd64 architecture, for both management and workload clusters on vSphere

For elaborate details on the vSphere deployment, please check the TCE docs site

### Networking

Since our VMC environment is all in a private network, we use a self hosted GitHub Action runner within the environment to be able to access the vCenter server

### Running the test

The E2E test's starting point is a bash script - `test/vsphere/run-tce-vsphere-management-and-workload-cluster.sh`. This script can be invoked standalone, or one can use `Makefile` target `vsphere-management-and-workload-cluster-e2e-test` to run the E2E tests like this `make vsphere-management-and-workload-cluster-e2e-test`

The script gives meaningful errors when required parameters (environment variables) are missing

### Scenario

The happy path scenario that we are testing is

- Run a management cluster on vSphere platform
- Run a workload cluster on vSphere platform
- Delete the workload cluster on vSphere platform
- Delete the management cluster on vSphere platform

This is all done using `tanzu` commands. If any of them fails, it means the test failed as the above are the Subjects Under Test.

### Handling failures

Given E2E tests test end to end a lot of things, a lot of things can go wrong, so failures are inevitable. And it's important to handle the failures

When there are expected / unexpected failures in the E2E test, we stop and try to do a cleanup. The cleanup is done using [`govc`](https://github.com/vmware/govmomi/tree/master/govc) which is CLI tool to interact with vCenter and works well for automation scripts

The cleanup helps us to cleanup any test resources so that no cloud resources used for testing are lying around. If anything is lying around not cleaned up we have to pay for it unnecessarily

For any failure we use `govc` to cleanup. Why? This is because the `tanzu` command failures cannot always be reverted using `tanzu`. For example, if management cluster creation fails, it's hard to use `tanzu` to delete it if was created half way through. And if it was not created at all, that is no resources created, it cannot be deleted using `tanzu`. It's hard to detect all these case with simple high level E2E tests. So whenever there's a failure, we fallback to using `govc` to do the cleanup

The `govc` cleanup order is always - first cleanup management cluster and only then cleanup workload cluster (if there's any). Why? As management cluster will keep reconciling to keep the workload cluster up and running as workload cluster's `cluster` resource would be in the management cluster, so we don't want to delete workload cluster first using `govc` and then have `govc` and management cluster fighting each other where management cluster reconciles to create the workload cluster and `govc` tries to delete it. Also, even if `govc` succeeds in deleting the workload cluster and moves on to deleting management cluster, it's possible that management cluster tries to create the workload cluster again by the time we delete the management cluster. But `govc` usually does clean up in seconds - like 10-20 seconds. In any case, we follow that order of cleanup in the E2E tests

We also have logs to mention if manual cleanup is needed in case `govc` cleanup fails for some reason. In future we will try to have alerts on slack or maybe create github issues in repo to track manual cleanup tasks to be taken care of. Alongside we can make the cleanup script more resilient by also running it separately in another standalone script in another job / workflow in the CI instead of the same workflow as the E2E tests within the same script
