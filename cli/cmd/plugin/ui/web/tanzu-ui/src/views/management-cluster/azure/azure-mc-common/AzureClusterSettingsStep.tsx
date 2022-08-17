// React imports
import React, { useContext, useState, useEffect } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
import * as yup from 'yup';

// App imports
import { AzureVirtualMachine } from '../../../../swagger-api';
import { AzureStore } from '../store/Azure.store.mc';
import { AZURE_FIELDS, AZURE_NODE_PROFILE_NAMES } from '../azure-mc-basic/AzureManagementClusterBasic.constants';
import { ClusterName, clusterNameValidation } from '../../../../shared/components/FormInputComponents/ClusterName/ClusterName';
import { FormAction } from '../../../../shared/types/types';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import {
    NodeInstanceType,
    nodeInstanceTypeValidation,
    NodeProfile,
} from '../../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import OsImageSelect from '../../../../shared/components/FormInputComponents/OsImageSelect/OsImageSelect';
import PageNotification, { Notification } from '../../../../shared/components/PageNotification/PageNotification';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import UseUpdateTabStatus from '../../../../shared/components/wizard/UseUpdateTabStatus.hooks';
import { getResource } from '../../../providers/azure/AzureResources.reducer';
// NOTE: icons must be imported
const nodeInstanceTypes: NodeInstanceType[] = [
    {
        id: AZURE_NODE_PROFILE_NAMES.SINGLE_NODE,
        label: 'Single node',
        icon: 'block',
        description: 'Create a single control plane node with a Standard_D2s_v3 instance type',
    },
    {
        id: AZURE_NODE_PROFILE_NAMES.HIGH_AVAILABILITY,
        label: 'High availability',
        icon: 'blocks-group',
        description: 'Create a multi-node control plane with a Standard_D2s_v3 instance type',
    },
    {
        id: AZURE_NODE_PROFILE_NAMES.PRODUCTION_READY,
        label: 'Production-ready (High availability)',
        icon: 'blocks-group',
        isSolidIcon: true,
        description: 'Create a multi-node control plane with a Standard_D4s_v3 instance type',
    },
];

type AZURE_CLUSTER_SETTING_STEP_FIELDS = AZURE_FIELDS.CLUSTER_NAME | AZURE_FIELDS.NODE_PROFILE | AZURE_FIELDS.OS_IMAGE;

interface AzureClusterSettingFormInputs {
    [AZURE_FIELDS.CLUSTER_NAME]: string;
    [AZURE_FIELDS.NODE_PROFILE]: string;
    [AZURE_FIELDS.OS_IMAGE]: string;
}

function createYupSchemaObject() {
    return {
        [AZURE_FIELDS.NODE_PROFILE]: nodeInstanceTypeValidation(),
        [AZURE_FIELDS.CLUSTER_NAME]: clusterNameValidation(),
    };
}

export function AzureClusterSettingsStep(props: Partial<StepProps>) {
    const { updateTabStatus, currentStep, submitForm, goToStep } = props;
    const { azureState, azureDispatch } = useContext(AzureStore);
    const [notification, setNotification] = useState<Notification | null>(null);
    const azureClusterSettingsFormSchema = yup.object(createYupSchemaObject()).required();
    const methods = useForm<AzureClusterSettingFormInputs>({
        resolver: yupResolver(azureClusterSettingsFormSchema),
        mode: 'all',
    });

    const {
        handleSubmit,
        formState: { errors },
        setValue,
    } = methods;

    // update tab status bar
    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }
    const osImages = (getResource(AZURE_FIELDS.OS_IMAGE, azureState) || []) as AzureVirtualMachine[];

    let initialSelectedNodeProfileId = azureState[AZURE_FIELDS.NODE_PROFILE];
    if (!initialSelectedNodeProfileId) {
        initialSelectedNodeProfileId = nodeInstanceTypes[0].id;
        setValue(AZURE_FIELDS.NODE_PROFILE, initialSelectedNodeProfileId);
    }
    const [selectedInstanceTypeId, setSelectedInstanceTypeId] = useState(initialSelectedNodeProfileId);

    function dismissAlert() {
        setNotification(null);
    }

    const canContinue = (): boolean => {
        return (
            Object.keys(errors).length === 0 &&
            azureState[STORE_SECTION_FORM][AZURE_FIELDS.CLUSTER_NAME] &&
            azureState[STORE_SECTION_FORM][AZURE_FIELDS.NODE_PROFILE]
        );
    };

    const onSubmit: SubmitHandler<AzureClusterSettingFormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            submitForm(currentStep);
            goToStep(currentStep + 1);
        }
    };

    const onFieldChange = (field: AZURE_CLUSTER_SETTING_STEP_FIELDS, data: any) => {
        azureDispatch({
            type: INPUT_CHANGE,
            field,
            payload: data,
        } as FormAction);
    };

    const onClusterNameChange = (clusterName: string) => {
        onFieldChange(AZURE_FIELDS.CLUSTER_NAME, clusterName);
    };

    const onInstanceTypeChange = (instanceType: string) => {
        onFieldChange(AZURE_FIELDS.NODE_PROFILE, instanceType);
        setSelectedInstanceTypeId(instanceType);
    };

    return (
        <FormProvider {...methods}>
            <div className="wizard-content-container">
                <h2 cds-layout="m-t:md m-b:xl" cds-text="title">
                    Azure Management Cluster Settings
                </h2>
                <div cds-layout="grid gap:m" key="section-holder">
                    <div cds-layout="col:4" key="cluster-name-section">
                        <ClusterName
                            field={AZURE_FIELDS.CLUSTER_NAME}
                            clusterNameChange={onClusterNameChange}
                            placeholderClusterName="my-azure-cluster"
                            defaultClusterName={azureState[STORE_SECTION_FORM][AZURE_FIELDS.CLUSTER_NAME]}
                        />
                    </div>
                    <div cds-layout="col:8" key="instance-type-section">
                        <NodeProfile
                            field={AZURE_FIELDS.NODE_PROFILE}
                            nodeInstanceTypes={nodeInstanceTypes}
                            nodeInstanceTypeChange={onInstanceTypeChange}
                            selectedInstanceId={selectedInstanceTypeId}
                        />
                    </div>
                    <div cds-layout="col:6">
                        <PageNotification notification={notification} closeCallback={dismissAlert}></PageNotification>
                    </div>
                    <div cds-layout="col:12">
                        <OsImageSelect
                            osImageTitle={'Azure Machine Image'}
                            images={osImages}
                            field={AZURE_FIELDS.OS_IMAGE}
                            onOsImageSelected={(value) => {
                                onFieldChange(AZURE_FIELDS.OS_IMAGE, value);
                            }}
                        />
                    </div>
                </div>
                <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                    NEXT
                </CdsButton>
            </div>
        </FormProvider>
    );
}
