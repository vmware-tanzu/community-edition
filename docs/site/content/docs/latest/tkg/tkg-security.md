# Tanzu Kubernetes Grid Security and Networking

This topic lists resources for securing Tanzu Kubernetes Grid infrastructure.

## <a id="ports"></a> Ports and Protocols

Networking ports and protocols used by Tanzu Kubernetes Grid are listed in the [VMWare Ports and Protocols](https://ports.vmware.com/home/Tanzu-Kubernetes-Grid) tool.

For each internal communication path, [VMWare Ports and Protocols](https://ports.vmware.com/home/Tanzu-Kubernetes-Grid) lists:

- Product
- Version
- Source
- Destination
- Ports
- Protocols
- Purpose
- Service Description
- Classification (Outgoing, Incoming, or Bidirectional)

You can use this information to configure firewall rules.

## <a id="rules"></a> Tanzu Kubernetes Grid Firewall Rules

| **Name**         | **Source**       | **Destination**  | **Service**      | **Purpose**      |
| ---------------- | ---------------- | ---------------- | ---------------- | ---------------- |
| workload-to-mgmt | tkg-workload-cluster-network    | tkg-management-cluster-network   | TCP:6443   | Allow workload cluster to register with management cluster   |
| mgmt-to-workload | tkg-management-cluster-network  | tkg-workload-cluster-network   | TCP:6443, 5556 | Allow management network to configure workload cluster    |
| allow-mgmt-subet | tkg-management-cluster-network  | tkg-management-cluster-network   | all    | Allow all internal cluster communication    |
| allow-wl-subnet  | tkg-workload-cluster-network    | tkg-workload-cluster-network   | all    | Allow all internal cluster communication    |
| jumpbox-to-k8s   | Jumpbox IP   | tkg-management-cluster-network, tkg-workload-cluster-network | TCP:6443 | Allow Jumpbox to create management cluster and manage clusters.    |
| dhcp             | any          | NSX-T: any / no NSX-T: DHCP IP    | DHCP   | Allows hosts to get DHCP addresses. |
| to-harbor        | tkg-management-cluster-network, tkg-workload-cluster-network, jumpbox IP | Harbor IP    | HTTPS    | Allows components to retrieve container images    |
| to-vcenter       | tkg-management-cluster-network, tkg-workload-cluster-network, jumpbox IP | vCenter IP   | HTTPS    | Allows components to access vSphere to create VMs and Storage Volumes   |
| dns-ntp-outbound | tkg-management-cluster-network, tkg-workload-cluster-network, jumpbox IP | DNS, NTP servers   | DNS, NTP   | Core services   |
| ssh-to-jumpbox   | any          | Jumpbox IP   | SSH    | Access from outside to the jumpbox    |
| deny-all         | any          | any    | all    | deny by default  |

## <a id="benchmarking"></a> CIS Benchmarking for Clusters

To assess cluster security, you can run Center for Internet Security (CIS) [Kubernetes benchmark](https://www.cisecurity.org/benchmark/kubernetes/) tests on clusters deployed by Tanzu Kubernetes Grid.

For clusters that do not pass all sections of the tests,
see explanations and possible remediations in the **Expected Test Failures for CIS Benchmark Inspection on Provisioned Tanzu Kubernetes Clusters** table under [About the CIS Benchmark Inspection and Provisioned Tanzu Kubernetes Clusters](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-concepts/GUID-1BB98612-A9BF-494C-8446-1DB2E80BF5F9.html#about-the-cis-benchmark-inspection-and-provisioned-tanzu-kubernetes-clusters-0), in the Tanzu Mission Control documentation.

For any cluster listed in Tanzu Mission Control, you can run CIS benchmarking by opening the cluster's **Overview** tab, clicking **Run Inspection**, and selecting **Inspection type** > **CIS benchmark**.

To register a cluster with Tanzu Mission Control, which lists it in the interface, see:

- **Management clusters**: [Register Your Management Cluster with Tanzu Mission Control](mgmt-clusters/register_tmc.md)
- **Tanzu Kubernetes (workload) clusters**: [Add a Workload Cluster into Tanzu Mission Control Management](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-using/GUID-78908829-CB4E-459F-AA81-BEA415EC9A11.html) in the Tanzu Mission Control documentation
