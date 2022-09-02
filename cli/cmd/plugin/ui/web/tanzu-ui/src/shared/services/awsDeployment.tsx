// React imports
import { useContext } from 'react';
import { useNavigate } from 'react-router-dom';

// App imports
import { AWS_FIELDS, AWS_NODE_PROFILE_NAMES } from '../../views/management-cluster/aws/aws-mc-basic/AwsManagementClusterBasic.constants';
import { AwsStore } from '../../views/management-cluster/aws/store/Aws.store.mc';
import { AWSManagementClusterParams, AwsService, ConfigFileInfo, IdentityManagementConfig } from '../../swagger-api';
import { DEPLOYMENT_STATUS_CHANGED } from '../../state-management/actions/Deployment.actions';
import { DeploymentStates, DeploymentTypes } from '../constants/Deployment.constants';
import { NavRoutes } from '../constants/NavRoutes.constants';
import { Providers } from '../constants/Providers.constants';
import { Store } from '../../state-management/stores/Store';
import { STORE_SECTION_FORM } from '../../state-management/reducers/Form.reducer';
import { STORE_SECTION_RESOURCES } from '../../state-management/reducers/Resources.reducer';
import { TOGGLE_APP_STATUS } from '../../state-management/actions/Ui.actions';

const useAwsDeployment = () => {
    const { dispatch } = useContext(Store);
    const { awsState } = useContext(AwsStore);
    const navigate = useNavigate();
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    // TODO: more dynamic population of this payload
    const getAwsRequestPayload = () => {
        const awsData = awsState[STORE_SECTION_FORM];
        const nodeType = awsState[STORE_SECTION_RESOURCES][AWS_FIELDS.NODE_TYPE];
        const awsClusterParams: AWSManagementClusterParams = {
            awsAccountParams: {
                profileName: awsData[AWS_FIELDS.PROFILE],
                sessionToken: awsData[AWS_FIELDS.SESSION_TOKEN],
                region: awsData[AWS_FIELDS.REGION],
                accessKeyID: awsData[AWS_FIELDS.ACCESS_KEY_ID],
                secretAccessKey: awsData[AWS_FIELDS.SECRET_ACCESS_KEY],
            },
            loadbalancerSchemeInternal: false,
            sshKeyName: awsData[AWS_FIELDS.EC2_KEY_PAIR],
            // TODO: may need to do discover around existing cloudformation stack
            createCloudFormationStack: false,
            clusterName: awsData[AWS_FIELDS.CLUSTER_NAME],
            controlPlaneFlavor: awsData[AWS_FIELDS.CLUSTER_PLAN],
            controlPlaneNodeType: nodeType[awsData[AWS_FIELDS.NODE_PROFILE]],
            bastionHostEnabled: awsData[AWS_FIELDS.ENABLE_BASTION_HOST],
            machineHealthCheckEnabled: awsData[AWS_FIELDS.ENABLE_MACHINE_HEALTH_CHECK],
            vpc: {
                cidr: awsData[AWS_FIELDS.VPC_CIDR],
                vpcID: '',
                // azs: createAZList(awsData[AWS_FIELDS.SELECTED_AZ_OBJECTS]),
                azs: createAZPayLoadObject(),
            },
            enableAuditLogging: awsData[AWS_FIELDS.ENABLE_AUDIT_LOGGING],
            networking: {
                networkName: '',
                clusterDNSName: '',
                clusterNodeCIDR: '',
                clusterServiceCIDR: awsData[AWS_FIELDS.CLUSTER_SERVICE_CIDR],
                clusterPodCIDR: awsData[AWS_FIELDS.CLUSTER_POD_CIDR],
                cniType: awsData[AWS_FIELDS.CLUSTER_NETWORKING_CNI_PROVIDER],
            },
            ceipOptIn: awsData[AWS_FIELDS.ENABLE_CEIP_PARTICIPATION],
            labels: {},
            // TODO: define a default OS image to set via aws.defaults.tsx
            os: awsData[AWS_FIELDS.OS_IMAGE],
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

    const createAZPayLoadObject = () => {
        const azs: { [key: string]: string }[] = [];
        const defaultAZ: { [key: string]: string }[] = awsState[STORE_SECTION_RESOURCES][AWS_FIELDS.DEFAULT_AZ];
        defaultAZ.forEach((az) => {
            const azObject: { [key: string]: string } = {
                name: awsState[STORE_SECTION_FORM][az.name],
                workNodeType: awsState[STORE_SECTION_FORM][az.workNodeType],
                publicSubnetID: '',
                privateSubnetID: '',
            };
            azs.push(azObject);
        });

        return azs;
    };

    return {
        getAwsRequestPayload,
        deployOnAws,
    };
};
export default useAwsDeployment;
