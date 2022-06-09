// React imports
import React, { ChangeEvent, useContext, useState } from 'react';

// Library imports
import { ClarityIcons, blockIcon, blocksGroupIcon, clusterIcon } from '@cds/core/icon';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { useForm } from 'react-hook-form';

// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

interface FormInputs {
    CLUSTER_NAME: string;
    CONTROL_PLANE_NODE_COUNT: string;
    WORKER_NODE_COUNT: string;
}

const unmanagedClusterProviders = [
    {
        label: 'calico',
        value: 'CALICO',
    },
    {
        label: 'anthrea',
        value: 'ANTHREA',
    },
    {
        label: 'none',
        value: 'NONE',
    },
];
function UnmanagedClusterSettingsAdvanced(props: Partial<StepProps>) {
    const { handleValueChange, currentStep } = props;
    const {
        register,
        formState: { errors },
    } = useForm<FormInputs>();

    const [selectedProvider, setSelectedProvider] = useState('KIND');

    const handleClusterNameChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'CLUSTER_NAME', event.target.value, currentStep, errors);
        }
    };

    const handleControlPlaneNodeCountChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'CONTROL_PLANE_NODE_COUNT', event.target.value, currentStep, errors);
        }
    };

    const handleWorkerNodeCountChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'WORKER_NODE_COUNT', event.target.value, currentStep, errors);
        }
    };

    const handleProviderChange = (event: ChangeEvent<HTMLSelectElement>) => {
        setSelectedProvider(event.target.value);
    };
    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <h2>Cluster settings</h2>
            <div cds-layout="grid gap:md">
                <div cds-layout="col@sm:4">
                    <CdsInput>
                        <label cds-layout="p-b:md">Cluster name</label>
                        <input
                            {...register('CLUSTER_NAME')}
                            placeholder="Cluster name"
                            onChange={handleClusterNameChange}
                            defaultValue="Test Cluster"
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
                <div cds-layout="col:12">
                    <div cds-layout="col:4">
                        <CdsInput>
                            <label cds-layout="p-b:md">Control Plane Node Count</label>
                            <input
                                {...register('CONTROL_PLANE_NODE_COUNT')}
                                placeholder="Control Plane Node Count"
                                onChange={handleControlPlaneNodeCountChange}
                                defaultValue="1"
                            ></input>
                            {errors['CONTROL_PLANE_NODE_COUNT'] && (
                                <CdsControlMessage status="error">{errors['CONTROL_PLANE_NODE_COUNT'].message}</CdsControlMessage>
                            )}
                        </CdsInput>
                        <p className="description" cds-layout="m-t:sm">
                            The number of control plane nodes to deploy; default is 1
                        </p>
                    </div>
                    <div cds-layout="col:4">
                        <CdsInput>
                            <label cds-layout="p-b:md">Worker Node Count</label>
                            <input
                                {...register('WORKER_NODE_COUNT')}
                                placeholder="Worker Node Count"
                                onChange={handleWorkerNodeCountChange}
                                defaultValue="0"
                            ></input>
                            {errors['WORKER_NODE_COUNT'] && (
                                <CdsControlMessage status="error">{errors['WORKER_NODE_COUNT'].message}</CdsControlMessage>
                            )}
                        </CdsInput>
                        <p className="description" cds-layout="m-t:sm">
                            The number of worker nodes to deploy; default is O
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default UnmanagedClusterSettingsAdvanced;
