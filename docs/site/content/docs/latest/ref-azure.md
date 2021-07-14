# Reference for Azure account

If you encounter issues deploying a cluster to vSphere, review the following troubleshooting and reference content:

## <a id="nsgs"></a> Network Security Groups on Azure {nsg}

If you do not specify a VNET when deploying a management cluster, the deployment process creates a new VNET along with the NSGs required for the management cluster.

The following NSGs are required:

- One control plane NSG shared by the control plane nodes of all clusters, including the management cluster and the workload clusters that it manages.
- One worker NSG for each cluster, for the cluster's worker nodes.

Management and workload clusters on Azure require the following Network Security Groups (NSGs) to be defined on their VNET.
   - A subnet for the management cluster control plane node
   - A Network Security Group on the control plane subnet with the following inbound security rules, to enable SSH and Kubernetes API server connections:
      - Allow TCP over port 22 for any source and destination
      - Allow TCP over port 6443 for any source and destination.
      Port 6443 is where the Kubernetes API is exposed on VMs in the clusters you create.
   - A subnet and Network Security Group for the management cluster worker nodes.

For each workload cluster that you deploy later, you need to create a worker NSG named `<CLUSTER-NAME>-node-nsg`, where `<CLUSTER-NAME>` is the name of the workload cluster.
This worker NSG must have the same VNET and region as its management cluster.

## Microsoft Azure account
Your Microsoft Azure account should have the following permissions and requirements:
   - Permissions required to register an app. See [Permissions required for registering an app](https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal#permissions-required-for-registering-an-app) in the Azure documentation.
   - Sufficient VM core (vCPU) quotas for your clusters. A standard Azure account has a quota of 10 vCPU per region. Tanzu Community Edition clusters require 2 vCPU per node, which translates to:
     - Management cluster:
         - `dev` plan: 4 vCPU (1 main, 1 worker)
         - `prod` plan: 8 vCPU (3 main , 1 worker)
     - Each workload cluster:
         - `dev` plan: 4 vCPU (1 main, 1 worker)
         - `prod` plan: 12 vCPU (3 main , 3 worker)
     - For example, assuming a **single management cluster** and all clusters with the same plan:
   <table width="100%" border="0">
   <tr>
     <th width="17%">Plan</th>
     <th width="22%">Workload Clusters</th>
     <th width="22%">vCPU for Workload</th>
     <th width="22%">vCPU for Management</th>
     <th width="17%">Total vCPU</th>
   </tr>
   <tr>
     <td rowspan="2">Dev</td>
     <td>1</td>
     <td>4</td>
     <td rowspan="2">4</td>
     <td>8</td>
   </tr>
   <tr>
     <td>5</td>
     <td>20</td>
     <td>24</td>
   </tr>
   <tr>
     <td rowspan="2">Prod</td>
     <td>1</td>
     <td>12</td>
     <td rowspan="2">8</td>
     <td>20</td>
   </tr>
   <tr>
     <td>5</td>
     <td>60</td>
     <td>68</td>
   </tr>
   </table>
   - Sufficient public IP address quotas for your clusters, including the quota for Public IP Addresses - Standard, Public IP Addresses - Basic, and Static Public IP Addresses. A standard Azure account has a quota of 10 public IP addresses per region. Every Tanzu Community Edition cluster requires 2 Public IP addresses regardless of how many control plane nodes and worker nodes it has. For each Kubernetes Service object with type `LoadBalancer`, 1 Public IP address is required.
   - Run a DNS lookup on all `imageRepository` values to find their CNAMEs.


