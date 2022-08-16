// React imports
import { useContext } from 'react';
import { useNavigate } from 'react-router-dom';

// App imports
import { AzureManagementClusterParams, AzureService, ConfigFileInfo } from '../../swagger-api';
import { AzureStore } from '../../views/management-cluster/azure/store/Azure.store.mc';
import { DeploymentStates, DeploymentTypes } from '../constants/Deployment.constants';
import { DEPLOYMENT_STATUS_CHANGED } from '../../state-management/actions/Deployment.actions';
import { NavRoutes } from '../constants/NavRoutes.constants';
import { Providers } from '../constants/Providers.constants';
import { retrieveAzureInstanceType } from '../constants/defaults/azure.defaults';
import { Store } from '../../state-management/stores/Store';
import { STORE_SECTION_FORM } from '../../state-management/reducers/Form.reducer';
import { TOGGLE_APP_STATUS } from '../../state-management/actions/Ui.actions';
import { AZURE_FIELDS } from '../../views/management-cluster/azure/azure-mc-basic/AzureManagementClusterBasic.constants';

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
                subscriptionId: formInfo[AZURE_FIELDS.SUBSCRIPTION_ID],
                tenantId: formInfo[AZURE_FIELDS.TENANT_ID],
                clientId: formInfo[AZURE_FIELDS.CLIENT_ID],
                clientSecret: formInfo[AZURE_FIELDS.CLIENT_SECRET],
                azureCloud: formInfo[AZURE_FIELDS.AZURE_ENVIRONMENT],
            },
            ceipOptIn: formInfo[AZURE_FIELDS.CEIP_OPT_IN],
            clusterName: formInfo[AZURE_FIELDS.CLUSTER_NAME],
            controlPlaneFlavor: formInfo[AZURE_FIELDS.CONTROL_PLANE_FLAVOR],
            controlPlaneMachineType: retrieveAzureInstanceType(formInfo[AZURE_FIELDS.NODE_PROFILE]),
            controlPlaneSubnet: formInfo[AZURE_FIELDS.CONTROL_PLANE_SUBNET],
            controlPlaneSubnetCidr: formInfo[AZURE_FIELDS.CONTROL_PLANE_SUBNET_CIDR],
            enableAuditLogging: formInfo[AZURE_FIELDS.ACTIVATE_AUDIT_LOGGING],
            isPrivateCluster: formInfo[AZURE_FIELDS.PRIVATE_AZURE_CLUSTER],
            location: formInfo[AZURE_FIELDS.REGION],
            machineHealthCheckEnabled: formInfo[AZURE_FIELDS.MACHINE_HEALTH_CHECK_ENABLED],
            networking: {
                clusterPodCIDR: formInfo[AZURE_FIELDS.CLUSTER_POD_CIDR],
                clusterServiceCIDR: formInfo[AZURE_FIELDS.CLUSTER_SERVICE_CIDR],
                cniType: formInfo[AZURE_FIELDS.CNI_TYPE],
            },
            os: formInfo[AZURE_FIELDS.OS_IMAGE],
            resourceGroup: formInfo[AZURE_FIELDS.RESOURCE_GROUP],
            sshPublicKey: formInfo[AZURE_FIELDS.SSH_PUBLIC_KEY],
            vnetCidr: formInfo[AZURE_FIELDS.VNET_CIDR],
            vnetName: formInfo[AZURE_FIELDS.VNET_NAME],
            vnetResourceGroup: formInfo[AZURE_FIELDS.RESOURCE_GROUP],
            workerMachineType: retrieveAzureInstanceType(formInfo[AZURE_FIELDS.NODE_PROFILE]),
            workerNodeSubnet: formInfo[AZURE_FIELDS.WORKER_NODE_SUBNET],
            workerNodeSubnetCidr: formInfo[AZURE_FIELDS.WORKER_NODE_SUBNET_CIDR],
        };
        return azureClusterParams;
    };

    const deployOnAzure = async () => {
        const azureClusterParams: AzureManagementClusterParams = getAzureRequestPayload();
        try {
            const configFileInfo: ConfigFileInfo = await AzureService.applyTkgConfigForAzure(azureClusterParams);
            await AzureService.createAzureManagementCluster(azureClusterParams);
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
