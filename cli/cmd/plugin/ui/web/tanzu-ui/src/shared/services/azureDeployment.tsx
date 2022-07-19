// React imports
import { useContext } from 'react';
import { useNavigate } from 'react-router-dom';

// App imports
import { DEPLOYMENT_STATUS_CHANGED } from '../../state-management/actions/Deployment.actions';
import { TOGGLE_APP_STATUS } from '../../state-management/actions/Ui.actions';
import { Store } from '../../state-management/stores/Store';
import { AzureStore } from '../../state-management/stores/Azure.store';
import { AWSManagementClusterParams, AzureManagementClusterParams, AzureService, ConfigFileInfo } from '../../swagger-api';
import { DeploymentStates, DeploymentTypes } from '../constants/Deployment.constants';
import { NavRoutes } from '../constants/NavRoutes.constants';
import { Providers } from '../constants/Providers.constants';
import { STORE_SECTION_FORM } from '../../state-management/reducers/Form.reducer';
import { retrieveAzureInstanceType } from '../constants/defaults/azure.defaults';

const useAzureDeployment = () => {
    const { dispatch } = useContext(Store);
    const { azureState } = useContext(AzureStore);
    const navigate = useNavigate();
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    const getAzureRequestPayload = () => {
        const formInfo = azureState[STORE_SECTION_FORM];
        const azureClusterParams: AzureManagementClusterParams = {
            azureAccountParams: {
                subscriptionId: formInfo.SUBSCRIPTION_ID,
                tenantId: formInfo.TENANT_ID,
                clientId: formInfo.CLIENT_ID,
                clientSecret: formInfo.CLIENT_SECRET,
                azureCloud: formInfo.AZURE_ENVIRONMENT,
            },
            ceipOptIn: formInfo.CEIP_OPT_IN,
            clusterName: formInfo.CLUSTER_NAME,
            controlPlaneFlavor: formInfo.CONTROL_PLANE_FLAVOR,
            controlPlaneMachineType: retrieveAzureInstanceType(formInfo.NODE_PROFILE),
            controlPlaneSubnet: formInfo.CONTROL_PLANE_SUBNET,
            controlPlaneSubnetCidr: formInfo.CONTROL_PLANE_SUBNET_CIDR,
            enableAuditLogging: formInfo.ACTIVATE_AUDIT_LOGGING,
            isPrivateCluster: formInfo.PRIVATE_AZURE_CLUSTER,
            location: formInfo.REGION,
            machineHealthCheckEnabled: formInfo.MACHINE_HEALTH_CHECK_ENABLED,
            networking: {
                clusterPodCIDR: formInfo.CLUSTER_POD_CIDR,
                clusterServiceCIDR: formInfo.CLUSTER_SERVICE_CIDR,
                cniType: formInfo.CNI_TYPE,
            },
            os: {
                name: formInfo.IMAGE_NAME,
                osInfo: {
                    arch: formInfo.IMAGE_OS_ARCH,
                    name: formInfo.IMAGE_OS_NAME,
                    version: formInfo.IMAGE_OS_VERSION,
                },
            },
            resourceGroup: formInfo.RESOURCE_GROUP,
            sshPublicKey: formInfo.SSH_PUBLIC_KEY,
            vnetCidr: formInfo.VNET_CIDR,
            vnetName: formInfo.VNET_NAME,
            vnetResourceGroup: formInfo.RESOURCE_GROUP,
            workerMachineType: retrieveAzureInstanceType(formInfo.NODE_PROFILE),
            workerNodeSubnet: formInfo.WORKER_NODE_SUBNET,
            workerNodeSubnetCidr: formInfo.WORKER_NODE_SUBNET_CIDR,
        };
        return azureClusterParams;
    };

    const deployOnAzure = async () => {
        const awsClusterParams: AWSManagementClusterParams = getAzureRequestPayload();
        try {
            const configFileInfo: ConfigFileInfo = await AzureService.applyTkgConfigForAzure(awsClusterParams);
            await AzureService.createAzureManagementCluster(awsClusterParams);
            dispatch({
                type: DEPLOYMENT_STATUS_CHANGED,
                payload: {
                    type: DeploymentTypes.MANAGEMENT_CLUSTER,
                    status: DeploymentStates.RUNNING,
                    provider: Providers.AZURE,
                    configPath: configFileInfo.path,
                },
            });
        } catch (e) {
            console.log(`Error when calling config or create API: ${e}`);
            return;
        }

        dispatch({
            type: TOGGLE_APP_STATUS,
        });
        navigateToProgress();
    };

    return {
        getAzureRequestPayload,
        deployOnAzure,
    };
};
export default useAzureDeployment;
