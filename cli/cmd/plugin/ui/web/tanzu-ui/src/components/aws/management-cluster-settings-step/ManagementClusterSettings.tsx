// React imports
import React, { ChangeEvent, useContext, useState } from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import {
    ClarityIcons,
    blockIcon,
    blocksGroupIcon,
    clusterIcon,
} from '@cds/core/icon';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { useForm } from 'react-hook-form';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { CdsIcon } from '@cds/react/icon';
import { CdsButton } from '@cds/react/button';

// App imports
import { AwsStore } from '../../../state-management/stores/Store.aws';
import { INPUT_CHANGE } from '../../../state-management/actions/Form.actions';
import { NavRoutes } from '../../../shared/constants/NavRoutes.constants';
import { AwsService } from '../../../swagger-api/services/AwsService';
import { AWSManagementClusterParams, ConfigFileInfo, IdentityManagementConfig } from '../../../swagger-api';
import { DEPLOYMENT_STATUS_CHANGED } from '../../../state-management/actions/Deployment.actions';
import { DeploymentStates, DeploymentTypes } from '../../../shared/constants/Deployment.constants';
import { Providers } from '../../../shared/constants/Providers.constants';
import { StepProps } from '../../../shared/components/wizard/Wizard';
import { Store } from '../../../state-management/stores/Store';
import { TOGGLE_APP_STATUS } from '../../../state-management/actions/Ui.actions';

import './ManagementClusterSettings.scss';

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

interface FormInputs {
    CLUSTER_NAME: string;
}

