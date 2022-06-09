// React imports
import React, { ChangeEvent, useState } from 'react';

// Library imports
import { ClarityIcons, blockIcon, blocksGroupIcon, clusterIcon } from '@cds/core/icon';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { useForm } from 'react-hook-form';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { CdsIcon } from '@cds/react/icon';
import { CdsButton } from '@cds/react/button';

// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

interface FormInputs {
    CLUSTER_NAME: string;
    CLUSTER_SERVICE_CIDR: string;
    CLUSTER_POD_CIDR: string;
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
function UnmanagedClusterNetworkSettings(props: Partial<StepProps>) {
    const { handleValueChange, currentStep } = props;

    const {
        register,
        formState: { errors },
    } = useForm<FormInputs>();

    const handleClusterNameChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'CLUSTER_NAME', event.target.value, currentStep, errors);
        }
    };

    const handleClusterServiceCidrChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'CLUSTER_SERVICE_CIDR', event.target.value, currentStep, errors);
        }
    };

    const handleClusterPodCidrChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'WORKER_NODE_COUNT', event.target.value, currentStep, errors);
        }
    };

    const [selectedProvider, setSelectedProvider] = useState('CALICO');

    const handleProviderChange = (event: ChangeEvent<HTMLSelectElement>) => {
        setSelectedProvider(event.target.value);
    };
    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <h3>Network settings</h3>
            <h5>
                Your cluster will use the following networking solution for pod-to-pod communications.
                <br></br>
                <br></br>
                Container Network Interface (CNI) provider: calico
            </h5>
            <div cds-layout="grid">
                <CdsFormGroup cds-layout="col@sm:8">
                    <div cds-layout="col@sm:12 vertical gap:lg">
                        <CdsRadioGroup layout="vertical-inline" onChange={handleProviderChange}>
                            <label>
                                Cluster provider <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                            </label>
                            {unmanagedClusterProviders.map((unmanagedClusterProviders, index) => {
                                return (
                                    <CdsRadio cds-layout="m:md m-l:none" key={index}>
                                        <label>{unmanagedClusterProviders.label}</label>
                                        <input
                                            type="radio"
                                            key={index}
                                            value={unmanagedClusterProviders.value}
                                            checked={selectedProvider === unmanagedClusterProviders.value}
                                            readOnly
                                        />
                                    </CdsRadio>
                                );
                            })}
                        </CdsRadioGroup>
                        <div cds-layout="vertical gap:xl">
                            <div cds-layout="horizontal gap:xl">
                                <div cds-layout="col:4">
                                    <CdsInput layout="vertical" control-width="shrink">
                                        <label cds-layout="p-b:md">
                                            Cluster service CIDR <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                                        </label>
                                        <input
                                            {...register('CLUSTER_SERVICE_CIDR')}
                                            placeholder="CLUSTER SERVICE CIDR"
                                            onChange={handleClusterServiceCidrChange}
                                            defaultValue="100.64.0.0/13"
                                        ></input>
                                        {errors['CLUSTER_SERVICE_CIDR'] && (
                                            <CdsControlMessage status="error">{errors['CLUSTER_SERVICE_CIDR'].message}</CdsControlMessage>
                                        )}
                                    </CdsInput>
                                </div>
                                <div cds-layout="col:4">
                                    <CdsInput layout="vertical" control-width="shrink">
                                        <label cds-layout="p-b:md">
                                            Cluster pod CIDR <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                                        </label>
                                        <input
                                            {...register('CLUSTER_POD_CIDR')}
                                            placeholder="CLUSTER POD CIDR"
                                            onChange={handleClusterPodCidrChange}
                                            defaultValue="100.96.0.0/11"
                                        ></input>
                                        {errors['CLUSTER_POD_CIDR'] && (
                                            <CdsControlMessage status="error">{errors['CLUSTER_POD_CIDR'].message}</CdsControlMessage>
                                        )}
                                    </CdsInput>
                                </div>
                            </div>
                            <div cds-layout="col:4">
                                <CdsInput>
                                    <label cds-layout="p-b:md">
                                        Node to host port mapping <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                                    </label>
                                    <input
                                        {...register('CLUSTER_NAME')}
                                        placeholder="Cluster name"
                                        onChange={handleClusterNameChange}
                                        defaultValue={'127.0.0.1:80:80/tcp'}
                                    ></input>
                                    {errors['CLUSTER_NAME'] && (
                                        <CdsControlMessage status="error">{errors['CLUSTER_NAME'].message}</CdsControlMessage>
                                    )}
                                    <CdsControlMessage className="description" cds-layout="m-t:sm">
                                        Can only contain lowercase alphanumeric characters and dashes.
                                        <br></br>
                                        <br></br>
                                        The name will be used to reference your cluster in the Tanzu CLI and kubectl.
                                    </CdsControlMessage>
                                </CdsInput>
                            </div>
                        </div>
                    </div>
                    <div cds-layout="grid col:12 p-t:lg">
                        <CdsButton cds-layout="col:start-1" status="success">
                            <CdsIcon shape="cluster" size="sm"></CdsIcon>
                            Create Management cluster
                        </CdsButton>
                    </div>
                </CdsFormGroup>
            </div>
        </div>
    );
}

export default UnmanagedClusterNetworkSettings;
