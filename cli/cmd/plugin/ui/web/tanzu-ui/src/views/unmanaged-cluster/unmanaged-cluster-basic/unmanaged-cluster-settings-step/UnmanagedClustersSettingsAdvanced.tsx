// React imports
import React, { ChangeEvent, useState, useContext } from 'react';
import { FormProvider, useForm, SubmitHandler } from 'react-hook-form';

// Library imports
import { blockIcon, blocksGroupIcon, ClarityIcons, clusterIcon } from '@cds/core/icon';
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsInput } from '@cds/react/input';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import { ClusterName, clusterNameValidation } from '../../../../shared/components/FormInputComponents/ClusterName/ClusterName';
import { FormAction } from '../../../../shared/types/types';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { K8sProviders } from '../../../../shared/constants/K8sProviders.constants';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { UmcStore } from '../../../../state-management/stores/Store.umc';
import { UNMANAGED_CLUSTER_FIELDS } from '../../unmanaged-cluster-common/UnmanagedCluster.constants';
import { UNMANAGED_PLACEHOLDER_VALUES } from '../../unmanaged-cluster-common/unmanaged.defaults';
import UseUpdateTabStatus from '../../../../shared/components/wizard/UseUpdateTabStatus.hooks';

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

interface FormInputs {
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]: string;
    [UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT]: number;
    [UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT]: number;
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

function createYupSchemaObject() {
    return {
        [UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]: clusterNameValidation(),
        [UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT]: yup.string().nullable().required('Please enter a control plane node count'),
        [UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT]: yup.string().nullable().required('Please enter a worker node count'),
    };
}
const unmanagedClusterAdvancedSettingStepFormSchema = yup.object(createYupSchemaObject()).required();

function UnmanagedClusterSettingsAdvanced(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm, updateTabStatus } = props;
    const { umcState, umcDispatch } = useContext(UmcStore);
    const methods = useForm<FormInputs>({
        resolver: yupResolver(unmanagedClusterAdvancedSettingStepFormSchema),
        mode: 'all',
    });
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = methods;

    // update tab status bar
    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }

    const canContinue = (): boolean => {
        return Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const onFieldChange = (field: UNMANAGED_CLUSTER_FIELDS, data: string) => {
        umcDispatch({
            type: INPUT_CHANGE,
            field,
            payload: data,
        } as FormAction);
    };

    const handleFieldChange = (event: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const fieldName = event.target.name as UNMANAGED_CLUSTER_FIELDS;
        const newValue = event.target.value;
        umcDispatch({
            type: INPUT_CHANGE,
            field: fieldName,
            payload: newValue,
        } as FormAction);
    };

    return (
        <FormProvider {...methods}>
            <div className="cluster-settings-container" cds-layout="m:lg">
                <div cds-layout="p-b:lg" cds-text="title">
                    Cluster settings
                </div>
                <div cds-layout="grid gap:lg" key="section-holder">
                    <div cds-layout="col:6" key="cluster-name-section">
                        <ClusterName
                            field={UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME}
                            clusterNameChange={(value) => {
                                onFieldChange(UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME, value);
                            }}
                            placeholderClusterName={UNMANAGED_PLACEHOLDER_VALUES.CLUSTER_NAME}
                        />
                    </div>
                    <div cds-layout="col:8" key="cluster-name-section">
                        {ClusterNodeCountSelect()}
                    </div>
                    <div cds-layout="col:6" key="cluster-name-section">
                        {ClusterProvider()}
                    </div>
                </div>
                <div cds-layout="p-t:lg">
                    <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                        NEXT
                    </CdsButton>
                </div>
            </div>
        </FormProvider>
    );

    function ClusterProvider() {
        return (
            <CdsRadioGroup
                layout="vertical-inline"
                onChange={(e: any) => onFieldChange(UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROVIDER, e.target.value)}
            >
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
                                checked={
                                    umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROVIDER] ===
                                    unmanagedClusterProviders.value
                                }
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
                            onChange={handleFieldChange}
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
                            onChange={handleFieldChange}
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
}

export default UnmanagedClusterSettingsAdvanced;
