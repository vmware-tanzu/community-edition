// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';
import { FieldError, FieldErrors, RegisterOptions, SubmitHandler, useForm, UseFormRegisterReturn } from 'react-hook-form';
// Library imports
import { blockIcon, blocksGroupIcon, ClarityIcons } from '@cds/core/icon';
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsSelect } from '@cds/react/select';
import { CdsTextarea } from '@cds/react/textarea';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
// App imports
import { clusterNameValidation, ClusterNameSection } from '../../../../shared/components/FormInputSections/ClusterNameSection';
import {
    nodeInstanceTypeValidation,
    NodeInstanceType,
    NodeProfileSection,
} from '../../../../shared/components/FormInputSections/NodeProfileSection';
import { FormAction, VsphereResourceAction } from '../../../../shared/types/types';
import { getResource } from '../../../providers/vsphere/VsphereResources.reducer';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
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

type VSPHERE_CLUSTER_SETTING_STEP_FIELDS =
    | VSPHERE_FIELDS.CLUSTERNAME
    | VSPHERE_FIELDS.INSTANCETYPE
    | VSPHERE_FIELDS.SSHKEY
    | VSPHERE_FIELDS.VMTEMPLATE;

interface VsphereClusterSettingFormInputs {
    [VSPHERE_FIELDS.CLUSTERNAME]: string;
    [VSPHERE_FIELDS.INSTANCETYPE]: string;
    [VSPHERE_FIELDS.SSHKEY]: string;
    [VSPHERE_FIELDS.VMTEMPLATE]: string;
}

export function VsphereClusterSettingsStep(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm, handleValueChange } = props;
    const { vsphereState, vsphereDispatch } = useContext(VsphereStore);
    const vsphereClusterSettingsFormSchema = yup.object(createYupSchemaObject()).required();
    const methods = useForm<VsphereClusterSettingFormInputs>({
        resolver: yupResolver(vsphereClusterSettingsFormSchema),
    });

    const {
        handleSubmit,
        formState: { errors },
        register,
        setValue,
    } = methods;

    const osImages = (getResource('osImages', vsphereState) || []) as VSphereVirtualMachine[];
    const osTemplates = osImages.filter((osImage) => osImage.isTemplate);
    // if there's only ONE template, then pretend the user has selected it (unless we've already done that)
    if (osTemplates.length === 1 && vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.VMTEMPLATE] !== osTemplates[0].moid) {
        const moid = osTemplates[0].moid || '';
        setValue(VSPHERE_FIELDS.VMTEMPLATE, moid);
        vsphereDispatch({
            type: INPUT_CHANGE,
            field: VSPHERE_FIELDS.VMTEMPLATE,
            payload: moid,
        } as FormAction);
    }
    let initialSelectedInstanceTypeId = vsphereState[VSPHERE_FIELDS.INSTANCETYPE];
    if (!initialSelectedInstanceTypeId) {
        initialSelectedInstanceTypeId = nodeInstanceTypes[0].id;
        setValue(VSPHERE_FIELDS.INSTANCETYPE, initialSelectedInstanceTypeId);
    }
    const [selectedInstanceTypeId, setSelectedInstanceTypeId] = useState(initialSelectedInstanceTypeId);

    const canContinue = (): boolean => {
        return Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<VsphereClusterSettingFormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const onFieldChange = (data: string, field: VSPHERE_CLUSTER_SETTING_STEP_FIELDS) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, field, data, currentStep, errors);
            setValue(field, data, { shouldValidate: true });
        }
    };

    const onClusterNameChange = (clusterName: string) => {
        onFieldChange(clusterName, VSPHERE_FIELDS.CLUSTERNAME);
    };

    const onInstanceTypeChange = (instanceType: string) => {
        onFieldChange(instanceType, VSPHERE_FIELDS.INSTANCETYPE);
        setSelectedInstanceTypeId(instanceType);
    };

    return (
        <div>
            <div className="wizard-content-container">
                <h2 cds-layout="m-t:lg">vSphere Management Cluster Settings</h2>
                <div cds-layout="grid gap:m" key="section-holder">
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
                    <div cds-layout="col:12">
                        {VmTemplateSection(VSPHERE_FIELDS.VMTEMPLATE, osTemplates, errors, register, onFieldChange)}
                    </div>
                    <div cds-layout="col:12">{SshKeySection(VSPHERE_FIELDS.SSHKEY, errors, register, onFieldChange)}</div>
                </div>
                <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                    NEXT
                </CdsButton>
            </div>
        </div>
    );
}

function VmTemplateDropdownOptions(vmTemplates: VSphereVirtualMachine[]) {
    if (vmTemplates && vmTemplates.length === 1) {
        return (
            <option key={vmTemplates[0].moid} value={vmTemplates[0].moid}>
                {vmTemplates[0].name}
            </option>
        );
    }
    return (
        <>
            <option />
            {vmTemplates.map((dc) => (
                <option key={dc.moid}>{dc.name}</option>
            ))}
        </>
    );
}

function VmTemplateSection(
    field: VSPHERE_CLUSTER_SETTING_STEP_FIELDS,
    vmTemplates: VSphereVirtualMachine[],
    errors: { [key: string]: FieldError | undefined },
    register: any,
    onOsImageSelected: (osImage: string, field: VSPHERE_CLUSTER_SETTING_STEP_FIELDS) => void
) {
    const handleOsImageSelect = (event: ChangeEvent<HTMLSelectElement>) => {
        onOsImageSelected(event.target.value || '', field);
    };
    const fieldError = errors[field];
    return (
        <div cds-layout="m:lg">
            <CdsSelect layout="vertical" controlWidth="shrink">
                <label>OS Image</label>
                <select {...register(VSPHERE_FIELDS.VMTEMPLATE)} onChange={handleOsImageSelect}>
                    {VmTemplateDropdownOptions(vmTemplates)}
                </select>
            </CdsSelect>
            {fieldError && (
                <div>
                    &nbsp;<CdsControlMessage status="error">{fieldError.message}</CdsControlMessage>
                </div>
            )}
        </div>
    );
}

function SshKeySection(
    field: VSPHERE_CLUSTER_SETTING_STEP_FIELDS,
    errors: { [key: string]: FieldError | undefined },
    register: any,
    onSshKeyEntered: (sshKey: string, field: VSPHERE_CLUSTER_SETTING_STEP_FIELDS) => void
) {
    const handleSshKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
        onSshKeyEntered(event.target.value || '', field);
    };
    const fieldError = errors[field];
    return (
        <div cds-layout="m:lg">
            <CdsFormGroup layout="vertical">
                <CdsTextarea layout="vertical">
                    <label>SSH key</label>
                    <textarea {...register(field)} onChange={handleSshKeyChange}></textarea>
                    {fieldError && <CdsControlMessage status="error">{fieldError.message}</CdsControlMessage>}
                </CdsTextarea>
            </CdsFormGroup>
        </div>
    );
}

function createYupSchemaObject() {
    return {
        [VSPHERE_FIELDS.SSHKEY]: yupStringRequired('Please enter an SSH key'),
        [VSPHERE_FIELDS.VMTEMPLATE]: yupStringRequired('Please select an OS image'),
        [VSPHERE_FIELDS.INSTANCETYPE]: nodeInstanceTypeValidation(),
        [VSPHERE_FIELDS.CLUSTERNAME]: clusterNameValidation(),
    };
}

function yupStringRequired(errorMessage: string) {
    return yup.string().nullable().required(errorMessage);
}
