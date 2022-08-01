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
import { retrieveAwsInstanceType } from '../constants/defaults/aws.defaults';
import { STORE_SECTION_FORM } from '../../state-management/reducers/Form.reducer';

const useAwsDeployment = () => {
    const { dispatch } = useContext(Store);
    const { awsState } = useContext(AwsStore);
    const navigate = useNavigate();
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    // TODO: more dynamic population of this payload
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
            // TODO: may need to do discover around existing cloudformation stack
            createCloudFormationStack: false,
            clusterName: awsState[STORE_SECTION_FORM].CLUSTER_NAME,
            controlPlaneFlavor: awsState[STORE_SECTION_FORM].CLUSTER_PLAN,
            controlPlaneNodeType: retrieveAwsInstanceType(awsState[STORE_SECTION_FORM].NODE_PROFILE),
            bastionHostEnabled: awsState[STORE_SECTION_FORM].ENABLE_BASTION_HOST,
            machineHealthCheckEnabled: awsState[STORE_SECTION_FORM].ENABLE_MACHINE_HEALTH_CHECK,
            vpc: {
                cidr: awsState[STORE_SECTION_FORM].VPC_CIDR,
                vpcID: '',
                // TODO: single subregion name populated from region selection; but does not support multi-az/HA
                azs: [
                    {
                        name: awsState[STORE_SECTION_FORM].REGION + 'a',
                        workerNodeType: retrieveAwsInstanceType(awsState[STORE_SECTION_FORM].NODE_PROFILE),
                        publicSubnetID: '',
                        privateSubnetID: '',
                    },
                ],
            },
            enableAuditLogging: awsState[STORE_SECTION_FORM].ENABLE_AUDIT_LOGGING,
            networking: {
                networkName: '',
                clusterDNSName: '',
                clusterNodeCIDR: '',
                clusterServiceCIDR: awsState[STORE_SECTION_FORM].CLUSTER_SERVICE_CIDR,
                clusterPodCIDR: awsState[STORE_SECTION_FORM].CLUSTER_POD_CIDR,
                cniType: awsState[STORE_SECTION_FORM].CLUSTER_NETWORKING_CNI_PROVIDER,
            },
            ceipOptIn: awsState[STORE_SECTION_FORM].ENABLE_CEIP_PARTICIPATION,
            labels: {},
            os: awsState[STORE_SECTION_FORM].OS_IMAGE,
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
