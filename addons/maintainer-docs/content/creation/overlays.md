---
title: ytt Overlays
weight: 3
---

Overlays allow you to apply a custom configuration on top of the upstream source of the package. For each configurable item in your schema, there should be a corresponding operation in a template to apply that parameter.

## Example Usage

Instructions for how to write overlays are outside the scope of this guide. Full tutorials and examples are available at the [ytt](https://carvel.dev/ytt/) site. However, here is an example of using overlays to modify the namespace that cert-manager uses.

1. Create the following three overlay files to modify the various places in the cert-manager manifest where the namespace is referenced. This example could use just one overlay, but it's convenient to have things separated at times.

    ```shell
    cat > bundle/config/overlays/annotations.yaml <<EOF
    #@ load("@ytt:data", "data")
    #@ load("@ytt:overlay", "overlay")

    #@overlay/match by=overlay.subset({"kind":"CustomResourceDefinition"}), expects=6
    ---
    metadata:
      annotations:
        #@overlay/match missing_ok=True
        cert-manager.io/inject-ca-from-secret: #@ "{}/cert-manager-webhook-ca".format(data.values.namespace)

    #@overlay/match by=overlay.subset({"kind":"MutatingWebhookConfiguration"})
    ---
    metadata:
      annotations:
        #@overlay/match missing_ok=True
        cert-manager.io/inject-ca-from-secret: #@ "{}/cert-manager-webhook-ca".format(data.values.namespace)

    #@overlay/match by=overlay.subset({"kind":"ValidatingWebhookConfiguration"})
    ---
    metadata:
      annotations:
        #@overlay/match missing_ok=True
        cert-manager.io/inject-ca-from-secret: #@ "{}/cert-manager-webhook-ca".format(data.values.namespace)
    EOF
    ```

    ```shell
    cat > bundle/config/overlays/deployment.yaml <<EOF
    #@ load("@ytt:overlay", "overlay")
    #@ load("@ytt:data", "data")

    #@overlay/match by=overlay.subset({"kind": "Deployment", "metadata": {"name": "cert-manager-webhook"}})
    ---
    spec:
      template:
        spec:
          containers:
          #@overlay/match by="name"
          - name: cert-manager
            args:
              #@overlay/match by=lambda i,l,r: l.startswith("--dynamic-serving-dns-names=")
              - #@ "--dynamic-serving-dns-names=cert-manager-webhook,cert-manager-webhook.{},cert-manager-webhook.{}.svc".format(data.values.namespace, data.values.namespace)
    EOF
    ```

    ```shell
    cat > bundle/config/overlays/misc.yaml <<EOF
    #@ load("@ytt:data", "data")
    #@ load("@ytt:overlay", "overlay")

    #@overlay/match by=overlay.subset({"kind":"Namespace", "metadata": {"name": "cert-manager"}})
    ---
    apiVersion: v1
    kind: Namespace
    metadata:
      name: #@ data.values.namespace

    #@overlay/match by=overlay.subset({"metadata": {"namespace": "cert-manager"}}), expects=10
    ---
    metadata:
      namespace: #@ data.values.namespace

    #@ crb=overlay.subset({"kind":"ClusterRoleBinding"})
    #@ rb=overlay.subset({"kind":"RoleBinding"})
    #@overlay/match by=overlay.or_op(crb, rb), expects=13
    ---
    subjects:
    #@overlay/match by=overlay.subset({"namespace": "cert-manager"})
    - kind: ServiceAccount
      namespace: #@ data.values.namespace

    #@ vwc=overlay.subset({"kind":"ValidatingWebhookConfiguration"})
    #@ mwc=overlay.subset({"kind":"MutatingWebhookConfiguration"})
    #@overlay/match by=overlay.or_op(vwc, mwc), expects=2
    ---
    webhooks:
    #@overlay/match by="name"
    - name: webhook.cert-manager.io
      clientConfig:
        service:
          namespace: #@ data.values.namespace

    #@overlay/match by=overlay.subset({"kind":"CustomResourceDefinition"}), expects=6
    ---
    spec:
      conversion:
        webhook:
          clientConfig:
            #@overlay/match by="name"
            service:
              name: cert-manager-webhook
              namespace: #@ data.values.namespace
    EOF
    ```

2. One more file is required to hold configuration values. In this case, the only value that we can modify is the namespace, so we provide a data value for the namespace. The configuration parameters defined in this file will later be documented in the package CRD.

    ```shell
    cat > bundle/config/values.yaml <<EOF
    #@data/values
    ---

    #! The namespace in which to deploy cert-manager.
    namespace: custom-namespace
    EOF
    ```

3. To test if everything is working, run ytt. If everything is correct, ytt will output the transformed YAML. If there's a problem, you'll see it in the console.

    ```shell
    ytt --file bundle/config
    ```
