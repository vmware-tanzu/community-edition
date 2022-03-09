
# Deploy a Test Workload to a Workload Cluster

After you have provisioned a workload cluster, it is good practice to deploy and test a workload to validate cluster functionality.

Use the [kuard](https://github.com/kubernetes-up-and-running/kuard#demo-application-for-kubernetes-up-and-running) demo app to verify that your workload cluster is up and running.

## Prerequisites

You must have a workload cluster deployed.

## Procedure

1. Switch configuration context to the target workload cluster.

    ```sh
    kubectl config use-context <WORKLOAD-CLUSTER-NAME>
    ```

    For example:

    ```sh
    kubectl config use-context tce-cluster-1

    Switched to context "tce-cluster-1".
    ```

1. Deploy the kuard demo app.

    ```sh
    kubectl run --restart=Never --image=gcr.io/kuar-demo/kuard-amd64:blue kuard
    ```

    Expected result:

    ```sh
    pod/kuard created
    ```

1. Verify that the pod is running.

    ```sh
    kubectl get pods
    ```

    Expected result:

    ```sh
    NAME                     READY   STATUS    RESTARTS   AGE
    kuard                    1/1     Running   0          10d
    ```

1. Forward the pod container port 8080 to your local host port 8080.

    ```sh
    kubectl port-forward kuard 8080:8080
    ```

    Expected result:

    ```sh
    Forwarding from 127.0.0.1:8080 -> 8080
    Forwarding from [::1]:8080 -> 8080
    Handling connection for 8080
    ```

1. Using a browser go to ``http://localhost:8080``.
The kuard demo app web page appears which you can interact with and verify aspects of your cluster. For example, perform liveness and readiness probes.
1. Stop port forwarding by pressing Ctrl+C in the kubectl session.

1. Delete the kuard pod.

    ```sh
    kubectl delete pod kuard
    ```

    Expected result:

    ```sh
    pod "kuard" deleted
    ```

1. Verify that the pod is deleted.

    ```sh
    kubectl get pods
    ```
