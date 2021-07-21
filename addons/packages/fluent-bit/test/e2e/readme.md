# THIS CONTENT HAS MOVED TO THE DOCS BRANCH IN:  PLEASE MAKE ANY FURTHER UPDATES THERE

File is available here on docs branch: ``docs\site\content\docs\latest\fluent-bit-config``

## Fluent-bit

This package collect any data like metrics and logs from different sources,
enrich them with filters and send them to multiple destinations using
[fluent-bit](https://github.com/fluent/fluent-bit).

## Components

* fluent-bit: open source Log Processor and Forwarder.

## Configuration

NA

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy fluent-bit. |

### fluent-bit Configuration

_Currently there is no fluent-bit customization_.

## Usage Example

In order to test Fluent-bit , the below manual steps are to be performed.The same steps are used for automation by using go Ginkgo framework.
-First checking the installation is done without error or not by storing the output the command ``` tanzu package install fluent-bit.tce.vmware.com ``` and checking for the desired result in the output.
-Next checking the number of ports by the command ```Kubectl get daemonset.apps/fluent-bit -n fluent-bit -o jsonpath={..status.desiredNumberScheduled}```the output should be 2.
-Next checking the pods availability by the command  ```Kubectl get daemonset.apps/fluent-bit -n fluent-bit -o jsonpath={..status.numberAvailable}`}``` the output should be 2.
-Now forward the port by command so that health can be checked for fluent-bit ```Kubectl port-forward daemonset/fluent-bit -n fuent-bit 56017:2020```.
-In last checking the health ```curl -I http://127.0.0.1:56017/api/v1/health```.

## Prerequisites

-Before running the suite , the cluster should be up and running.Fluent-bit test suite will never create or destroy any cluster.

## How to run the tests

-To run individual suits , ```cd addons/packages/fluent-bit/test/e2e``` from the tce root path and run ```ginkgo -v -r```
