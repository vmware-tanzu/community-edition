// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package contour_test

// invalidInfrastructureProvider has an invalid infrastructureProvider value.
const invalidInfrastructureProvider string = `
#@data/values
---
infrastructureProvider: invalid-value
`

// noNamespace has an empty namespace value.
const noNamespace string = `
#@data/values
---
namespace: ""
`

// noContourReplicas has an empty contour.replicas value.
const noContourReplicas string = `
#@data/values
---
contour:
  replicas: 0
`

// noCertificatesDuration has certificates.useCertManager set to true and
// an empty certificates.duration value.
const noCertificatesDuration string = `
#@data/values
---
certificates:
  useCertManager: true
  duration: 0
`

// noCertificatesRenewBefore has certificates.useCertManager set to true and
// an empty certificates.renewBefore value.
const noCertificatesRenewBefore string = `
#@data/values
---
certificates:
  useCertManager: true
  renewBefore: 0
`

// noHostPortHttp has envoy.hostPorts.enable set to true and
// an empty envoy.hostPorts.http value.
const noHostPortHttp string = `
#@data/values
---
envoy:
  hostPorts:
    enable: true
    http: 0
`

// noHostPortHttps has envoy.hostPorts.enable set to true and
// an empty envoy.hostPorts.https value.
const noHostPortHttps string = `
#@data/values
---
envoy:
  hostPorts:
    enable: true
    https: 0
`

// invalidEnvoyLogLevel has an invalid envoy.logLevel value.
const invalidEnvoyLogLevel string = `
#@data/values
---
envoy:
  logLevel: invalid-log-level
`

// noTerminationGracePeriodSeconds has an empty envoy.terminationGracePeriodSeconds value.
const noTerminationGracePeriodSeconds string = `
#@data/values
---
envoy:
  terminationGracePeriodSeconds: 0
`

// invalidEnvoyServiceType has an invalid envoy.service.type value.
const invalidEnvoyServiceType string = `
#@data/values
---
envoy:
  service:
    type: InvalidServiceType
`

// invalidEnvoyServiceExternalTrafficPolicy has an invalid envoy.service.externalTrafficPolicy value.
const invalidEnvoyServiceExternalTrafficPolicy string = `
#@data/values
---
envoy:
  service:
    externalTrafficPolicy: InvalidExternalTrafficPolicy
`

// invalidAwsLoadBalancerType has infrastructureProvider set to aws and
// an invalid envoy.service.aws.loadBalancerType value.
const invalidAwsLoadBalancerType string = `
#@data/values
---
infrastructureProvider: aws
envoy:
  service:
    aws:
      loadBalancerType: invalid-type
`

// nonDefaultNamespace has namespace set to non-default-namespace.
const nonDefaultNamespace string = `
#@data/values
---
namespace: non-default-namespace
`

// contourConfigFileContents has contour.configFileContents set.
const contourConfigFileContents string = `
#@data/values
---
contour:
  #@overlay/replace
  configFileContents:
    foo:
      bar: baz
    boo: bam
`

// dockerInfrastructureProvider has infrastructureProvider set to docker.
const dockerInfrastructureProvider string = `
#@data/values
---
infrastructureProvider: docker
`

// awsInfrastructureProvider has infrastructureProvider set to aws.
const awsInfrastructureProvider string = `
#@data/values
---
infrastructureProvider: aws
`

// vsphereInfrastructureProvider has infrastructureProvider set to vsphere.
const vsphereInfrastructureProvider string = `
#@data/values
---
infrastructureProvider: vsphere
`

// azureInfrastructureProvider has infrastructureProvider set to azure.
const azureInfrastructureProvider string = `
#@data/values
---
infrastructureProvider: azure
`

// vsphereClusterIPEnvoyService has infrastructureProvider set to vsphere
// and envoy.service.type set to ClusterIP.
const vsphereClusterIPEnvoyService string = `
#@data/values
---
infrastructureProvider: vsphere
envoy:
  service:
    type: ClusterIP
`

