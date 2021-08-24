# Implementing Service Discovery with External DNS

The external DNS service reserves DNS hostnames for applications, using a declarative, Kubernetes-native interface.
It is packaged as an extension in the Tanzu Kubernetes Grid extensions bundle

This topic explains how to deploy the external DNS service to a workload or shared services cluster in Tanzu Kubernetes Grid.

On infrastructures with load balancing (AWS, Azure, and vSphere with NSX Advanced Load Balancer), VMware recommends installing the External DNS service alongside the Harbor service, as described in [Harbor Registry and External DNS](harbor-registry.md#external-dns), especially in production or other environments where Harbor availability is important.

The procedures in this topic apply to vSphere, Amazon EC2, and Azure deployments.

## <a id="prereqs"></a> Prerequisites

- You have deployed a management cluster on vSphere, Amazon EC2, or Azure, in either an Internet-connected or Internet-restricted environment.

   If you are using Tanzu Kubernetes Grid in an Internet-restricted environment, you performed the procedure in [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](../mgmt-clusters/airgapped-environments.md) before you deployed the management cluster.
- You have downloaded and unpacked the bundle of Tanzu Kubernetes Grid extensions. For information about where to obtain the bundle, see [Download and Unpack the Tanzu Kubernetes Grid Extensions Bundle](index.md#unpack-bundle).  
- You have installed the Carvel tools. For information about installing the Carvel tools, see [Install the Carvel Tools](../install-cli.md#install-carvel).
- Determine the FQDN for the external DNS service. If you are using external DNS for Harbor on a shared services cluster upgraded from Tanzu Kubernetes Grid v1.2, note the following:

   - If your Harbor registry in v1.2 used a fully-qualified domain name (FQDN) that you control, such as `myharbor.mycompany.com`, use this FQDN for the external DNS service.
   - If your Harbor registry in v1.2 used a fictitious domain name such as `harbor.system.tanzu`, you cannot upgrade workload clusters automatically, but must instead create a new v1.3 workload cluster and migrate the workloads to the new cluster manually.

**IMPORTANT**: The extensions folder `tkg-extensions-v1.3.x+vmware.1` contains subfolders for each type of extension, for example, `authentication`, `ingress`, `registry`, and so on. At the top level of the folder there is an additional subfolder named `extensions`. The `extensions` folder also contains subfolders for `authentication`, `ingress`, `registry`, and so on. Take care to run commands from the location provided in the instructions. Commands are usually run from within the `extensions` folder.

## <a id="prepare-tkc"></a> Prepare a Cluster for External DNS Deployment

To prepare a cluster for running the External DNS service:

1. If you are running External DNS on a shared services cluster, alongside Harbor, and the cluster has not yet been created, create the cluster by following [Create a Shared Services Cluster](index.md#shared).

1. Deploy the Contour service on the cluster.
External DNS requires Contour to be present on its cluster, to provide ingress control. For how to deploy Contour, see [Deploy Contour on the Tanzu Kubernetes Cluster](ingress-contour.md#deploy).

## <a id="provider"></a>Choose the External DNS Provider

The external-dns extension has been validated with AWS (Route53), Azure,
and RFC2136 (BIND). The extension supports Ingress with either Contour
or Service type Load Balancer. Below are instructions for each of these
options.

### <a id="aws"></a>AWS (Route53)

1. Create a hosted zone within Route53 with the domain that shared
    services will be using.
1. Take note of the “Hosted zone ID” as this ID will be used in the
    external-dns-data-values.yaml.
1. Create an IAM user for external-dns with the following policy
    document and ensure “Programmatic access” is checked. If desired you
    may fine-tune the policy to permit updates only to the hosted zone
    that you just created.

    ```
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "route53:ChangeResourceRecordSets"
          ],
          "Resource": [
            "arn:aws:route53:::hostedzone/*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "route53:ListHostedZones",
            "route53:ListResourceRecordSets"
          ],
          "Resource": [
            "*"
          ]
        }
      ]
    }

  ```
  
1.  Take note of the “Access key ID” and “Secret access key” as they
    will be needed for configuring external-dns.
1.  Create a Kubernetes secret to supply credentials for route53:

  ```

  kubectl -n tanzu-system-service-discovery create secret generic
  route53-credentials
  --from-literal=aws_access_key_id=YOUR_ACCESS_KEY_ID_HERE
  --from-literal=aws_secret_access_key=YOUR_SECRET_ACCESS_KEY_HERE

  ```
1.  Copy the appropriate example values.
    - **With Contour**: If you have deployed Contour and would like the
      external-dns extension to use Contour HTTPProxy resources as
      sources for external-dns, copy the example file with:
        ```
        cp external-dns-data-values-aws-with-contour.yaml.example external-dns-data-values.yaml
        ```
    - **Without Contour**: To use a Service type LoadBalancer as the
      only source for external-dns, copy the example file with:
        ```
        cp external-dns-data-values-aws.yaml.example external-dns-data-values.yaml
        ```

1.  Change the values inside `external-dns-data-values.yaml` as
    appropriate, making sure to fill in your hosted zone ID and domain
    name.
     For additional configuration options, you can use `ytt` overlays as described in [Extensions and Shared Services](../ytt.md#extensions) in _Customizing Clusters, Plans, and Extensions with ytt Overlays_ and in the extensions mods examples in the [TKG Lab repository](https://github.com/Tanzu-Solutions-Engineering/tkg-lab).

### <a id="rfc2136"></a>RFC2136 (BIND) Server

The RFC2136 provider allows you to use any RFC2136-compatible DNS server
as a provider for external-dns such as BIND.

1. Find/Create a TSIG key for your server
    
  1. If your DNS is provided for you, ask for a TSIG key authorized
     to update and transfer the zone you wish to update. The key
     will look something like this:
     ```
      key "externaldns-key" {
         algorithm hmac-sha256;
         secret "/2avn5M4ndEztbDqy66lfQ+PjRZta9UXLtToW6NV5nM=";
      };
     ```

  1. If you are managing your own DNS server then you can create a
     TSIG key using `tsig-keygen -a hmac-sha256 externaldns`. Copy
     the result to your DNS servers configuration. In the case of
     BIND you would add the key to the named.conf file and
     configure the zone with the allow-transfer and update policy
     fields. For example:
     ```
      key "externaldns-key" {
         algorithm hmac-sha256;
         secret "/2avn5M4ndEztbDqy66lfQ+PjRZta9UXLtToW6NV5nM=";
      };
      zone "k8s.example.org" {
         type master;
         file "/etc/bind/zones/k8s.zone";
            allow-transfer {
              key "externaldns-key";
            };
            update-policy {
              grant externaldns-key zonesub ANY;
            };
       };
      ```

  1. The above assumes you also have a zone file that might look
     something like this:
     ```
      $TTL 60 ; 1 minute
      @         IN SOA  k8s.example.org.  root.k8s.example.org. (
                        16  ; serial
                        60  ; refresh (1 minute)
                        60  ; retry (1 minute)
                        60  ; expire (1 minute)
                        60  ; minimum (1 minute)
                        )
                   NS   ns.k8s.example.org.
      ns           A    1.2.3.4
      ```
1. Copy the appropriate example values.
    
    - **With Contour**: If you have deployed Contour and would like the
      external-dns extension to use Contour HTTPProxy resources as
      sources for external-dns, copy the example file with:
        ```
        cp external-dns-data-values-rfc2136-with-contour.yaml.example external-dns-data-values.yaml
        ```
    - **Without Contour**: To use a Service type LoadBalancer as the
      only source for external-dns, copy the example file with:
        ```
        cp external-dns-data-values-rfc2136.yaml.example external-dns-data-values.yaml
        ```

1. Change the values inside `external-dns-data-values.yaml` as
   appropriate, making sure to fill in your DNS server IP, domain name,
   TSIG secret, and TSIG key name.

### <a id="Azure"></a>Microsoft Azure

1. Log in to the `az` CLI: `az login`

1. Set your subscription: `az account set -s <subscriptionId GUID>`

1. Create a service principal: `az ad sp create-for-rbac -n <name of service principal>`
  - This returns a json that looks similar to the following:
    ```
    {
      "appId": "a72a7cfd-7cb0-4b02-b130-03ee87e6ca89",
      "displayName": "foo",
      "name": "http://foo",
      "password": "515c55da-f909-4e17-9f52-236ffe1d3033",
      "tenant": "b35138ca-3ced-4b4a-14d6-cd83d9ea62f0"
    }
    ```

1. Assign permissions to the service principal.
  1. Discover the id of the resource group
    ```
    az group show --name <resource group> --query i
    ```
  1. Assign the reader role to the service principal for the
     resource group scope. You will need the appId from the output
     of the creation of the service principal.
    ```
    az role assignment create --role "Reader" --assignee <appId GUID> --scope <resource group resource id>
    ```
  1. Discover the id of the dns zone.
    ```
    az network dns zone show --name <dns zone name> -g <resource group name> --query i
    ```
  1. Assign the contributor role to the service principal for the
     dns zone scope.
    ```
    az role assignment create --role "Contributor" --assignee <appId GUID> --scope <dns zone resource id>
    ```

1. To connect the external-dns extension to the Azure DNS service you
   will create a configuration file called azure.json on your local
   machine with contents that look like the following:
   ```

   {
     "tenantId": "01234abc-de56-ff78-abc1-234567890def",
     "subscriptionId": "01234abc-de56-ff78-abc1-234567890def",
     "resourceGroup": "MyDnsResourceGroup",
     "aadClientId": "01234abc-de56-ff78-abc1-234567890def",
     "aadClientSecret": "uKiuXeiwui4jo9quae9o"
   }

  ```
    - The `tenantId` can be retrieved from: `az account show --query
        "tenantId"`
    - The `subscriptionId` can be retrieved from: `az account show --query
        "id"`
    - The `resourceGroup` is the name of the resource group that your dns
        zone is within.
    - The `aadClientId` is the `appId` from the output of the Service
        Principal.
    - The `aadClientSecret` is the password from the output of the Service
        Principal.

1. Create a Kubernetes secret with the `azure.json` configuration file.
  ```

  kubectl -n tanzu-system-service-discovery create secret generic
  azure-config-file --from-file=azure.json

  ```

1. Copy the appropriate example values.
    - **With Contour**:  If you have deployed Contour and would like the
      external-dns extension to use Contour HTTPProxy resources as
      sources for external-dns, copy the example file with:
      ```
      cp external-dns-data-values-azure-with-contour.yaml.example external-dns-data-values.yaml
      ```
    - **Without Contour**: To use a Service type LoadBalancer as the
      only source for external-dns, copy the example file with:
      ```
      cp external-dns-data-values-azure.yaml.example external-dns-data-values.yaml
      ```

1.  Change the values inside `external-dns-data-values.yaml` as
    appropriate, making sure to fill in the Azure resource group and
    domain name.

## <a id="deploy"></a>Deploy the External DNS Extension

1. Set the context of `kubectl` to the shared services cluster or other cluster where you are deploying External DNS.

  ```

  kubectl config use-context tkg-services-admin@tkg-services

  ```

1. From the unpacked `tkg-extensions` folder, navigate to the
    `external-dns` extension folder
  ```

  cd extensions/service-discovery/external-dns

  ```

1. Install kapp-controller:
  ```

  kubectl apply -f ../../kapp-controller.yaml  

  ```

1. Create external-dns namespace:
  ```

  kubectl apply -f namespace-role.yaml

  ```

1. Create a secret with data values:
  ```

  kubectl create secret generic external-dns-data-values
  --from-file=values.yaml=external-dns-data-values.yaml -n
  tanzu-system-service-discovery

  ```

1.  Deploy the `ExternalDNS` extension:
  ```

  kubectl apply -f external-dns-extension.yaml

  ```

1.  Ensure the extension is deployed successfully:
  ```

  kubectl get app external-dns -n tanzu-system-service-discovery  

  ```
  `ExternalDNS` app status should change to **Reconcile Succeeded** once
  `ExternalDNS` is deployed successfully.

## <a id="validate"></a>Validating External DNS

If configured with Contour, External DNS will automatically watch the
specified namespace for HTTPProxy resources and create DNS records for
services with hostnames that match the configured domain filter.

External DNS will also automatically watch for Kubernetes Services with
the annotation `external-dns.alpha.kubernetes.io/hostname` and create DNS records for
services whose annotations match the configured domain filter.

For example, a service with the annotation  
`external-dns.alpha.kubernetes.io/hostname: foo.k8s.example.org`

will cause External DNS to create a DNS record for `foo.k8s.example.org`,
and you can validate that the record exists by examining the zone that
you created in whichever provider you created.
