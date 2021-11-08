# Build Harbor Package

## How to generate the manifests from Harbor Helm Chart

The manifests from 02 to 11 are generated from [Harbor Helm Chart](https://github.com/goharbor/harbor-helm). Helm CLI 3 and [yq](https://github.com/mikefarah/yq) are required to generate these manifests.

1. Clone the harbor-helm repo.

    ```sh
    git clone https://github.com/goharbor/harbor-helm.git

    git checkout v1.7.3
    ```

2. Create a script named `generate-manifests.sh` in the harbor-helm directory.

    ```shell
    #!/usr/bin/env bash
    ## Note yq (https://github.com/mikefarah/yq) v4 is required

    rm -rf manifests
    mkdir -p manifests

    sed -i 's/"%s-%s" .Release.Name $name/"harbor"/' templates/_helpers.tpl
    sed -i '/heritage:/d' templates/_helpers.tpl
    sed -i '/release:/d' templates/_helpers.tpl
    sed -i '/chart:/d' templates/_helpers.tpl

    valuesFile=$(mktemp /tmp/values.XXXXXX.yaml)
    release=harbor

    cat <<EOF >> $valuesFile
    expose:
      tls:
        certSource: secret
        secret:
          secretName: harbor-tls
          notarySecretName: harbor-tls
    persistence:
      resourcePolicy: remove
    caSecretName: harbor-tls
    chartmuseum:
      enabled: false
    secretKey: -the-secret-key-
    core:
      secret: the-secret-of-the-core
      xsrfKey: -xsrf-key-must-be-32-characters-
      secretName: harbor-token-service
    jobservice:
      secret: the-secret-of-the-jobservice
    notary:
      secretName: harbor-notary-signer
    registry:
      secret: the-secret-of-the-registry
    internalTLS:
      enabled: true
      certSource: secret
      core:
        secretName: harbor-core-internal-tls
      jobservice:
        secretName: harbor-jobservice-internal-tls
      registry:
        secretName: harbor-registry-internal-tls
      portal:
        secretName: harbor-portal-internal-tls
      chartmuseum:
        secretName: harbor-chartmuseum-internal-tls
      clair:
        secretName: harbor-clair-internal-tls
      trivy:
        secretName: harbor-trivy-internal-tls
    metrics:
      enabled: true
    EOF

    ix=2
    for item in `ls templates`; do
        if [ -d templates/$item ]; then
        filename="$(printf "manifests/%02d-%s.yaml" $ix $item)"
        for subitem in `ls templates/$item`; do
            content=`helm template $release . -s templates/$item/$subitem -f $valuesFile`
            content=`echo "$content" | sed '/^# Source: /d' | sed 's/#/#!/g'`
            if [[ $content != "{}"  ]]; then
            if [[ $content = *[!\ ]* ]]; then
                content=`echo "$content" | yq e '.metadata.namespace="harbor"' -`
                # echo '---' >> $filename
                echo "$content" >> $filename
                if [[ $item != "ingress" ]]; then
                sed -i '/checksum/d' $filename
                sed -i '/annotations/d' $filename
                fi
            fi
            fi
        done

        if [[ -z $(grep '[^[:space:]]' $filename) ]]; then
            rm -rf $filename
        else
            ((ix=ix+1))
        fi
        fi
    done

    sed -i '/tls.crt/d' manifests/04-exporter.yaml
    sed -i '/tls.key/d' manifests/04-exporter.yaml

    rm $valuesFile

    git checkout templates/_helpers.tpl
    ```

3. Run `generate-manifests.sh` to generate the manifests and the output files are in the `manifests` directory.

    ```sh
    chmod +x ./generate-manifests.sh
    ./generate-manifests.sh
    ```
