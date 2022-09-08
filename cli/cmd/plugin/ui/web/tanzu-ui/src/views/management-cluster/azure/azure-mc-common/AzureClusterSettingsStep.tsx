// React imports
import React, { useContext, useEffect, useState } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
import * as yup from 'yup';

// App imports
import { AzureDefaults } from './default-service/AzureDefaults.service';
import { AzureVirtualMachine } from '../../../../swagger-api';
import { AzureStore } from '../store/Azure.store.mc';
import { AZURE_FIELDS, AZURE_NODE_PROFILE_NAMES } from '../azure-mc-basic/AzureManagementClusterBasic.constants';
import { ClusterName, clusterNameValidation } from '../../../../shared/components/FormInputComponents/ClusterName/ClusterName';
import { getResource, STORE_SECTION_RESOURCES } from '../../../../state-management/reducers/Resources.reducer';
import { BATCH_SET } from '../../../../state-management/actions/Form.actions';
import {
    NodeProfileType,
    nodeProfileValidation,
    NodeProfile,
} from '../../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import OsImageSelect from '../../../../shared/components/FormInputComponents/OsImageSelect/OsImageSelect';
import PageNotification, { Notification } from '../../../../shared/components/PageNotification/PageNotification';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import UseUpdateTabStatus from '../../../../shared/components/wizard/UseUpdateTabStatus.hooks';

const nodeInstanceTypes: NodeProfileType[] = [
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
        [AZURE_FIELDS.NODE_PROFILE]: nodeProfileValidation(),
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
    // NOTE: we assume that the osImages were set in the store during the credentials step
    const osImages = getResource<AzureVirtualMachine[]>(AZURE_FIELDS.OS_IMAGE, azureState) || [];

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
    const selectedimage = azureState[STORE_SECTION_FORM][AZURE_FIELDS.OS_IMAGE];
    useEffect(() => {
        if (selectedimage) {
            setValue(AZURE_FIELDS.OS_IMAGE, selectedimage.name);
        }
    }, [selectedimage, setValue]);

    const selectedNodeProfile = azureState[STORE_SECTION_FORM][AZURE_FIELDS.NODE_PROFILE];

    useEffect(() => {
        if (selectedNodeProfile) {
            setValue(AZURE_FIELDS.NODE_PROFILE, selectedNodeProfile);
        }
    }, [selectedNodeProfile, setValue]);

    const onSubmit: SubmitHandler<AzureClusterSettingFormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            submitForm(currentStep);
            goToStep(currentStep + 1);
        }
    };

    const onFieldChange = (field: AZURE_CLUSTER_SETTING_STEP_FIELDS, data: any) => {
        const nodeTypes = azureState[STORE_SECTION_RESOURCES][AZURE_FIELDS.NODE_TYPE];
        const nodeProfile = azureState[STORE_SECTION_FORM][AZURE_FIELDS.NODE_PROFILE];
        const nodeType = AzureDefaults.getDefaultNodeType(nodeTypes, nodeProfile);

        azureDispatch({
            type: BATCH_SET,
            payload: {
                [field]: data,
                [AZURE_FIELDS.CONTROL_PLANE_MACHINE_TYPE]: nodeType,
                [AZURE_FIELDS.WORKER_MACHINE_TYPE]: nodeType,
                [AZURE_FIELDS.CONTROL_PLANE_FLAVOR]: nodeProfile !== AZURE_NODE_PROFILE_NAMES.SINGLE_NODE.valueOf() ? 'prod' : 'dev',
            },
        });
    };

    const onClusterNameChange = (clusterName: string) => {
        onFieldChange(AZURE_FIELDS.CLUSTER_NAME, clusterName);
    };

    const onInstanceTypeChange = (instanceType: string) => {
        onFieldChange(AZURE_FIELDS.NODE_PROFILE, instanceType);
    };

    return (
        <FormProvider {...methods}>
            <div className="wizard-content-container">
                <h3 cds-layout="m-t:md m-b:xl" cds-text="title">
                    Azure Management Cluster Settings
                </h3>
                <div cds-layout="col:12">
                    <PageNotification notification={notification} closeCallback={dismissAlert}></PageNotification>
                </div>
                <div cds-layout="horizontal gap:md align:fill" key="section-holder">
                    <div cds-layout="vertical gap:xxl p-b:lg" key="cluster-name-section">
                        <ClusterName
                            field={AZURE_FIELDS.CLUSTER_NAME}
                            clusterNameChange={onClusterNameChange}
                            placeholderClusterName="my-azure-cluster"
                            defaultClusterName={azureState[STORE_SECTION_FORM][AZURE_FIELDS.CLUSTER_NAME]}
                        />
                        <OsImageSelect
                            osImageLabel={'OS Image with Kubernetes'}
                            images={osImages}
                            field={AZURE_FIELDS.OS_IMAGE}
                            onOsImageSelected={(value) => {
                                onFieldChange(AZURE_FIELDS.OS_IMAGE, value);
                            }}
                            selectedImage={azureState[STORE_SECTION_FORM][AZURE_FIELDS.NODE_PROFILE]}
                        />
                    </div>
                    <div key="instance-type-section">
                        <NodeProfile
                            field={AZURE_FIELDS.NODE_PROFILE}
                            nodeProfileTypes={nodeInstanceTypes}
                            nodeProfileTypeChange={onInstanceTypeChange}
                            selectedProfileId={azureState[STORE_SECTION_FORM][AZURE_FIELDS.NODE_PROFILE]}
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
