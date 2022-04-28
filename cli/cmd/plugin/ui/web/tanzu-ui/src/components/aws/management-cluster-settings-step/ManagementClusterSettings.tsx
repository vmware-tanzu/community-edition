// React imports
import React, { ChangeEvent, useContext, useState } from 'react';

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
import { StepProps } from '../../../shared/components/wizard/Wizard';
import { AwsStore } from '../../../state-management/stores/Store.aws';
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
                'CLUSTER_NAME',
                event.target.value,
                currentStep,
                errors
            );
        }
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
                    <div cds-layout="m-t:xxl p-t:lg">
                        <CdsButton
                            status="success"
                            onClick={() => {
                                console.log(
                                    '//TODO: call management cluster creation api and navigate to progress page'
                                );
                            }}
                        >
                            <CdsIcon shape="cluster" size="sm"></CdsIcon>
                            Create Management cluster
                        </CdsButton>
                    </div>
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
            </div>
        </div>
    );
}

export default ManagementClusterSettings;
