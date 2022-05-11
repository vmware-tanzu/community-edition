// React imports
import { useContext } from 'react';
import { useNavigate } from 'react-router-dom';

// App imports
import { DEPLOYMENT_STATUS_CHANGED } from '../../state-management/actions/Deployment.actions';
import { TOGGLE_APP_STATUS } from '../../state-management/actions/Ui.actions';
import { DockerStore } from '../../state-management/stores/Docker.store';
import { Store } from '../../state-management/stores/Store';
import { DockerService, ConfigFileInfo, DockerManagementClusterParams, IdentityManagementConfig } from '../../swagger-api';
import { DeploymentStates, DeploymentTypes } from '../constants/Deployment.constants';
import { NavRoutes } from '../constants/NavRoutes.constants';
import { Providers } from '../constants/Providers.constants';

const useDockerDeployment = () => {
    const { dispatch } = useContext(Store);
    const { dockerState } = useContext(DockerStore);
    const navigate = useNavigate();
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    const getDockerRequestPayload = () => {
        const dockerClusterParams: DockerManagementClusterParams = {
            clusterName: dockerState.data.CLUSTER_NAME,
            networking: {
                clusterDNSName: '',
                clusterNodeCIDR: '',
                clusterServiceCIDR: dockerState.data.CLUSTER_SERVICE_CIDR,
                clusterPodCIDR: dockerState.data.CLUSTER_POD_CIDR,
                cniType: dockerState.data.CNI_TYPE,
            },
            identityManagement: {
                idm_type: IdentityManagementConfig.idm_type.NONE,
            },
        };
        return dockerClusterParams;
    };

    const deployOnDocker = async () => {
        const dockerClusterParams: DockerManagementClusterParams = getDockerRequestPayload();
        try {
            const configFileInfo: ConfigFileInfo = await DockerService.applyTkgConfigForDocker(dockerClusterParams);
            await DockerService.createDockerManagementCluster(dockerClusterParams);
            dispatch({
                type: DEPLOYMENT_STATUS_CHANGED,
                payload: {
                    type: DeploymentTypes.MANAGEMENT_CLUSTER,
                    status: DeploymentStates.RUNNING,
                    provider: Providers.DOCKER,
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
        getDockerRequestPayload,
        deployOnDocker,
    };
};
export default useDockerDeployment;
