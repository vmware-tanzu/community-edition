# THIS CONTENT HAS MOVED TO THE DOCS BRANCH IN:  PLEASE MAKE ANY FURTHER UPDATES THERE

File is available here on docs branch: ``docs\site\content\docs\latest\gatekeeper-config``

## Gatekeeper

This package provides custom admission control using
[gatekeeper](https://github.com/open-policy-agent/gatekeeper). Under the hood,
gatekeeper uses [Open Policy Agent](https://www.openpolicyagent.org) to enforce
policy when requests hit the Kubernetes API server.

## Components

* gatekeeper: Uses Open Policy Agent (OPA) to validate whether a request is
authorized.
* audit-controller: Identifies existing resources in the cluster that break
active policy.

## Configuration

The following configuration values can be set to customize the gatekeeper installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy gatekeeper. |

### gatekeeper Configuration

_Currently there is no gatekeeper customization_.

## Usage Example

This walkthrough demonstrates how to apply a policy that restricts root users
from running containers. The gatekeeper project maintains a library of policies
at
[github.com/open-policy-agent/gatekeeper-library](https://github.com/open-policy-agent/gatekeeper-library).
This walkthrough will leverage a policy from this repository.

1. Apply the following constraint template, which check for specified labels.

    ```yaml
    apiVersion: templates.gatekeeper.sh/v1beta1
    kind: ConstraintTemplate
    metadata:
      name: k8srequiredlabels
      annotations:
        description: Requires all resources to contain a specified label with a value
          matching a provided regular expression.
    spec:
      crd:
        spec:
          names:
            kind: K8sRequiredLabels
          validation:
            # Schema for the `parameters` field
            openAPIV3Schema:
              properties:
                message:
                  type: string
                labels:
                  type: array
                  items:
                    type: object
                    properties:
                      key:
                        type: string
                      allowedRegex:
                        type: string
      targets:
        - target: admission.k8s.gatekeeper.sh
          rego: |
            package k8srequiredlabels
            get_message(parameters, _default) = msg {
              not parameters.message
              msg := _default
            }
            get_message(parameters, _default) = msg {
              msg := parameters.message
            }
            violation[{"msg": msg, "details": {"missing_labels": missing}}] {
              provided := {label | input.review.object.metadata.labels[label]}
              required := {label | label := input.parameters.labels[_].key}
              missing := required - provided
              count(missing) > 0
              def_msg := sprintf("you must provide labels: %v", [missing])
              msg := get_message(input.parameters, def_msg)
            }
            violation[{"msg": msg}] {
              value := input.review.object.metadata.labels[key]
              expected := input.parameters.labels[_]
              expected.key == key
              # do not match if allowedRegex is not defined, or is an empty string
              expected.allowedRegex != ""
              not re_match(expected.allowedRegex, value)
              def_msg := sprintf("Label <%v: %v> does not satisfy allowed regex: %v", [key, value, expected.allowedRegex])
              msg := get_message(input.parameters, def_msg)
            }
    ```

1. Verify the `k8srequiredlabels` CRD was created.

    ```sh
    kubectl get crds | grep -i k8srequiredlabels
    ```

1. Create a constraint that requires the label `owner` to be specified.

    ```yaml
    apiVersion: constraints.gatekeeper.sh/v1beta1
    kind: K8sRequiredLabels
    metadata:
      name: all-must-have-owner
    spec:
      match:
        kinds:
          - apiGroups: [""]
            kinds: ["Namespace"]
      parameters:
        message: "All namespaces must have an `owner` label"
        labels:
          - key: owner
    ```

1. Create a namespace

    ```sh
    kubectl create ns test
    ```

1. Verify it fails to deploy due to missing label.

    ```text
    Error from server ([denied by all-must-have-owner] All namespaces must have an `owner` label): admission webhook "validation.gatekeeper.sh" denied the request: [denied by all-must-have-owner] All namespaces must have an `owner` label
    ```

1. Create a namespace with the owner label.

    ```yaml
    apiVersion: v1
    kind: Namespace
    metadata:
      name: test
      labels:
        owner: bearcanoe
    ```
