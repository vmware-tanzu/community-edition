# Secure automated Ingress on AWS

Now that you have provisioned your first workload cluster, it’s time to use it. If you want to deploy applications on this cluster you will need some packages installed on it that will give you some capabilities to improve the usage experience of the cluster.

In this guide, you are going to learn how to enable access to your applications in an easy and secure way.

We will use a workload cluster provisioned on Amazon Web Services (AWS), but the steps should be similar on other infrastructure providers. You need a [DNS registered in Route53](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/registrar.html). This DNS name will form the URLs of your applications. We will use **example.com** as the sample domain throughout the rest of the guide.

You will learn how to install and properly configure the following Kubernetes services:

* **Ingress controller**. This is the entry point to your cluster. Requests will most likely be routed to applications running in Kubernetes via an Ingress Controller.
* **Certificate management**. Because your services should be exposed using a secure transport (HTTPS), you need to generate TLS certificates that your ingress controller will use to secure your application's communication with the public.
* **DNS management**. You want your applications to listen on a specific domain name, but you don’t know (or want) to learn how to configure Route53 or your DNS provider. With automatic DNS management integration, your services will directly be exposed where you expect them.

With these services deployed, it becomes easier to deploy an application on a cluster and access it through a named URL.

## Deploy an application

Since you need an application to verify everything works, let’s deploy one with a few easy commands. This application will be deployed in your current namespace, so first you will create a namespace for this application. By doing this, you can easily clean up afterwards.

```shell
kubectl create namespace sample
```

Now that we have a namespace created, we can deploy the sample application. The sample application is a simple application that displays the Tanzu Community Edition logo.

```shell
kubectl create deployment sample --image=projects.registry.vmware.com/tce/logo:latest -n sample
```

You can now verify that the deployment is running:

```shell
kubectl get deployment sample -n sample
```

The output should look similar to:

```shell
NAME     READY   UP-TO-DATE   AVAILABLE   AGE
sample   1/1     1            1           11s
```

