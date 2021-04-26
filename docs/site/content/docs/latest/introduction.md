DRAFT DRAFT DRAFT DRAFT DRAFT DRAFT
WIP WIP WIP
NOT READY FOR REVEIEW

Tanzu Community can be deployed across several platforms including VMware vSphere, AWS, Azure, and on a desktop Docker instance.

Management Cluster
Description 1: After you deploy the Tanzu ClI, this is the first element you deploy. The management cluster  provides management and operations for your instance. It runs Cluster-API which is used to create your workload clusters, as well as creating shared services for all your clusters within the instance. A management cluster can be deployed using the Tanzu Kubernets Grid installer ui.

Description 2: This is the first architectural components to be deployed for creating a TKG instance. The management cluster is a dedicated cluster for management and operation of your whole TKG instance infrastructure. A management cluster will have Antrea networking enabled by default. This runs cluster API to create the additional clusters for your workloads to run, as well as the shared and in-cluster services for all clusters within the instance to use.
Note: It is not recommended that the management cluster be used as a general-purpose compute environment for your application workloads.

Workload cluster

Description 1: After you have deployed the managment cluster, you deploy the workload clusters. The workload clusters are deployed by the management cluster. The workload clusters are deployed using the Tanzu CLI only, there is no UI available for this.  The workload clusters will be used to run your applications. You can deploy multiple clusters at the Kubernetes versions you need. The management and workload clusters can be different versions. 

Description 2: Once you have deployed your management cluster, you can deploy additional CNCF conformant Kubernetes clusters and manage their full lifecycle. These clusters are designed to run your application workloads, managed via your management cluster. These clusters can run different Kubernetes versions as required. These clusters use Antrea networking by default.

These clusters are referred to as Workload Clusters when working with the Tanzu CLI.