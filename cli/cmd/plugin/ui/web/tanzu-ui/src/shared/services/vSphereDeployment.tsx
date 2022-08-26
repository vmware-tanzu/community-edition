import { ConfigFileInfo, VsphereManagementClusterParams, VsphereService } from '../../swagger-api';
import { DeploymentStates, DeploymentTypes } from '../constants/Deployment.constants';

// App imports
import { DEPLOYMENT_STATUS_CHANGED } from '../../state-management/actions/Deployment.actions';
import { NavRoutes } from '../constants/NavRoutes.constants';
import { Providers } from '../constants/Providers.constants';
import { STORE_SECTION_FORM } from '../../state-management/reducers/Form.reducer';
import { Store } from '../../state-management/stores/Store';
import { TOGGLE_APP_STATUS } from '../../state-management/actions/Ui.actions';
import { VsphereStore } from '../../views/management-cluster/vsphere/Store.vsphere.mc';
// React imports
import { useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import { VSPHERE_FIELDS } from '../../views/management-cluster/vsphere/VsphereManagementCluster.constants';

const useVSphereDeployment = () => {
    const { dispatch } = useContext(Store);
    const { vsphereState } = useContext(VsphereStore);

    const navigate = useNavigate();

    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    const getVSphereRequestPayload = () => {
        const vSphereClusterParams: VsphereManagementClusterParams = {
            vsphereCredentials: {
                host: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.SERVERNAME],
                insecure: !vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.USETHUMBPRINT],
                password: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.PASSWORD],
                thumbprint: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.THUMBPRINT],
                username: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.USERNAME],
            },
            clusterName: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.CLUSTERNAME],
            datacenter: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.DATACENTER],
            resourcePool: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.Pool].name,
            datastore: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.DataStore].name,
            folder: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.VMFolder].name,
            controlPlaneNodeType: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.CONTROL_PLANE_INSTANCE_TYPE],
            controlPlaneFlavor: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.CONTROL_PLANE_FLAVOR],
            workerNodeType: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.WORKER_INSTANCE_TYPE],
            numOfWorkerNode: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.NUM_WORKER_NODES] || 0,
            kubernetesVersion: '',
            ipFamily: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.IPFAMILY],
            networking: {
                networkName: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.Network].name,
                clusterDNSName: '',
                clusterNodeCIDR: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.CLUSTER_NODE_CIDR],
                clusterServiceCIDR: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.CLUSTER_SERVICE_CIDR],
                clusterPodCIDR: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.CLUSTER_POD_CIDR],
                cniType: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.CNI_TYPE],
                /*
                httpProxyConfiguration: {
                    enabled: false,
                    HTTPProxyURL: '',
                    HTTPProxyUsername: '',
                    HTTPProxyPassword: '',
                    HTTPSProxyURL: '',
                    HTTPSProxyUsername: '',
                    HTTPSProxyPassword: '',
                    noProxy: '',
                },
*/
            },
            os: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.VMTEMPLATE],
            ssh_key: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.SSHKEY],
            machineHealthCheckEnabled: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.MACHINE_HEALTH_CHECK_ACTIVATED] || false,
            ceipOptIn: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.CEIP_OPT_IN] || false,
            enableAuditLogging: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.ENABLE_AUDIT_LOGGING] || false,
            /*
            annotations: { a: '' },
            labels: { '': '' },
*/
            controlPlaneEndpoint: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.CLUSTER_ENDPOINT],
            identityManagement: undefined,
            aviConfig: undefined,
        };
        return vSphereClusterParams;
    };

    const deployOnVSphere = async () => {
        const vSphereClusterParams: VsphereManagementClusterParams = getVSphereRequestPayload();
        try {
            const configFileInfo: ConfigFileInfo = await VsphereService.applyTkgConfigForVsphere(vSphereClusterParams);
            await VsphereService.createVSphereManagementCluster(vSphereClusterParams);
            dispatch({
                type: DEPLOYMENT_STATUS_CHANGED,
                payload: {
                    type: DeploymentTypes.MANAGEMENT_CLUSTER,
                    status: DeploymentStates.RUNNING,
                    provider: Providers.VSPHERE,
                    configPath: configFileInfo.path,
                },
            });
        } catch (e) {
            console.error(`Error when calling config or create API: ${e}`);
            return;
        }

        dispatch({
            type: TOGGLE_APP_STATUS,
        });
        navigateToProgress();
    };

    return {
        getVSphereRequestPayload,
        deployOnVSphere,
    };
};
export default useVSphereDeployment;
