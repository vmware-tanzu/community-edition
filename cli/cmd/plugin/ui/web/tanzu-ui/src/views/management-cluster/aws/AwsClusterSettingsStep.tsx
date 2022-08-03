// React imports
import React, { useContext, useState, useEffect } from 'react';

// Library imports
import { ClarityIcons, blockIcon, blocksGroupIcon, clusterIcon } from '@cds/core/icon';
import { FormProvider, useForm } from 'react-hook-form';
import { CdsIcon } from '@cds/react/icon';
import { CdsButton } from '@cds/react/button';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
import * as yup from 'yup';

// App imports
import { AwsStore } from '../../../state-management/stores/Store.aws';
import { AwsService, AWSVirtualMachine } from '../../../swagger-api';
import { AWS_FIELDS } from './aws-mc-basic/AwsManagementClusterBasic.constants';
import { ClusterName, clusterNameValidation } from '../../../shared/components/FormInputComponents/ClusterName/ClusterName';
import { FormAction } from '../../../shared/types/types';
import { INPUT_CHANGE } from '../../../state-management/actions/Form.actions';
import { StepProps } from '../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../state-management/reducers/Form.reducer';
import {
    NodeProfile,
    NodeInstanceType,
    nodeInstanceTypeValidation,
} from '../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import OsImageSelect from '../../../shared/components/FormInputComponents/OsImageSelect/OsImageSelect';
import PageNotification, { Notification, NotificationStatus } from '../../../shared/components/PageNotification/PageNotification';
import UseUpdateTabStatus from '../../../shared/components/wizard/UseUpdateTabStatus.hooks';

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

type AWS_CLUSTER_SETTING_STEP_FIELDS = 'NODE_PROFILE' | 'IMAGE_INFO' | 'CLUSTER_NAME';

const nodeInstanceTypes: NodeInstanceType[] = [
    {
        id: 'SINGLE_NODE',
        label: 'Single node',
        icon: 'block',
        description: 'Create a single control plane node with a medium instance type',
    },
    {
        id: 'HIGH_AVAILABILITY',
        label: 'High availability',
        icon: 'blocks-group',
        description: 'Create a multi-node control plane with a medium instance type',
    },
    {
        id: 'PRODUCTION_READY',
        label: 'Production-ready (High availability)',
        icon: 'blocks-group',
        isSolidIcon: true,
        description: 'Create a multi-node control plane with a large instance type',
    },
];

interface AwsClusterSettingFormInputs {
    [AWS_FIELDS.CLUSTER_NAME]: string;
    [AWS_FIELDS.NODE_PROFILE]: string;
    [AWS_FIELDS.IMAGE_INFO]: string;
}

function createYupSchemaObject() {
    return {
        IMAGE_INFO: yupStringRequired('Please select an OS image'),
        NODE_PROFILE: nodeInstanceTypeValidation(),
        CLUSTER_NAME: clusterNameValidation(),
    };
}

function yupStringRequired(errorMessage: string) {
    return yup.string().nullable().required(errorMessage);
}

function AwsClusterSettingsStep(props: Partial<StepProps>) {
    const { updateTabStatus, currentStep, goToStep } = props;
    const { awsState, awsDispatch } = useContext(AwsStore);
    const awsClusterSettingsFormSchema = yup.object(createYupSchemaObject()).required();
    const methods = useForm<AwsClusterSettingFormInputs>({
        resolver: yupResolver(awsClusterSettingsFormSchema),
        mode: 'all',
    });
    const [notification, setNotification] = useState<Notification | null>(null);
    const [images, setImages] = useState<AWSVirtualMachine[]>([]);
    const {
        handleSubmit,
        formState: { errors },
        setValue,
    } = methods;

    function dismissAlert() {
        setNotification(null);
    }

    // update tab status bar
    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }

    const goToNextStep = () => {
        if (goToStep && currentStep) {
            goToStep(currentStep + 1);
        }
    };

    const region = awsState[STORE_SECTION_FORM].REGION;
    useEffect(() => {
        const setImageInfo = (image: any) => {
            awsDispatch({
                type: INPUT_CHANGE,
                field: AWS_FIELDS.IMAGE_INFO,
                payload: image,
            } as FormAction);
        };
        const fetchImages = async () => {
            try {
                const data = await AwsService.getAwsosImages(region);
                setImages(data);
                setImageInfo(data[0]);
            } catch (e: any) {
                setNotification({
                    status: NotificationStatus.DANGER,
                    message: `Unable to retrieve OS Images: ${e}`,
                } as Notification);
            }
        };
        fetchImages();
    }, [awsDispatch, region]);

    let initialSelectedInstanceTypeId = awsState[STORE_SECTION_FORM].NODE_PROFILE;

    if (!initialSelectedInstanceTypeId) {
        initialSelectedInstanceTypeId = nodeInstanceTypes[0].id;
        setValue(AWS_FIELDS.NODE_PROFILE, initialSelectedInstanceTypeId);
    }
    const [selectedInstanceTypeId, setSelectedInstanceTypeId] = useState(initialSelectedInstanceTypeId);

    const onFieldChange = (field: AWS_CLUSTER_SETTING_STEP_FIELDS, data: any) => {
        awsDispatch({
            type: INPUT_CHANGE,
            field,
            payload: data,
        } as FormAction);
    };

    const onInstanceTypeChange = (instanceType: string) => {
        setSelectedInstanceTypeId(instanceType);
        onFieldChange(AWS_FIELDS.NODE_PROFILE, instanceType);
    };

    return (
        <FormProvider {...methods}>
            <div className="cluster-settings-container" cds-layout="m:lg">
                <h3>Management Cluster settings</h3>
                <div cds-layout="grid gap:m" key="section-holder">
                    <div cds-layout="col:4" key="cluster-name-section">
                        <ClusterName
                            field={AWS_FIELDS.CLUSTER_NAME}
                            clusterNameChange={(value) => {
                                onFieldChange(AWS_FIELDS.CLUSTER_NAME, value);
                            }}
                            placeholderClusterName={'my-aws-cluster'}
                            defaultClusterName={awsState[STORE_SECTION_FORM][AWS_FIELDS.CLUSTER_NAME]}
                        />
                    </div>
                    <div cds-layout="col:8" key="instance-type-section">
                        <NodeProfile
                            field={AWS_FIELDS.NODE_PROFILE}
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
                            osImageTitle="Amazon Machine Image(AMI)"
                            images={images}
                            field={AWS_FIELDS.IMAGE_INFO}
                            onOsImageSelected={(value) => {
                                onFieldChange(AWS_FIELDS.IMAGE_INFO, value);
                            }}
                        />
                    </div>

                    <div cds-layout="grid col:12 p-t:lg">
                        <CdsButton onClick={handleSubmit(goToNextStep)}>NEXT</CdsButton>
                    </div>
                </div>
            </div>
        </FormProvider>
    );
}

export default AwsClusterSettingsStep;