// vsphereLocalExternalTrafficPolicy has infrastructureProvider set to vsphere
// and envoy.service.externalTrafficPolicy set to Local.
const vsphereLocalExternalTrafficPolicy string = `
#@data/values
---
infrastructureProvider: vsphere
envoy:
  service:
    externalTrafficPolicy: Local
`

// contourReplicasThree has contour.replicas set to 3.
const contourReplicasThree string = `
#@data/values
---
contour:
  replicas: 3
`

// awsLoadBalancerTypeNLB has infrastructureProvider set to aws and
// envoy.service.aws.loadBalancerType set to nlb.
const awsLoadBalancerTypeNLB = `
#@data/values
---
infrastructureProvider: aws
envoy:
  service:
    aws:
      loadBalancerType: nlb
`

// awsLoadBalancerTypeClassic has infrastructureProvider set to aws and
// envoy.service.aws.loadBalancerType set to classic.
const awsLoadBalancerTypeClassic = `
#@data/values
---
infrastructureProvider: aws
envoy:
  service:
    aws:
      loadBalancerType: classic
`

// awsLoadBalancerTypeClassicProxyProtocolEnabled has infrastructureProvider set to aws and
// envoy.service.aws.loadBalancerType set to classic and contour.useProxyProtocol
// set to true.
const awsLoadBalancerTypeClassicProxyProtocolEnabled = `
#@data/values
---
infrastructureProvider: aws
envoy:
  service:
    aws:
      loadBalancerType: classic
contour:
  useProxyProtocol: true
`

// contourUseProxyProtocolEnabled has contour.useProxyProtocol set to true.
const contourUseProxyProtocolEnabled = `
#@data/values
---
contour:
  useProxyProtocol: true
`

// contourLogLevelDebug has contour.logLevel set to debug.
const contourLogLevelDebug = `
#@data/values
---
contour:
  logLevel: debug
`

// envoyServiceLoadBalancerIP has envoy.service.loadBalancerIP set to 7.7.7.7.
const envoyServiceLoadBalancerIP = `
#@data/values
---
envoy:
  service:
    loadBalancerIP: 7.7.7.7
`

// envoyServiceAnnotations has envoy.service.annotations set.
const envoyServiceAnnotations = `
#@data/values
---
envoy:
  service:
    #@overlay/replace
    annotations:
      foo: bar
      boo: baz
`

// envoyServiceAnnotations has envoy.service.annotations set and
// infrastructureProvider set to aws and envoy.service.aws set.
const envoyServiceAnnotationsAWSNLB = `
#@data/values
---
infrastructureProvider: aws
envoy:
  service:
    #@overlay/replace
    annotations:
      foo: bar
      boo: baz
    aws:
      loadBalancerType: nlb
`

// envoyServiceNodePorts has envoy.service.type set to NodePort and envoy.service.nodePorts set.
const envoyServiceNodePorts = `
#@data/values
---
envoy:
  service:
    type: NodePort
    nodePorts:
      http: 30080
      https: 30443
`

// envoyHostPorts has envoy.hostPorts set.
const envoyHostPorts = `
#@data/values
---
envoy:
  hostPorts:
    enable: true
    http: 80
    https: 443
`

// envoyHostNetworkingEnabled has envoy.hostNetwork set to true.
const envoyHostNetworkingEnabled = `
#@data/values
---
envoy:
  hostNetwork: true
`

// envoyTerminationGracePeriodSeconds has envoy.terminationGracePeriodSeconds set to 777.
const envoyTerminationGracePeriodSeconds = `
#@data/values
---
envoy:
  terminationGracePeriodSeconds: 777
`

// envoyLogLevelDebug has envoy.logLevel set to debug.
const envoyLogLevelDebug = `
#@data/values
---
envoy:
  logLevel: debug
`

// useCertManagerEnabled has certificates.useCertManager set to true.
const useCertManagerEnabled = `
#@data/values
---
certificates:
  useCertManager: true
`

// useCertManagerEnabledDurationRenewBeforeSet has certificates.useCertManager set to true
// and certificates.duration and certificates.renewBefore set.
const useCertManagerEnabledDurationRenewBeforeSet = `
#@data/values
---
certificates:
  useCertManager: true
  duration: 777h
  renewBefore: 77h
`
