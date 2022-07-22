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
                host: vsphereState[STORE_SECTION_FORM].SERVERNAME,
                insecure: false,
                password: vsphereState[STORE_SECTION_FORM].PASSWORD,
                thumbprint: '',
                username: vsphereState[STORE_SECTION_FORM].USERNAME,
            },
            clusterName: vsphereState[STORE_SECTION_FORM].CLUSTER_NAME,
            datacenter: vsphereState[STORE_SECTION_FORM].DATACENTER,
            resourcePool: '',
            datastore: vsphereState[STORE_SECTION_FORM].Datastore,
            folder: vsphereState[STORE_SECTION_FORM].VMFolder,
            controlPlaneNodeType: '',
            controlPlaneFlavor: '',
            workerNodeType: '',
            numOfWorkerNode: 0,
            kubernetesVersion: '',
            ipFamily: vsphereState[STORE_SECTION_FORM].IPFAMILY,
            networking: {
                networkName: '',
                clusterDNSName: '',
                clusterNodeCIDR: '',
                clusterServiceCIDR: '',
                clusterPodCIDR: '',
                cniType: '',
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
            },
            os: {
                name: '',
                moid: '',
                k8sVersion: '',
                isTemplate: false,
                osInfo: {
                    name: '',
                    version: '',
                    arch: '',
                },
            },
            ssh_key: vsphereState[STORE_SECTION_FORM].SSHKEY,
            machineHealthCheckEnabled: false,
            ceipOptIn: false,
            enableAuditLogging: false,
            annotations: { a: '' },
            labels: { '': '' },
            controlPlaneEndpoint: '',
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
