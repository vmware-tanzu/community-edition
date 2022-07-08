// React imports
import React, { ChangeEvent, useContext, useState } from 'react';
import { FieldError, FieldErrors, RegisterOptions, SubmitHandler, useForm, UseFormRegisterReturn } from 'react-hook-form';
// Library imports
import { blockIcon, blocksGroupIcon, ClarityIcons } from '@cds/core/icon';
import { CdsButton } from '@cds/react/button';
import { CdsSelect } from '@cds/react/select';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
// App imports
import { addClusterNameValidation, ClusterNameSection } from '../../../../shared/components/FormInputSections/ClusterNameSection';
import {
    addNodeInstanceTypeValidation,
    NodeInstanceType,
    NodeProfileSection,
} from '../../../../shared/components/FormInputSections/NodeProfileSection';
import { getResource } from '../../../providers/vsphere/VsphereResources.reducer';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { VsphereStore } from '../Store.vsphere.mc';
import { VSphereVirtualMachine } from '../../../../swagger-api';

// NOTE: icons must be imported
const nodeInstanceTypes: NodeInstanceType[] = [
    {
        id: 'single-node',
        label: 'Single node',
        icon: 'block',
        description: 'Create a single control plane node with a medium instance type',
    },
    {
        id: 'high-availability',
        label: 'High availability',
        icon: 'blocks-group',
        description: 'Create a multi-node control plane with a medium instance type',
    },
    {
        id: 'compute-optimized',
        label: 'Production-ready (High availability)',
        icon: 'blocks-group',
        isSolidIcon: true,
        description: 'Create a multi-node control plane with a large instance type',
    },
];
ClarityIcons.addIcons(blockIcon, blocksGroupIcon);

interface VsphereClusterSettingFormInputs {
    [VSPHERE_FIELDS.CLUSTERNAME]: string;
    [VSPHERE_FIELDS.INSTANCETYPE]: string;
}

export function VsphereClusterSettingsStep(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm, handleValueChange } = props;
    const { vsphereState } = useContext(VsphereStore);

    const [selectedInstanceTypeId, setSelectedInstanceTypeId] = useState(
        vsphereState[VSPHERE_FIELDS.INSTANCETYPE] || nodeInstanceTypes[0].id
    );
    const osImages = (getResource('osImages', vsphereState) || []) as VSphereVirtualMachine[];
    const osTemplates = osImages.filter((osImage) => osImage.isTemplate);

    let yupSchemaObject = addNodeInstanceTypeValidation(VSPHERE_FIELDS.INSTANCETYPE, {});
    yupSchemaObject = addClusterNameValidation(VSPHERE_FIELDS.CLUSTERNAME, yupSchemaObject);
    const vsphereClusterSettingsFormSchema = yup.object(yupSchemaObject).required();
    const methods = useForm<VsphereClusterSettingFormInputs>({
        resolver: yupResolver(vsphereClusterSettingsFormSchema),
    });

    const {
        handleSubmit,
        formState: { errors },
        register,
        setValue,
    } = methods;

    const canContinue = (): boolean => {
        return Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<VsphereClusterSettingFormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const onClusterNameChange = (clusterName: string) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, VSPHERE_FIELDS.CLUSTERNAME, clusterName, currentStep, errors);
            setValue(VSPHERE_FIELDS.CLUSTERNAME, clusterName, { shouldValidate: true });
        }
    };

    const onInstanceTypeChange = (instanceType: string) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, VSPHERE_FIELDS.INSTANCETYPE, instanceType, currentStep, errors);
            setValue(VSPHERE_FIELDS.INSTANCETYPE, instanceType, { shouldValidate: true });
        }
        setSelectedInstanceTypeId(instanceType);
    };

    const onOsImageChange = (osImage: string) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, VSPHERE_FIELDS.OSIMAGE, osImage, currentStep, errors);
        }
    };

    return (
        <div>
            <div className="wizard-content-container">
                <h2 cds-layout="m-t:lg">vSphere Management Cluster Settings</h2>
                <div cds-layout="grid gap:xxl" key="section-holder">
                    <div cds-layout="col:6" key="cluster-name-section">
                        {ClusterNameSection(VSPHERE_FIELDS.CLUSTERNAME, errors, register, onClusterNameChange, 'my-vsphere-cluster')}
                    </div>
                    <div cds-layout="col:6" key="instance-type-section">
                        {NodeProfileSection(
                            VSPHERE_FIELDS.INSTANCETYPE,
                            nodeInstanceTypes,
                            errors,
                            register,
                            onInstanceTypeChange,
                            selectedInstanceTypeId
                        )}
                    </div>
                </div>
                <div cds-layout="col:12">{OsImageSection(VSPHERE_FIELDS.OSIMAGE, osTemplates, errors, register, onOsImageChange)}</div>
            </div>
            <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                NEXT
            </CdsButton>
        </div>
    );
}

function OsImageSection(
    field: string,
    osImages: VSphereVirtualMachine[],
    errors: { [key: string]: FieldError | undefined },
    register: any,
    onOsImageSelected: (osImage: string) => void
) {
    const handleOsImageSelect = (event: ChangeEvent<HTMLSelectElement>) => {
        onOsImageSelected(event.target.value || '');
    };
    return (
        <div>
            {' '}
            <CdsSelect layout="vertical" controlWidth="shrink">
                <label cds-layout="p-b:xs">VM Template</label>
                <select {...register(VSPHERE_FIELDS.OSIMAGE)} onChange={handleOsImageSelect}>
                    <option />
                    {osImages.map((dc) => (
                        <option key={dc.moid}>{dc.name}</option>
                    ))}
                </select>
            </CdsSelect>
        </div>
    );
}