const nodeProfiles = [
    {
        label: 'Single node',
        icon: 'block',
        message:
            'Create one control plane-node with a general purpose instance type in a single region.',
        value: 'SINGLE_NODE',
    },
    {
        label: 'High availability',
        icon: 'blocks-group',
        message: `Create a multi-node control plane with general purpose instance types in a single,
            or multiple, regions. Provides a fault-tolerant control plane.`,
        value: 'HIGH_AVAILABILITY',
    },
    {
        label: 'Production-ready (High availability)',
        icon: 'blocks-group',
        isSolid: true,
        message: `Create a multi-node control plane with
            recommended, performant, instance types
            across multiple regions. Recommended for
            production workloads.`,
        value: 'PRODUCTION_READY',
    },
];
function ManagementClusterSettings(props: Partial<StepProps>) {
    const { handleValueChange, currentStep } = props;
    const { awsState } = useContext(AwsStore);
    const { dispatch } = useContext(Store);
    
    const navigate = useNavigate();

    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };
    const {
        register,
        formState: { errors },
    } = useForm<FormInputs>();
    const [selectedProfile, setSelectedProfile] = useState('SINGLE_NODE');
    const handleNodeProfileChange = (event: ChangeEvent<HTMLSelectElement>) => {
        setSelectedProfile(event.target.value);
    };

    const handleClusterNameChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(
                INPUT_CHANGE,
                'CLUSTER_NAME',
                event.target.value,
                currentStep,
                errors
            );
        }
    };
    const handleMCCreation = () => {

        let awsClusterParams: AWSManagementClusterParams = {
            awsAccountParams: {
                profileName: awsState.data.PROFILE,
                sessionToken: awsState.data.SESSION_TOKEN,
                region: awsState.data.REGION,
                accessKeyID: awsState.data.ACCESS_KEY_ID,
                secretAccessKey: awsState.data.SECRET_ACCESS_KEY
            },
            loadbalancerSchemeInternal: false,
            sshKeyName: awsState.data.EC2_KEY_PAIR,
            createCloudFormationStack: false,
            clusterName: awsState.data.CLUSTER_NAME,
            controlPlaneFlavor: awsState.data.CLUSTER_PLAN,
            controlPlaneNodeType: awsState.data.CONTROL,
            bastionHostEnabled: true,
            machineHealthCheckEnabled: true,
            vpc: {
                cidr: awsState.data.VPC_CIDR,
                vpcID: '',
                azs: [
                    {
                        name: 'us-west-2a',
                        workerNodeType: awsState.data.CLUSTER_WORKER_NODE_TYPE,
                        publicSubnetID: '',
                        privateSubnetID: ''
                    }
                ]
            },
            enableAuditLogging: false,
            networking: {
                networkName: '',
                clusterDNSName: '',
                clusterNodeCIDR: '',
                clusterServiceCIDR: awsState.data.CLUSTER_SERVICE_CIDR,
                clusterPodCIDR: awsState.data.CLUSTER_POD_CIDR,
                cniType: 'antrea'
            },
            ceipOptIn: true,
            labels: {},
            os: {
                name: 'ubuntu-20.04-amd64 (ami-0dd0327a3bfaa4dc8)',
                osInfo: {
                    arch: 'amd64',
                    name: 'ubuntu',
                    version: '20.04'
                }
            },
            annotations: {
                description: '',
                location: ''
            },
            identityManagement: {
                idm_type: IdentityManagementConfig.idm_type.NONE
            }
        }

        AwsService.applyTkgConfigForAws(awsClusterParams).then((data: ConfigFileInfo) => {
            const configFilePath = data.path;

            AwsService.createAwsManagementCluster(awsClusterParams).then(() => {
                dispatch({
                    type: DEPLOYMENT_STATUS_CHANGED,
                    payload: {
                        type: DeploymentTypes.MANAGEMENT_CLUSTER,
                        status: DeploymentStates.RUNNING,
                        provider: Providers.AWS,
                        configPath: configFilePath
                    }
                })
            });
        });


        dispatch({
            type: TOGGLE_APP_STATUS
        })

        navigateToProgress();
    };
    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <h3>Management Cluster settings</h3>
            <div cds-layout="grid gap:md">
                <div cds-layout="col@sm:4">
                    <CdsInput>
                        <label cds-layout="p-b:md">Cluster name</label>
                        <input
                            {...register('CLUSTER_NAME')}
                            placeholder="Cluster name"
                            onChange={handleClusterNameChange}
                            defaultValue={awsState.data.CLUSTER_NAME}
                        ></input>
                        {errors['CLUSTER_NAME'] && (
                            <CdsControlMessage status="error">
                                {errors['CLUSTER_NAME'].message}
                            </CdsControlMessage>
                        )}
                    </CdsInput>
                    <p className="description" cds-layout="m-t:sm">
                        Can only contain lowercase alphanumeric characters and
                        dashes.
                        <br></br>
                        <br></br>
                        The name will be used to reference your cluster in the
                        Tanzu CLI and kubectl.
                    </p>
                </div>
                <div cds-layout="col@sm:8 p-l:xl">
                    <CdsRadioGroup
                        layout="vertical"
                        onChange={handleNodeProfileChange}
                    >
                        <label>Select a control plane-node profile</label>
                        {nodeProfiles.map((nodeProfile, index) => {
                            return (
                                <CdsRadio
                                    cds-layout="m:lg m-l:xl p-b:sm"
                                    key={index}
                                >
                                    <label>
                                        {nodeProfile.label}
                                        <CdsIcon
                                            shape={nodeProfile.icon}
                                            size="md"
                                            className={
                                                selectedProfile ===
                                                nodeProfile.value ? 'node-icon selected' : 'node-icon'
                                            }
                                            solid={nodeProfile.isSolid}
                                        ></CdsIcon>
                                        <div className="radio-message">
                                            {nodeProfile.message}
                                        </div>
                                    </label>
                                    <input
                                        type="radio"
                                        key={index}
                                        value={nodeProfile.value}
                                        checked={
                                            selectedProfile ===
                                            nodeProfile.value
                                        }
                                        readOnly
                                    />
                                </CdsRadio>
                            );
                        })}
                    </CdsRadioGroup>
                </div>
                <div cds-layout="grid col:12 p-t:lg">
                    <CdsButton
                        cds-layout="col:start-1"
                        status="success"
                        onClick={handleMCCreation}
                    >
                        <CdsIcon shape="cluster" size="sm"></CdsIcon>
                        Create Management cluster
                    </CdsButton>
                    <CdsButton
                        cds-layout="col:end-12"
                        action="flat"
                    >
                        View configuration details
                    </CdsButton>
                </div>
            </div>
        </div>
    );
}

export default ManagementClusterSettings;
