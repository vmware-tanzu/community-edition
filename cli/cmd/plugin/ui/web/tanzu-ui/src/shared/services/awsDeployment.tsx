// React imports
import { useContext } from 'react';
import { useNavigate } from 'react-router-dom';

// App imports
import { DEPLOYMENT_STATUS_CHANGED } from '../../state-management/actions/Deployment.actions';
import { TOGGLE_APP_STATUS } from '../../state-management/actions/Ui.actions';
import { Store } from '../../state-management/stores/Store';
import { AwsStore } from '../../state-management/stores/Store.aws';
import { AWSManagementClusterParams, AwsService, ConfigFileInfo, IdentityManagementConfig } from '../../swagger-api';
import { DeploymentStates, DeploymentTypes } from '../constants/Deployment.constants';
import { NavRoutes } from '../constants/NavRoutes.constants';
import { Providers } from '../constants/Providers.constants';
import { STORE_SECTION_FORM } from '../../state-management/reducers/Form.reducer';

const useAwsDeployment = () => {
    const { dispatch } = useContext(Store);
    const { awsState } = useContext(AwsStore);
    const navigate = useNavigate();
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    const getAwsRequestPayload = () => {
        const awsClusterParams: AWSManagementClusterParams = {
            awsAccountParams: {
                profileName: awsState[STORE_SECTION_FORM].PROFILE,
                sessionToken: awsState[STORE_SECTION_FORM].SESSION_TOKEN,
                region: awsState[STORE_SECTION_FORM].REGION,
                accessKeyID: awsState[STORE_SECTION_FORM].ACCESS_KEY_ID,
                secretAccessKey: awsState[STORE_SECTION_FORM].SECRET_ACCESS_KEY,
            },
            loadbalancerSchemeInternal: false,
            sshKeyName: awsState[STORE_SECTION_FORM].EC2_KEY_PAIR,
            createCloudFormationStack: false,
            clusterName: awsState[STORE_SECTION_FORM].CLUSTER_NAME,
            controlPlaneFlavor: awsState[STORE_SECTION_FORM].CLUSTER_PLAN,
            controlPlaneNodeType: awsState[STORE_SECTION_FORM].CONTROL,
            bastionHostEnabled: true,
            machineHealthCheckEnabled: true,
            vpc: {
                cidr: awsState[STORE_SECTION_FORM].VPC_CIDR,
                vpcID: '',
                azs: [
                    {
                        name: 'us-west-2a',
                        workerNodeType: awsState[STORE_SECTION_FORM].CLUSTER_WORKER_NODE_TYPE,
                        publicSubnetID: '',
                        privateSubnetID: '',
                    },
                ],
            },
            enableAuditLogging: false,
            networking: {
                networkName: '',
                clusterDNSName: '',
                clusterNodeCIDR: '',
                clusterServiceCIDR: awsState[STORE_SECTION_FORM].CLUSTER_SERVICE_CIDR,
                clusterPodCIDR: awsState[STORE_SECTION_FORM].CLUSTER_POD_CIDR,
                cniType: 'antrea',
            },
            ceipOptIn: true,
            labels: {},
            os: {
                name: 'ubuntu-20.04-amd64 (ami-0dd0327a3bfaa4dc8)',
                osInfo: {
                    arch: 'amd64',
                    name: 'ubuntu',
                    version: '20.04',
                },
            },
            annotations: {
                description: '',
                location: '',
            },
            identityManagement: {
                idm_type: IdentityManagementConfig.idm_type.NONE,
            },
        };
        return awsClusterParams;
    };

    const deployOnAws = async () => {
        const awsClusterParams: AWSManagementClusterParams = getAwsRequestPayload();
        try {
            const configFileInfo: ConfigFileInfo = await AwsService.applyTkgConfigForAws(awsClusterParams);
            await AwsService.createAwsManagementCluster(awsClusterParams);
            dispatch({
                type: DEPLOYMENT_STATUS_CHANGED,
                payload: {
                    type: DeploymentTypes.MANAGEMENT_CLUSTER,
                    status: DeploymentStates.RUNNING,
                    provider: Providers.AWS,
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
        getAwsRequestPayload,
        deployOnAws,
    };
};
export default useAwsDeployment;
