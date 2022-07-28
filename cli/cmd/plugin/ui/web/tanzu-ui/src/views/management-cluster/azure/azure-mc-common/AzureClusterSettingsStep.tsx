// React imports
import React, { useContext, useState, useEffect } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
import * as yup from 'yup';

// App imports

import { AZURE_FIELDS, AZURE_NODE_PROFILE_NAMES } from '../azure-mc-basic/AzureManagementClusterBasic.constants';
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
import { AzureService, AzureVirtualMachine } from '../../../../swagger-api';
import RetrieveOSImages from '../../../../shared/components/FormInputComponents/RetrieveOSImages/RetrieveOSImages';

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

type AZURE_CLUSTER_SETTING_STEP_FIELDS = AZURE_FIELDS.CLUSTER_NAME | AZURE_FIELDS.NODE_PROFILE;

interface AzureClusterSettingFormInputs {
    [AZURE_FIELDS.CLUSTER_NAME]: string;
    [AZURE_FIELDS.NODE_PROFILE]: string;
    [AZURE_FIELDS.IMAGE_INFO]: string;
}

export function AzureClusterSettingsStep(props: Partial<StepProps>) {
    const { currentStep, deploy, handleValueChange } = props;
    const { azureState, azureDispatch } = useContext(AzureStore);
    const azureClusterSettingsFormSchema = yup.object(createYupSchemaObject()).required();
    const methods = useForm<AzureClusterSettingFormInputs>({
        resolver: yupResolver(azureClusterSettingsFormSchema),
    });
    const [images, setImages] = useState<AzureVirtualMachine[]>([]);
    const {
        handleSubmit,
        formState: { errors },
        register,
        setValue,
    } = methods;

    let initialSelectedNodeProfileId = azureState[AZURE_FIELDS.NODE_PROFILE];
    if (!initialSelectedNodeProfileId) {
        initialSelectedNodeProfileId = nodeInstanceTypes[0].id;
        setValue(AZURE_FIELDS.NODE_PROFILE, initialSelectedNodeProfileId);
    }
    const [selectedInstanceTypeId, setSelectedInstanceTypeId] = useState(initialSelectedNodeProfileId);

    const setImageParameters = (image) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, AZURE_FIELDS.IMAGE_INFO, image, currentStep, errors);
            setValue(AZURE_FIELDS.IMAGE_INFO, image.name);
        }
    };

    const canContinue = (): boolean => {
        return (
            Object.keys(errors).length === 0 &&
            azureState.dataForm[AZURE_FIELDS.CLUSTER_NAME] &&
            azureState.dataForm[AZURE_FIELDS.NODE_PROFILE]
        );
    };

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
        onFieldChange(instanceType, AZURE_FIELDS.NODE_PROFILE);
        setSelectedInstanceTypeId(instanceType);
    };

    const onOsImageSelected = (imageName: string) => {
        images.some((image) => {
            if (image.name === imageName) {
                setImageParameters(image);
            }
        });
    };

    useEffect(() => {
        AzureService.getAzureOsImages().then((data) => {
            setImages(data);
            setImageParameters(data[0]);
        });
    }, []);

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
                            field={AZURE_FIELDS.NODE_PROFILE}
                            nodeInstanceTypes={nodeInstanceTypes}
                            errors={errors}
                            register={register}
                            nodeInstanceTypeChange={onInstanceTypeChange}
                            selectedInstanceId={selectedInstanceTypeId}
                        />
                    </div>
                    <div cds-layout="col:12">
                        <RetrieveOSImages
                            osImageTitle={'Amazon Machine Image(AMI)'}
                            images={images}
                            field={'IMAGE_INFO'}
                            errors={errors}
                            register={register}
                            onOsImageSelected={onOsImageSelected}
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
        [AZURE_FIELDS.NODE_PROFILE]: nodeInstanceTypeValidation(),
        [AZURE_FIELDS.CLUSTER_NAME]: clusterNameValidation(),
    };
}