To access this application you need a [Kubernetes Service](https://kubernetes.io/docs/concepts/services-networking/service/) fronting it. This application is listening on port 8080, so you will expose this port within the Kubernetes cluster.

```shell
kubectl expose deployment sample --port=8080 -n sample
```

To get the details on the created service, execute:

```shell
kubectl get service sample -n sample
```

The output should look similar to:

```shell
NAME     TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
sample   ClusterIP   100.71.54.99   <none>        8080/TCP   19s
```

Now, to verify that the application is running, and since that service is only exposed within the cluster, you will need to do a port-forward to that service:

```shell
kubectl port-forward svc/sample 8080 -n sample
```

And now in a different terminal run:

```shell
curl localhost:8080
```

The output should look similar to:

![CLI Architecture](../img/tanzu-community-edition-ascciart.png)

This means that your application is up and running.

Now let’s move on to how you can access this same service in a secure way outside the cluster using a domain name.

## Ingress Controller

By default, when you deploy an application to Kubernetes, it is only accessible from inside of the cluster. If it is a web application and you want to make it public outside of the cluster, you have a number of options available.

### 1. NodePort

The first option associates Kubernetes Service object of type [NodePort](https://kubernetes.io/docs/concepts/services-networking/service/#nodeport) with the application.

This type of service will expose the application via a static port on each node in the cluster and if the node IP and port are made public from outside of the cluster it would then be accessible by an external user.

The port number for a NodePort service can either be assigned or chosen, but needs to be in a set range of ports set aside for NodePort services (default ports are `30000-32767`). Different applications need to use different ports and the standard HTTP service ports (`:80`,`:443`) would not be able to be used.

### 2. LoadBalancer

The second option is to associate a Kubernetes Service object of type **LoadBalancer** with the `Service`.

This type of service integrates with the underlying infrastructure provider used by the Kubernetes cluster to automatically assign an external IP address unique to the application and routes traffic for that IP to the application. This enables applications to be exposed on the standard HTTP service ports.

Because each application is assigned a unique IP address, you would need to have as many IP addresses available as you have applications you want to make public. This is because using a LoadBalancer type service alone, applications cannot share an IP address.

Where a Kubernetes cluster only hosts one web application that needs to be made public, using a LoadBalancer service type is often sufficient. If however, you need to expose multiple web applications, it is better to configure the Kubernetes cluster with an ingress controller.

### 3. Contour Ingress Controller

The third option is to install and configure an ingress controller. An ingress controller is a HTTP proxy service that can accept requests for many different hostnames and route traffic through to the appropriate application in the Kubernetes cluster.

Tanzu Community Edition includes the open source [Contour](https://projectcontour.io/) ingress controller. Contour utilizes [Envoy Proxy](https://www.envoyproxy.io/) for routing. Standard Kubernetes Ingress resources for creating ingresses are supported, along with extended resources which provide additional features and flexibility above what the standard Ingress type provides. The readme for the Tanzu Community Edition Contour package is [here](https://tanzucommunityedition.io/docs/package-readme-contour-1.19.1/).

When installing the Contour ingress controller, by default it will use a Kubernetes Service of type LoadBalancer to expose the ingress controller router externally to the Kubernetes cluster. Since you are going to be using AWS for this guide, this is the type of service you would want.

Complete the following steps:

1. Install the Tanzu Community Edition package repository into the `tanzu-package-repo-global` namespace.

    ```sh
    tanzu package repository add tce-repo --url projects.registry.vmware.com/tce/main:{{< pkg_repo_latest >}} --namespace tanzu-package-repo-global
    ```

1. To list details for the Contour ingress controller package run:

    ```shell
    tanzu package available get contour.community.tanzu.vmware.com
    ```

    The output should look similar to:

    ```shell
    - Retrieving package details for contour.community.tanzu.vmware.com...
    NAME:                 contour.community.tanzu.vmware.com
    DISPLAY-NAME:         Contour
    SHORT-DESCRIPTION:    An ingress controller
    PACKAGE-PROVIDER:     VMware
    LONG-DESCRIPTION:     An Envoy-based ingress controller that supports dynamic configuration updates and multi-team ingress delegation. See https://projectcontour.io for more information.
    MAINTAINERS:          [{Steve Kriss} {Steve Sloka} {Nick Young} {Sunjay Bhatia} {Nicholas Seemiller}]
    SUPPORT:              Go to https://projectcontour.io for documentation or the #contour channel on Kubernetes slack
    CATEGORY:             [ingress]
    ```

1. To list the versions available for the Contour package, run:

    ```shell
    tanzu package available list contour.community.tanzu.vmware.com
    ```

    The output should look similar to:

    ```shell
    - Retrieving package versions for contour.community.tanzu.vmware.com...
      NAME                                VERSION  RELEASED-AT
      contour.community.tanzu.vmware.com  1.17.1   2021-07-23 19:00:00 +0100 BST
      contour.community.tanzu.vmware.com  1.17.2   2021-07-23 19:00:00 +0100 BST
      contour.community.tanzu.vmware.com  1.18.1   2021-07-23 19:00:00 +0100 BST
      contour.community.tanzu.vmware.com  1.19.1   2021-10-26 01:00:00 +0100 BST
    ```

1. Complete the following steps to install Contour:

      1. You will need to provide some custom configuration specific to your environment. Note that although the options listed in the [Contour package readme](https://tanzucommunityedition.io/docs/package-readme-contour-1.19.1/) are shown in a flat namespace, the data input values need to be supplied as a hierarchical YAML file definition.

          As you want to have an **external-dns** managing your DNS registrations, you will add an annotation to Contour for the [wildcard DNS](https://en.wikipedia.org/wiki/Wildcard_DNS_record) you want to use. You will configure external-dns later to manage this registration. In this example, we will use `*.example.com`.

          Create a file `contour-config.yaml` which includes:

          ```shell
          envoy:
            service:
              annotations:
                external-dns.alpha.kubernetes.io/hostname: "*.example.com"
          ```

      1. Ensure that the file is saved and then install the Contour package by running:

          ```shell
          tanzu package install contour --package-name contour.community.tanzu.vmware.com --version 1.19.1 --values-file contour-config.yaml
          ```

          The output from the command should look similar to:

          ```shell
          | Installing package 'contour.community.tanzu.vmware.com'
          | Getting namespace 'default'
          | Getting package metadata for 'contour.community.tanzu.vmware.com'
          | Creating service account 'contour-default-sa'
          | Creating cluster admin role 'contour-default-cluster-role'
          | Creating cluster role binding 'contour-default-cluster-rolebinding'
          | Creating secret 'contour-default-values'
          - Creating package resource
          | Package install status: Reconciling

          Added installed package 'contour' in namespace 'default'
          ```

      1. You now have an AWS load balancer that points to your Contour Ingress controller. You can verify it by running:

          ```shell
          kubectl get svc/envoy -n projectcontour
          ```

          The output should look similar to:

          ```shell
          NAME    TYPE           CLUSTER-IP     EXTERNAL-IP                                                               PORT(S)                      AGE
          envoy   LoadBalancer   100.65.70.13   aca8e02fba5614ce0b828f234f4f404f-1196576133.eu-west-1.elb.amazonaws.com   80:30033/TCP,443:30922/TCP   78s
          ```

1. Now that you have an ingress controller, you can access your application through the controller instead of the previously used port-forwarding mechanism.

    To create an ingress for the sample application, execute:

    ```shell
    cat <<EOF | kubectl apply -n sample -f -
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: sample
    spec:
      rules:
      - host: sample.example.com
        http:
          paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: sample
                port:
                  number: 8080
    EOF
    ```

    **NOTE**: Change the hostname to the one you will be using later.

1. You can now access the ingress using the **external-ip** of your ingress controller and providing a **Host** header with the domain you’re using.

    **NOTE**: in the following example, be sure to use the correct hostname that was provided by running the `kubectl get svc/envoy -n projectcontour` command in step 4.3 above.

    ```shell
    curl -H 'Host: sample.example.com' aca8e02fba5614ce0b828f234f4f404f-1196576133.eu-west-1.elb.amazonaws.com
    ```

1. Contour provides an advanced resource type called [HttpProxy](https://projectcontour.io/docs/v1.18.1/config/fundamentals/) that provides [some benefits over Ingress](https://projectcontour.io/docs/v1.18.1/config/fundamentals/#key-httpproxy-benefits) resources.

    You will replace your ingress with the corresponding version of HTTPProxy.

    ```shell
    kubectl delete ingress sample -n sample

    cat <<EOF | kubectl apply -n sample -f -
    apiVersion: projectcontour.io/v1
    kind: HTTPProxy
    metadata:
      name: sample
    spec:
      virtualhost:
        fqdn: sample.example.com
      routes:
        - conditions:
          - prefix: /
          services:
            - name: sample
              port: 8080
    EOF
    ```

    You should see the same output as before. However, this time we are ensuring that a response is received using the proper hostname.

## External-dns

[ExternalDNS](https://github.com/kubernetes-sigs/external-dns) synchronizes exposed Kubernetes Services and Ingresses with DNS providers. What this means is that the external-dns controller will interact with your infrastructure provider (AWS in this case) and will register the DNS name in the DNS service of the infrastructure provider (Route53).

Before you proceed, you need to know that since your cluster will be interacting with the infrastructure provider, you need to set up some authentication/authorization policies for everything to work.

The [External DNS package Readme](https://tanzucommunityedition.io/docs/package-readme-external-dns-0.10.0/) guides you through the process of configuring Route53 and the required Amazon IAM policies. See the following sections in the [External DNS package Readme](https://tanzucommunityedition.io/docs/package-readme-external-dns-0.10.0/):

* Create a permissions policy that allows external DNS updates
* Create an IAM user with the sole permission of updating DNS
* If you don’t have already a Hosted Zone for your domain name, create it now
* Create a Kubernetes secret with your IAM credentials

After having configured Route53 and AWS IAM you can now proceed to install the external-dns package.

Create a file external-dns-config.yaml, which contains:

```shell
---
deployment:
  args:
    - --source=service
    - --domain-filter=example.com
    - --provider=aws
    - --policy=upsert-only
    - --aws-zone-type=public
    - --registry=txt
    - --txt-owner-id=example.com
  env:
    - name: AWS_ACCESS_KEY_ID
      valueFrom:
        secretKeyRef:
          name: external-dns-iam-credentials
          key: access-key-id
    - name: AWS_SECRET_ACCESS_KEY
      valueFrom:
        secretKeyRef:
          name: external-dns-iam-credentials
          key: secret-access-key
```

**NOTE**: Make sure to modify the values with your domain and also with the name of the IAM credentials secret.

Now you can proceed to install the external-dns package:

```shell
tanzu package install external-dns --package-name external-dns.community.tanzu.vmware.com --version 0.8.0 --values-file external-dns-config.yaml
```

The output should look similar to:

```shell
| Installing package 'external-dns.community.tanzu.vmware.com'
| Getting namespace 'default'
| Getting package metadata for 'external-dns.community.tanzu.vmware.com'
| Creating service account 'external-dns-default-sa'
| Creating cluster admin role 'external-dns-default-cluster-role'
| Creating cluster role binding 'external-dns-default-cluster-rolebinding'
| Creating secret 'external-dns-default-values'
- Creating package resource
| Package install status: Reconciling

Added installed package 'external-dns' in namespace 'default'
```

You can now test your application via the generated domain name.

```shell
curl sample.example.com
```

As expected, you should see the same output as before. However, this time we are using the correct domain name to ensure that the ingress controller, external-dns plugin, and the DNS configuration is installed, configured and working correctly.

But security is very important, and all internet communications should happen using a secure transport (HTTPS), so you’re going to add an additional package that will automate using a secure transport for your applications.

## Cert-manager

[Cert-manager](https://cert-manager.io/docs/) automates certificate management in cloud native environments. It provides certificates-as-a-service capabilities. You will install the cert-manager package on your cluster with the provided package, and then you will create some cert-manager resources to provide you with a wildcard certificate for your cluster, so any of your applications will be able to use it. Another option can be to have every application request a certificate individually.

[Cert-manager package docs](https://tanzucommunityedition.io/docs/package-readme-cert-manager-1.6.1/) states that the only possible configuration is the namespace where cert-manager will be installed. Since the default is ok, you will go ahead and install the package without any configuration:

```shell
tanzu package install cert-manager --package-name cert-manager.community.tanzu.vmware.com --version 1.5.1
```

The output should look similar to:

```shell
| Installing package 'cert-manager.community.tanzu.vmware.com'
| Getting namespace 'default'
| Getting package metadata for 'cert-manager.community.tanzu.vmware.com'
| Creating service account 'cert-manager-default-sa'
| Creating cluster admin role 'cert-manager-default-cluster-role'
| Creating cluster role binding 'cert-manager-default-cluster-rolebinding'
- Creating package resource
| Package install status: Reconciling


Added installed package 'cert-manager' in namespace 'default'
```

Now that cert-manager is installed, you need to create a Certificate Issuer that will create the certificates you need. As we will use Let’s Encrypt, you need a [certificate issuer of type ACME](https://cert-manager.io/docs/configuration/acme/). Since you want a wildcard certificate, you need to use a [DNS01 solver](https://cert-manager.io/docs/configuration/acme/dns01/route53/), which will integrate with Route53 to prove Let’s Encrypt that you own that domain. Cert-manager will create some records on your domain’s hosted zone and during the process, Let’s Encrypt will query to validate the ownership of the domain. But luckily all this is done by cert-manager, so you will only need to create a couple of Kubernetes objects. But first, you will need to create an IAM user with a policy that will allow cert-manager to do what it needs.

For convenience, you will modify the [policy created for external-dns](../package-readme-external-dns-0.8.0/#1-aws-permissions) so that it can also be used by cert-manager.

In your aws console, find the policy you created before and add the rules required by cert-manager, which are:

```shell
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "route53:GetChange",
      "Resource": "arn:aws:route53:::change/*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "route53:ChangeResourceRecordSets",
        "route53:ListResourceRecordSets"
      ],
      "Resource": "arn:aws:route53:::hostedzone/*"
    },
    {
      "Effect": "Allow",
      "Action": "route53:ListHostedZonesByName",
      "Resource": "*"
    }
  ]
}
```

Your external-dns policy should now look like this:

```shell
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "route53:GetChange",
            "Resource": "arn:aws:route53:::change/*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "route53:ChangeResourceRecordSets",
                "route53:ListResourceRecordSets"
            ],
            "Resource": "arn:aws:route53:::hostedzone/*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "route53:ListHostedZones",
                "route53:ListResourceRecordSets",
                "route53:ListHostedZonesByName"
            ],
            "Resource": "*"
        }
    ]
}
```

Now, you can create a ClusterIssuer resource that will instruct cert-manager how to communicate with Let’s Encrypt to generate certificates. Since cert-manager will also need to access the secret with the IAM user credentials, deploy it to the cert-manager namespace as well. Follow the [same process as before](../package-readme-external-dns-0.8.0/#4-create-a-kubernetes-secret) but using cert-manager namespace.

Now, create the cluster issuer, making sure to adapt to your own values:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: wildcard-lets-encrypt
spec:
  acme:
    email: jorge@none.com
    privateKeySecretRef:
      name: aws-creds
    server: https://acme-v02.api.letsencrypt.org/directory
    solvers:
    - dns01:
        route53:
          region: eu-west1
          accessKeyID: AKIARKM7D2FFCPM7MXM5
          secretAccessKeySecretRef:
            name: aws-creds
            key: secret-access-key
    - http01:
        ingress:
          class: contour
EOF
```

**NOTE**: we don't recommend this means of propagating IAM secret data in production clusters, you should use a more secure solution like [Hashicorp Vault](https://www.vaultproject.io/), [kubernetes-external-secrets](https://github.com/external-secrets/kubernetes-external-secrets) or [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/).

This cluster issuer can be used for wildcard certificates (DNS01) as well as regular ingress certificates (HTTP01). In this guide you will only use the former.

Now that the Cluster Issuer is installed, let’s go ahead and request a wildcard certificate for your domain:

```shell
cat <<EOF | kubectl apply -n tanzu-certificates -f -
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: wildcard
spec:
  secretName: wildcard
  issuerRef:
    name: wildcard-lets-encrypt
    kind: ClusterIssuer
  dnsNames:
  - "*.example.com"
EOF
```

This will create a secret named wildcard in the tanzu-certificates namespace, that will hold your wildcard certificate. Note that this is the namespace where cert-manager is installed by default.

A certificate request will be issued, and you can check the process by monitoring the **Certificate** resource you just created.

```shell
kubectl get certs -n tanzu-certificates
```

Once the certificate is there, you should see something like:

```shell
NAME       READY   SECRET     AGE
wildcard   True    wildcard   0m20s
```

Now, you can modify the sample application ingress, to use your wildcard certificate. But before you do that, you need to instruct contour to provide the wildcard certificate to any ingress, using a feature called [TLSDelegation](https://projectcontour.io/docs/v1.18.1/config/tls-delegation/):

```shell
cat <<EOF | kubectl apply -n tanzu-certificates -f -
apiVersion: projectcontour.io/v1
kind: TLSCertificateDelegation
metadata:
  name: wildcard
spec:
  delegations:
  - secretName: wildcard
    targetNamespaces: ["*"]
```

Now, you can apply this updated version of the HTTPProxy (or Ingress) that uses TLS:

```shell
cat <<EOF | kubectl apply -n sample -f -
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: sample
spec:
  virtualhost:
    fqdn: sample.example.com
    tls:
      secretName: tanzu-certificates/wildcard
  routes:
    - conditions:
      - prefix: /
      services:
        - name: sample
          port: 8080
EOF
```

**NOTE**: If you're using Ingress, TLS Delegation [no longer works with the latest of the networking.k8s.io apiGroup (v1)](https://github.com/projectcontour/contour/issues/3544). You need to use previous version of the API (v1beta1), or use the more advanced HTTPProxy.

You can now verify that your application works. Open a browser and type [**https://sample.example.com**](https://sample.example.com). You should see the same result as before.

**NOTE**: Make sure you use your domain name in the previous URL.

Click the lock icon on the left of the address bar to verify that the certificate used is the wildcard certificate that was issued.

![Secure Connection](../img/secure-connection.png)

## Summary

This guide has walked you through the process of configuring your environment to have your services served via a domain name in a secure way.
