// React imports
import React, { ChangeEvent, useState } from 'react';

// Library imports
import { ClarityIcons, blockIcon, blocksGroupIcon, clusterIcon } from '@cds/core/icon';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { useForm } from 'react-hook-form';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { CdsIcon } from '@cds/react/icon';
import { CdsButton } from '@cds/react/button';

// App imports
import './ManagementClusterSettings.scss';
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';
import { INPUT_CHANGE } from '../../../state-management/actions/Form.actions';
import { StepProps } from '../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../state-management/reducers/Form.reducer';

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

interface FormInputs {
    CLUSTER_NAME: string;
}
interface MCSettings extends StepProps {
    message?: string;
    deploy: () => void;
    defaultData?: { [key: string]: any };
}

const nodeProfiles = [
    {
        label: 'Single node',
        icon: 'block',
        message: 'Create one control plane node with a general purpose instance type in a single region.',
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
function ManagementClusterSettings(props: Partial<MCSettings>) {
    const { handleValueChange, currentStep, deploy, defaultData, message } = props;
    const {
        register,
        formState: { errors },
    } = useForm<FormInputs>();
    const [selectedProfile, setSelectedProfile] = useState('SINGLE_NODE');
    const handleNodeProfileChange = (event: ChangeEvent<HTMLSelectElement>) => {
        setSelectedProfile(event.target.value);
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'NODE_PROFILE', event.target.value, currentStep, errors);
        }
    };

    const handleClusterNameChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'CLUSTER_NAME', event.target.value, currentStep, errors);
        }
    };
    const handleMCCreation = () => {
        if (deploy) {
            deploy();
        }
    };
    return (
        <div className="cluster-settings-container" cds-layout="p:lg">
            <h2 cds-layout="m-t:md m-b:xl" cds-text="title">
                Management Cluster settings
            </h2>
            <div cds-layout="grid gap:md">
                <div cds-layout="col@sm:4">
                    <CdsInput>
                        <label cds-layout="p-b:md">Cluster name</label>
                        <input
                            {...register('CLUSTER_NAME')}
                            aria-label="cluster name"
                            placeholder="Cluster name"
                            onChange={handleClusterNameChange}
                            defaultValue={defaultData ? defaultData[STORE_SECTION_FORM]?.CLUSTER_NAME : undefined}
                        ></input>
                        {errors['CLUSTER_NAME'] && <CdsControlMessage status="error">{errors['CLUSTER_NAME'].message}</CdsControlMessage>}
                    </CdsInput>
                    <p className="description" cds-layout="m-t:sm">
                        Can only contain lowercase alphanumeric characters and dashes.
                        <br></br>
                        <br></br>
                        The name will be used to reference your cluster in the Tanzu CLI and kubectl.
                    </p>
                </div>
                <div cds-layout="col@sm:8 p-l:xl">
                    {message && (
                        <CdsAlertGroup status="info" cds-layout="m-b:md">
                            <CdsAlert>{message}</CdsAlert>
                        </CdsAlertGroup>
                    )}
                    {!message && (
                        <CdsRadioGroup layout="vertical" onChange={handleNodeProfileChange}>
                            <label>Select a control plane node profile</label>
                            {nodeProfiles.map((nodeProfile, index) => {
                                return (
                                    <CdsRadio cds-layout="m:lg m-l:xl p-b:sm" key={index} data-testid="cds-radio">
                                        <label>
                                            {nodeProfile.label}
                                            <CdsIcon
                                                shape={nodeProfile.icon}
                                                size="md"
                                                className={selectedProfile === nodeProfile.value ? 'node-icon selected' : 'node-icon'}
                                                solid={nodeProfile.isSolid}
                                            ></CdsIcon>
                                            <div className="radio-message">{nodeProfile.message}</div>
                                        </label>
                                        <input
                                            type="radio"
                                            key={index}
                                            value={nodeProfile.value}
                                            checked={selectedProfile === nodeProfile.value}
                                            readOnly
                                        />
                                    </CdsRadio>
                                );
                            })}
                        </CdsRadioGroup>
                    )}
                </div>
                <div cds-layout="grid col:12 p-t:lg">
                    <CdsButton cds-layout="col:start-1" status="success" onClick={handleMCCreation}>
                        <CdsIcon shape="cluster" size="sm"></CdsIcon>
                        Create Management cluster
                    </CdsButton>
                    <CdsButton cds-layout="col:end-12" action="flat">
                        View configuration details
                    </CdsButton>
                </div>
            </div>
        </div>
    );
}

export default ManagementClusterSettings;
