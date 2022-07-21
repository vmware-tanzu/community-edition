// React imports
import React, { useContext, useState } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
import * as yup from 'yup';

// App imports

import { AZURE_FIELDS } from '../AzureManagementCluster.constants';
import { AzureStore } from '../../../../state-management/stores/Azure.store';
import { ClusterName, clusterNameValidation } from '../../../../shared/components/FormInputComponents/ClusterName/ClusterName';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import {
    NodeInstanceType,
    nodeInstanceTypeValidation,
    NodeProfile,
} from '../../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { SubmitHandler, useForm } from 'react-hook-form';

// NOTE: icons must be imported
const nodeInstanceTypes: NodeInstanceType[] = [
    {
        id: 'single-node',
        label: 'Single node',
        icon: 'block',
        description: 'Create a single control plane node with a Standard_D2s_v3 instance type',
    },
    {
        id: 'high-availability',
        label: 'High availability',
        icon: 'blocks-group',
        description: 'Create a multi-node control plane with a Standard_D2s_v3 instance type',
    },
    {
        id: 'compute-optimized',
        label: 'Production-ready (High availability)',
        icon: 'blocks-group',
        isSolidIcon: true,
        description: 'Create a multi-node control plane with a Standard_D4s_v3 instance type',
    },
];

type AZURE_CLUSTER_SETTING_STEP_FIELDS = AZURE_FIELDS.CLUSTER_NAME | AZURE_FIELDS.INSTANCE_TYPE;

interface AzureClusterSettingFormInputs {
    [AZURE_FIELDS.CLUSTER_NAME]: string;
    [AZURE_FIELDS.INSTANCE_TYPE]: string;
}

export function AzureClusterSettingsStep(props: Partial<StepProps>) {
    const { currentStep, deploy, handleValueChange } = props;
    const { azureState, azureDispatch } = useContext(AzureStore);
    const azureClusterSettingsFormSchema = yup.object(createYupSchemaObject()).required();
    const methods = useForm<AzureClusterSettingFormInputs>({
        resolver: yupResolver(azureClusterSettingsFormSchema),
    });

    const {
        handleSubmit,
        formState: { errors },
        register,
        setValue,
    } = methods;

    let initialSelectedInstanceTypeId = azureState[AZURE_FIELDS.INSTANCE_TYPE];
    if (!initialSelectedInstanceTypeId) {
        initialSelectedInstanceTypeId = nodeInstanceTypes[0].id;
        setValue(AZURE_FIELDS.INSTANCE_TYPE, initialSelectedInstanceTypeId);
    }
    const [selectedInstanceTypeId, setSelectedInstanceTypeId] = useState(initialSelectedInstanceTypeId);

    const canContinue = (): boolean => {
        return Object.keys(errors).length === 0;
    };

    // TODO: just deactivate button until no errors
    const onSubmit: SubmitHandler<AzureClusterSettingFormInputs> = (data) => {
        if (canContinue() && deploy) {
            deploy();
        }
    };

    const onFieldChange = (data: string, field: AZURE_CLUSTER_SETTING_STEP_FIELDS) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, field, data, currentStep, errors);
            setValue(field, data, { shouldValidate: true });
        }
    };

    const onClusterNameChange = (clusterName: string) => {
        onFieldChange(clusterName, AZURE_FIELDS.CLUSTER_NAME);
    };

    const onInstanceTypeChange = (instanceType: string) => {
        onFieldChange(instanceType, AZURE_FIELDS.INSTANCE_TYPE);
        setSelectedInstanceTypeId(instanceType);
    };

    return (
        <div>
            <div className="wizard-content-container">
                <h2 cds-layout="m-t:lg">Azure Management Cluster Settings</h2>
                <div cds-layout="grid gap:m" key="section-holder">
                    <div cds-layout="col:6" key="cluster-name-section">
                        <ClusterName
                            field={AZURE_FIELDS.CLUSTER_NAME}
                            errors={errors}
                            register={register}
                            clusterNameChange={onClusterNameChange}
                            placeholderClusterName={'my-azure-cluster'}
                        />
                    </div>
                    <div cds-layout="col:6" key="instance-type-section">
                        <NodeProfile
                            field={AZURE_FIELDS.INSTANCE_TYPE}
                            nodeInstanceTypes={nodeInstanceTypes}
                            errors={errors}
                            register={register}
                            nodeInstanceTypeChange={onInstanceTypeChange}
                            selectedInstanceId={selectedInstanceTypeId}
                        />
                    </div>
                </div>
                <CdsButton cds-layout="col:start-1" status="success" onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                    <CdsIcon shape="cluster" size="sm"></CdsIcon>
                    Create Management cluster
                </CdsButton>
            </div>
        </div>
    );
}

function createYupSchemaObject() {
    return {
        [AZURE_FIELDS.INSTANCE_TYPE]: nodeInstanceTypeValidation(),
        [AZURE_FIELDS.CLUSTER_NAME]: clusterNameValidation(),
    };
}
