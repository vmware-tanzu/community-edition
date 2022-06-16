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
<<<<<<< HEAD
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
=======
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea

// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
<<<<<<< HEAD
import { isValidCidr } from '../../../../shared/validations/Validation.service';
=======
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

interface FormInputs {
<<<<<<< HEAD
    NODE_HOST_PORT_MAPPING: string;
=======
    CLUSTER_NAME: string;
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea
    CLUSTER_SERVICE_CIDR: string;
    CLUSTER_POD_CIDR: string;
}

<<<<<<< HEAD
const unmanagedClusterNetworkSettingStepFormSchema = yup
    .object({
        CLUSTER_SERVICE_CIDR: yup
            .string()
            .nullable()
            .required('Please enter a CIDER for your cluster service')
            .test('', 'Cluster name must contain only lower case letters and hyphen', (value) => value !== null && isValidCidr(value)),
        CLUSTER_POD_CIDR: yup
            .string()
            .nullable()
            .required('Please enter a CIDER for your cluster pod')
            .test('', 'Cluster name must contain only lower case letters and hyphen', (value) => value !== null && isValidCidr(value)),
    })
    .required();

=======
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea
const unmanagedClusterProviders = [
    {
        label: 'calico',
        value: 'CALICO',
    },
    {
        label: 'antrea',
        value: 'ANTREA',
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
<<<<<<< HEAD
    } = useForm<FormInputs>({ resolver: yupResolver(unmanagedClusterNetworkSettingStepFormSchema) });
=======
    } = useForm<FormInputs>();
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea

    const handleFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, event.target.name, event.target.value, currentStep, errors);
        }
    };

    const [selectedProvider, setSelectedProvider] = useState('CALICO');

    const handleProviderChange = (event: ChangeEvent<HTMLSelectElement>) => {
        setSelectedProvider(event.target.value);
    };

    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <div cds-layout="p-b:lg" cds-text="title">
                Network settings
            </div>
            <div cds-layout="grid">
                <CdsFormGroup cds-layout="col@sm:8">
                    <div cds-layout="vertical gap:lg">
                        {ClusterProvider()}
                        {ClusterCidr()}
                        {NodeHostPortMapping()}
<<<<<<< HEAD
                        <CdsButton cds-layout="col:start-1" status="success">
                            <CdsIcon shape="cluster" size="sm"></CdsIcon>
                            Create Unmanaged cluster
=======

                        <CdsButton cds-layout="col:start-1" status="success">
                            <CdsIcon shape="cluster" size="sm"></CdsIcon>
                            Create Management cluster
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea
                        </CdsButton>
                    </div>
                </CdsFormGroup>
            </div>
        </div>
    );

    function ClusterProvider() {
        return (
            <CdsRadioGroup layout="vertical-inline" onChange={handleProviderChange}>
                <label>
                    Container Network Interface (CNI) provider <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
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
        );
    }

    function ClusterCidr() {
        return (
            <div cds-layout="horizontal gap:xl">
                <div cds-layout="col:4">
                    <CdsInput layout="vertical" control-width="shrink">
                        <label cds-layout="p-b:md">
                            Cluster service CIDR <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                        </label>
                        <input
                            {...register('CLUSTER_SERVICE_CIDR')}
                            placeholder="CLUSTER SERVICE CIDR"
                            onChange={handleFieldChange}
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
                            onChange={handleFieldChange}
                            defaultValue="100.96.0.0/11"
                        ></input>
                        {errors['CLUSTER_POD_CIDR'] && (
                            <CdsControlMessage status="error">{errors['CLUSTER_POD_CIDR'].message}</CdsControlMessage>
                        )}
                    </CdsInput>
                </div>
            </div>
        );
    }

    function NodeHostPortMapping() {
        return (
            <div cds-layout="grid">
                <div cds-layout="col:6">
                    <CdsInput>
                        <label cds-layout="horizontal p-b:md">
                            Node to host port mapping <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                        </label>
                        <input
<<<<<<< HEAD
                            {...register('NODE_HOST_PORT_MAPPING')}
=======
                            {...register('CLUSTER_NAME')}
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea
                            placeholder="Cluster name"
                            onChange={handleFieldChange}
                            defaultValue={'127.0.0.1:80:80/tcp'}
                        ></input>
<<<<<<< HEAD
                        {errors['NODE_HOST_PORT_MAPPING'] && (
                            <CdsControlMessage status="error">{errors['NODE_HOST_PORT_MAPPING'].message}</CdsControlMessage>
                        )}
=======
                        {errors['CLUSTER_NAME'] && <CdsControlMessage status="error">{errors['CLUSTER_NAME'].message}</CdsControlMessage>}
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea
                        <CdsControlMessage className="description" cds-layout="m-t:sm">
                            Ports to map between container node and the host (format: <q>127.0.0.1:80:80/tcp</q>, <q>80:80/tcp</q>,{' '}
                            <q>80:80</q>, or just <q>80</q>)
                        </CdsControlMessage>
                    </CdsInput>
                </div>
            </div>
        );
    }
}

export default UnmanagedClusterNetworkSettings;
