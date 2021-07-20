These tests are going to test the fluent-bit Package in the following way
- First checking the installation is done without error or not by storing the output the command ``` tanzu package install fluent-bit.tce.vmware.com ``` and checking for the desired result in the output.
- Next checking the number of ports by the command ```Kubectl get daemonset.apps/fluent-bit -n fluent-bit -o jsonpath={..status.desiredNumberScheduled}```the output should be 2.
- Next checking the pods availability by the command  ```Kubectl get daemonset.apps/fluent-bit -n fluent-bit -o jsonpath={..status.numberAvailable}`}``` the output should be 2.
- Now forward the port by command so that health can be checked for fluent-bit ```Kubectl port-forward daemonset/fluent-bit -n fuent-bit 56017:2020```.
- In last checking the health ```curl -I http://127.0.0.1:56017/api/v1/health```.
