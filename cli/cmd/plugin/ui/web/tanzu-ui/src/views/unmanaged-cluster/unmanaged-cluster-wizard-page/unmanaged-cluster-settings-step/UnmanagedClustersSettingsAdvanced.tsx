// React imports
import React, { ChangeEvent, useState, useContext } from 'react';
import { useForm, SubmitHandler } from 'react-hook-form';

// Library imports
import { ClarityIcons, blockIcon, blocksGroupIcon, clusterIcon } from '@cds/core/icon';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { CdsIcon } from '@cds/react/icon';
import { CdsButton } from '@cds/react/button';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { isK8sCompliantString } from '../../../../shared/validations/Validation.service';
import { UmcStore } from '../../../../state-management/stores/Store.umc';
import { UNMANAGED_CLUSTER_FIELDS } from '../UnmanagedCluster.constants';
import { UNMANAGED_PLACEHOLDER_VALUES } from '../../../../shared/constants/defaults/unmanaged.defaults';
import { K8sProviders } from '../../../../shared/constants/K8sProviders.constants';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

const unmanagedClusterAdvancedSettingStepFormSchema = yup
    .object({
        CLUSTER_NAME: yup
            .string()
            .nullable()
            .required('Please enter a name for your unmanaged cluster')
            .test(
                '',
                'Cluster name must contain only lower case letters and hyphen',
                (value) => value !== null && isK8sCompliantString(value)
            ),
    })
    .required();

interface FormInputs {
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]: string;
    [UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT]: string;
    [UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT]: string;
}

const unmanagedClusterProviders = [
    {
        label: K8sProviders.KIND,
        value: K8sProviders.KIND,
    },
    {
        label: K8sProviders.MINIKUBE,
        value: K8sProviders.MINIKUBE,
    },
];

function UnmanagedClusterSettingsAdvanced(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const { umcState } = useContext(UmcStore);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<FormInputs>({ resolver: yupResolver(unmanagedClusterAdvancedSettingStepFormSchema) });

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const handleClusterNameChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME, event.target.value, currentStep, errors);
        }
    };

    const handleControlPlaneNodeCountChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT, event.target.value, currentStep, errors);
        }
    };

    const handleWorkerNodeCountChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT, event.target.value, currentStep, errors);
        }
    };

    const [selectedProvider, setSelectedProvider] = useState(umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROVIDER]);

    const handleProviderChange = (event: ChangeEvent<HTMLSelectElement>) => {
        setSelectedProvider(event.target.value);
    };

    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <div cds-layout="p-b:lg" cds-text="title">
                Cluster settings
            </div>
            <div cds-layout="grid">
                <div cds-layout="col@sm:8">
                    <div cds-layout="vertical gap:lg">
                        <div cds-layout="grid gap:md">
                            <div cds-layout="col@sm:6">{ClusterName()}</div>
                        </div>
                        {ClusterNodeCountSelect()}
                        {ClusterProvider()}
                        <div cds-layout="horizontal gap:md">
                            <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );

    function ClusterProvider() {
        return (
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
        );
    }

    function ClusterNodeCountSelect() {
        const errorControlPlaneNodeCount = errors[UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT];
        const errorWorkerNodeCount = errors[UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT];

        return (
            <div cds-layout="grid gap:lg">
                <div cds-layout="col:4">
                    <CdsInput layout="vertical" controlWidth="shrink">
                        <label cds-layout="p-b:md">Control Plane Node Count</label>
                        <input
                            {...register(UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT)}
                            placeholder={UNMANAGED_PLACEHOLDER_VALUES[UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT]}
                            onChange={handleControlPlaneNodeCountChange}
                            defaultValue={umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT]}
                        ></input>
                        {errorControlPlaneNodeCount && (
                            <CdsControlMessage status="error">{errorControlPlaneNodeCount.message}</CdsControlMessage>
                        )}
                    </CdsInput>
                    <p className="description" cds-layout="m-t:sm">
                        The number of control plane nodes to deploy; default is 1
                    </p>
                </div>
                <div cds-layout="col:4">
                    <CdsInput layout="vertical" controlWidth="shrink">
                        <label cds-layout="p-b:md">Worker Node Count</label>
                        <input
                            {...register(UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT)}
                            placeholder={UNMANAGED_PLACEHOLDER_VALUES[UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT]}
                            onChange={handleWorkerNodeCountChange}
                            defaultValue={umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT]}
                        ></input>
                        {errorWorkerNodeCount && <CdsControlMessage status="error">{errorWorkerNodeCount.message}</CdsControlMessage>}
                    </CdsInput>
                    <p className="description" cds-layout="m-t:sm">
                        The number of worker nodes to deploy; default is 0
                    </p>
                </div>
            </div>
        );
    }
    function ClusterName() {
        const errorClusterName = errors[UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME];
        return (
            <div>
                <CdsInput layout="vertical">
                    <label cds-layout="p-b:xs" cds-text="section">
                        Cluster name
                    </label>
                    <input
                        {...register(UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME)}
                        placeholder={UNMANAGED_PLACEHOLDER_VALUES[UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]}
                        onChange={handleClusterNameChange}
                        defaultValue={umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]}
                    ></input>
                    {errorClusterName && <CdsControlMessage status="error">{errorClusterName.message}</CdsControlMessage>}
                </CdsInput>
                <div>
                    <p className="description" cds-layout="m-t:sm">
                        Can only contain lowercase alphanumeric characters and dashes.
                    </p>
                    <p className="description" cds-layout="m-t:sm">
                        The name will be used to reference your cluster in the Tanzu CLI and kubectl.
                    </p>
                </div>
            </div>
        );
    }
}

export default UnmanagedClusterSettingsAdvanced;
