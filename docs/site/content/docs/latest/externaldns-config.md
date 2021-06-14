# Configuring the ExternalDNS package

[ExternalDNS](https://github.com/kubernetes-sigs/external-dns) synchronizes exposed Kubernetes Services and Ingresses with DNS providers.

## Components

* ExternalDNS deployment

## Configuration

The following configuration values can be set to customize the ExternalDNS installation.

### Global

| Value       | Required/Optional | Description                                    |
|:-------------|:-------------------|:------------------------------------------------|
| `namespace` | Optional          | The namespace in which to deploy ExternalDNS. |

### ExternalDNS Configuration

| Value                        | Required/Optional&nbsp;&nbsp; | Description                                       |
|:------------------------------|:--------------------|:--------------------------------------------------|
| `deployment.args`            | Required           | Args passed via command-line to ExternalDNS     |
| `deployment.env`             | Optional           | Environment variables to pass to ExternalDNS    |
| `deployment.securityContext` | Optional           | Security context of the ExternalDNS container   |
| `deployment.volumeMounts`    | Optional           | Volume mounts of the ExternalDNS container      |
| `deployment.volumes`         | Optional           | Volumes of the ExternalDNS pod                  |

Follow [the ExternalDNS docs](https://github.com/kubernetes-sigs/external-dns#running-externaldns)
for guidance on how to configure ExternalDNS for your DNS provider.

### Configuration sample
The following example shows a simple `bind` (rfc2136) implementation. 

#### Before you begin
Ensure the ExternalDNS package is installed, for more information about installing packages, see [Packages Introduction](packages-intro.md).

#### Procedure
1. Run the following command to generate an empty configuration file in the current directory:

    `tanzu package configure external-dns.tce.vmware.com`

2. Update the empty configuration file based on the following sample:  
Note: Comments which begin with `#@` are important `ytt` directives and should remain unchanged in
your final configuration file.

    ```yaml
    #@data/values
    #@overlay/match-child-defaults missing_ok=True
    ---

    #! The namespace in which to deploy ExternalDNS.
    namespace: external-dns

    #! Deployment related configuration
    deployment:
      #@overlay/replace
      args:
      - --source=service
      - --source=contour-httpproxy
      - --txt-owner-id=k8s
      - --domain-filter=k8s.example.org
      - --namespace=tanzu-system-service-discovery
      - --provider=rfc2136
      - --rfc2136-host=100.69.97.77
      - --rfc2136-port=53
      - --rfc2136-zone=k8s.example.org
      - --rfc2136-tsig-secret=MTlQs3NNU=
      - --rfc2136-tsig-secret-alg=hmac-sha256
      - --rfc2136-tsig-keyname=externaldns-key
      - --rfc2136-tsig-axfr
      env: []
      securityContext: []
      volumeMounts: []
      volumes: []
    ```

### Configuring with Contour HTTPProxy

Follow [this tutorial](https://github.com/kubernetes-sigs/external-dns/blob/v0.7.6/docs/tutorials/contour.md) for guidance on providing arguments to ExternalDNS to enable HTTPProxy support. The ExternalDNS package is preconfigured with the correct RBAC permissions to watch for HTTPProxies, so this part of the tutorial may be skipped.

## Amazon Web Services Route 53 Example

This example guides you through setting up the ExternalDNS package with the AWS Route 53 DNS service. This example is based on the instructions for Setting Up ExternalDNS for Services on AWS in the [ExternalDNS documentation](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md).

### Before you begin 
Ensure you have met the following prerequisites:

* Your Cluster is on AWS
* You have a Domain managed by Route 53
* You have the ability to create AWS IAM users and permissions

### 1. AWS Permissions

Start by creating a permissions policy that allows external DNS updates. 

1. In the [AWS Console](https://console.aws.amazon.com/iam/home#/policies$new?step=edit), select the JSON tab, and paste in the following policy. For more information, see the [ExternalDNS documentation](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md#iam-policy).

    ```json
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

    Note: This policy allows updating of any hosted zone. You can limit the zones effected by replacing the wildcard with the hosted zone you will be using for this example.

    ![Create Policy Step 1](/docs/img/create-policy-step1.png)

2. Continue through the wizard and complete the policy. For simplicity, name the policy as the documentation suggests, as `AllowExternalDNSUpdates` and create the policy.

    ![Create Policy Step 2](/docs/img/create-policy-step2.png)

### 2. IAM User in AWS

1. Create an IAM user in the [AWS console](https://console.aws.amazon.com/iam/home#/users$new?step=details) called `external-dns-user`. This user will have the sole permission for updating DNS.
For Access Type, select Programmatic access.

    ![Create User Step 1](/docs/img/create-user-step1.png)

2. Attach the `AllowExternalDNSUpdates` policy to the new user created in the previous step. Select `Attach existing policies directly` and search for and then select the policy.

    ![Create User Step 2](/docs/img/create-user-step2.png)

3. Continue to the review page, review your choices, and select Create user.

    ![Create User Step 3](/docs/img/create-user-step3.png)

4. Copy the access keys. These credentials will be used to give ExternalDNS access to this user and permission to modify your DNS settings. This will be your only opportunity to see the `secret-access-key`. Make a note of the Access Key ID and Secret access key.

    ![Create User Step 4](/docs/img/create-user-step4.png)

### 3. Hosted Zone

You can follow the instructions in the [ExternalDNS documentation](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md#set-up-a-hosted-zone), or alternatively, follow these steps. The ExternalDNS documentation creates a subdomain on the hosted zone. You can do this, or just use the hosted zone itself. There is an extra step if you choose the subdomain route that is not reflected in the ExternalDNS documentation. This example will follow the ExternalDNS and call out the additional step.

For this example, we are using the domain `k8squid.com`, and a subdomain of `external-dns-test`. 

1. Create the new hosted zone.

    ```shell
    aws route53 create-hosted-zone --name "external-dns-test.k8squid.com." --caller-reference "external-dns-test-$(date +%s)"
    /hostedzone/Z09346372A26K4C7GYTEI
    ```

2. Obtain the name servers assigned to the new subdomain.

    ```shell
    aws route53 list-resource-record-sets --output json --hosted-zone-id "/hostedzone/Z09346372A26K4C7GYTEI" --query "ResourceRecordSets[?Type == 'NS']" | jq -r '.[0].ResourceRecords[].Value'
    ns-451.awsdns-56.com.
    ns-1214.awsdns-23.org.
    ns-1625.awsdns-11.co.uk.
    ns-515.awsdns-00.net.
    ```

    Take note of the new hosted zone id and name servers.

3. "Hook up your DNS zone with is parent zone", as the official documentation cryptically suggests. Go to the [AWS Route 53 Console](https://console.aws.amazon.com/route53/v2/hostedzones#) and select your domain. Create a new record. Enter the desired subdomain, select NS for the record type, and paste in the list of name servers from the previous step into the Value field.

    ![Create NS Record](/docs/img/create-ns-record.png)

    After creating the NS record on the hosted zone for your new subdomain, you've completed the prerequisites on AWS for this example.

### 4. Create a Kubernetes Secret

In an earlier section, you obtained AWS credentials. Use these credentials to make a secret in Kubernetes that ExternalDNS can reference. Start by creating a manifest for an opaque secret.

The secret must be created in the same namespace that the ExternalDNS package will run it. 

1. If the namespace does not exist, create it now and use it in the manifest below.

    ```shell
    kubectl create namespace my-external-dns
    ```

2. Create the manifest and apply it to your cluster: `kubectl apply -f secret.yaml`.
For the secret, you will need: the secret name, the namespace, and the AWS access key ID and Secret access key.

    ```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: << SECRET CREDENTIAL NAME >>
      namespace: << NAMESPACE >>
    type: Opaque
    stringData:
      access-key-id: << ACCESS KEY ID >>
      secret-access-key: << SECRET ACCESS KEY >>
    ```

### 5. Install the ExternalDNS package

Next you will, configure the ExternalDNS package to use your new AWS hosted zone. 

1. Run the following command to obtain the configuration file.

    ```shell
    tanzu package configure external-dns.tce.vmware.com
    ```

2. Edit the configuration file and provide the values to configure ExternalDNS with the Route 53 provider. In this example, provide the values for:

  * DOMAIN, e.g. `example.com`
  * HOSTED ZONE ID, e.g. `Z09346372A26K4C7GYTEI`
  * SECRET CREDENTIAL NAME, e.g whatever name was used in step 4.

    ```yaml
    #@data/values
    #@overlay/match-child-defaults missing_ok=True
    ---

    #! The namespace in which to deploy ExternalDNS.
    namespace: external-dns

    #! Deployment related configuration
    #@overlay/replace
    deployment:
      args:
        - --source=service
        - --source=ingress
        - --domain-filter=external-dns-test.<< DOMAIN >> # will make ExternalDNS see only the hosted zones matching provided domain, omit to process all available hosted zones
        - --provider=aws
        - --policy=upsert-only # would prevent ExternalDNS from deleting any records, omit to enable full synchronization
        - --aws-zone-type=public # only look at public hosted zones (valid values are public, private or no value for both)
        - --registry=txt
        - --txt-owner-id=<< HOSTED ZONE ID >>
      env:
        - name: AWS_ACCESS_KEY_ID
          valueFrom:
            secretKeyRef:
              name: << SECRET CREDENTIAL NAME >>
              key: access-key-id
        - name: AWS_SECRET_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: << SECRET CREDENTIAL NAME >>
              key: secret-access-key
      securityContext: []
      volumeMounts: []
      volumes: []
    ```

3. Once the configuration file is updated with your information, deploy the ExternalDNS package to your cluster.

    ```shell
    tanzu package install external-dns.tce.vmware.com --config external-dns.tce.vmware.com-values.yaml
    ```

    After a minute or so, check to see that the package has installed.

    ```shell
    kubectl get apps --all-namespaces
    NAMESPACE         NAME                             DESCRIPTION           SINCE-DEPLOY   AGE
    external-dns      external-dns.tce.vmware.com      Reconcile succeeded   26s            26s
    ```

### Result
ExternalDNS should now be installed and running on your cluster. To verify that it works, you can follow the example in the [official documentation using a service](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md#verify-externaldns-works-service-example). Be sure to substitute your domain name and hosted zone id in service manifest and relevant AWS CLI commands.

For more advanced use cases and documentation, see the official ExternalDNS [documentation](https://github.com/kubernetes-sigs/external-dns).