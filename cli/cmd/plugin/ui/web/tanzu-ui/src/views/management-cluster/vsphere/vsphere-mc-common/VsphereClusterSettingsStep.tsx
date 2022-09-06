// React imports
import React, { ChangeEvent, MouseEvent, useContext, useEffect, useState } from 'react';
import { FieldError, FormProvider, SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { blockIcon, blocksGroupIcon, ClarityIcons } from '@cds/core/icon';
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsTextarea } from '@cds/react/textarea';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';

// App imports
import { ClusterName, clusterNameValidation } from '../../../../shared/components/FormInputComponents/ClusterName/ClusterName';
import { createNodeProfileFieldUpdateObject, VSPHERE_NODE_PROFILES } from './VsphereNodeProfileUtil';
import { FormAction } from '../../../../shared/types/types';
import { INPUT_CHANGE, BATCH_SET } from '../../../../state-management/actions/Form.actions';
import {
    NodeProfileType,
    NodeProfile,
    nodeProfileValidation,
} from '../../../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { VsphereStore } from '../Store.vsphere.mc';
import { VsphereService, VSphereVirtualMachine } from '../../../../swagger-api';
import UseUpdateTabStatus from '../../../../shared/components/wizard/UseUpdateTabStatus.hooks';
import SpinnerSelect from '../../../../shared/components/Select/SpinnerSelect';

// NOTE: icons must be imported
const nodeProfileTypes: NodeProfileType[] = [
    {
        id: VSPHERE_NODE_PROFILES.SINGLE_NODE,
        label: 'Single node',
        icon: 'block',
        description: 'Create a single control plane node with a medium instance type',
    },
    {
        id: VSPHERE_NODE_PROFILES.HIGH_AVAILABILITY,
        label: 'High availability',
        icon: 'blocks-group',
        description: 'Create a multi-node control plane with a medium instance type',
    },
    {
        id: VSPHERE_NODE_PROFILES.PRODUCTION_READY,
        label: 'Production-ready (High availability)',
        icon: 'blocks-group',
        isSolidIcon: true,
        description: 'Create a multi-node control plane with a large instance type',
    },
];
ClarityIcons.addIcons(blockIcon, blocksGroupIcon);

type VSPHERE_CLUSTER_SETTING_STEP_FIELDS =
    | VSPHERE_FIELDS.CLUSTERNAME
    | VSPHERE_FIELDS.NODE_PROFILE_TYPE
    | VSPHERE_FIELDS.SSHKEY
    | VSPHERE_FIELDS.VMTEMPLATE;

interface VsphereClusterSettingFormInputs {
    [VSPHERE_FIELDS.CLUSTERNAME]: string;
    [VSPHERE_FIELDS.NODE_PROFILE_TYPE]: string;
    [VSPHERE_FIELDS.SSHKEY]: string;
    [VSPHERE_FIELDS.VMTEMPLATE]: string;
}

export function VsphereClusterSettingsStep(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm, updateTabStatus } = props;
    const { vsphereState, vsphereDispatch } = useContext(VsphereStore);
    const vsphereClusterSettingsFormSchema = yup.object(createYupSchemaObject()).required();
    const [osImagesLoading, setOsImagesLoading] = useState(false);
    const [osImages, setOsImages] = useState<VSphereVirtualMachine[]>([]);
    const methods = useForm<VsphereClusterSettingFormInputs>({
        resolver: yupResolver(vsphereClusterSettingsFormSchema),
        mode: 'all',
    });

    const {
        handleSubmit,
        formState: { errors },
        register,
        setValue,
    } = methods;

    // update tab status bar
    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }

    useEffect(() => {
        fetchOsImages();
    }, [vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.DATACENTER]]);

    const osTemplates = osImages.filter((osImage) => osImage.isTemplate);
    // if there's only ONE template, then pretend the user has selected it (unless we've already done that)
    useEffect(() => {
        if (osTemplates.length === 1 && vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.VMTEMPLATE] !== osTemplates[0].moid) {
            const moid = osTemplates[0].moid || '';
            setValue(VSPHERE_FIELDS.VMTEMPLATE, moid);
            vsphereDispatch({
                type: INPUT_CHANGE,
                field: VSPHERE_FIELDS.VMTEMPLATE,
                payload: osTemplates[0],
            } as FormAction);
        }
    }, [osImages]);

    const recordFieldsRelatedToNodeProfile = (profileId: VSPHERE_NODE_PROFILES) => {
        const fieldValues = createNodeProfileFieldUpdateObject(profileId);
        vsphereDispatch({ type: BATCH_SET, payload: fieldValues } as FormAction);
    };

    const initialSelectedNodeProfileId = vsphereState[VSPHERE_FIELDS.NODE_PROFILE_TYPE];
    useEffect(() => {
        let initialProfile = initialSelectedNodeProfileId;
        if (!initialProfile) {
            initialProfile = nodeProfileTypes[0].id;
            vsphereDispatch({
                type: INPUT_CHANGE,
                field: VSPHERE_FIELDS.NODE_PROFILE_TYPE,
                payload: initialProfile,
            } as FormAction);
        }
        setValue(VSPHERE_FIELDS.NODE_PROFILE_TYPE, initialProfile);
        recordFieldsRelatedToNodeProfile(initialProfile);
    }, []);

    const canContinue = (): boolean => {
        return Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<VsphereClusterSettingFormInputs> = () => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const onFieldChange = (data: string, field: VSPHERE_CLUSTER_SETTING_STEP_FIELDS) => {
        vsphereDispatch({
            type: INPUT_CHANGE,
            field,
            payload: data,
        } as FormAction);
    };

    const onClusterNameChange = (clusterName: string) => {
        onFieldChange(clusterName, VSPHERE_FIELDS.CLUSTERNAME);
    };

    const onNodeProfileChange = (profileId: string) => {
        onFieldChange(profileId, VSPHERE_FIELDS.NODE_PROFILE_TYPE);
        recordFieldsRelatedToNodeProfile(profileId as VSPHERE_NODE_PROFILES);
    };

    const handleRefresh = async (event: MouseEvent<HTMLAnchorElement>) => {
        event.preventDefault();
        fetchOsImages();
    };

    const fetchOsImages = async () => {
        try {
            setOsImagesLoading(true);
            const datacenterMoid = vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.DATACENTER]?.moid;
            if (datacenterMoid) {
                await VsphereService.getVSphereOsImages(datacenterMoid).then((fetchedOsImages) => {
                    setOsImages(fetchedOsImages);
                });
            }
        } catch (e: any) {
            console.log(`Unable to retrieve OS Images: ${JSON.stringify(e)}`);
        } finally {
            setOsImagesLoading(false);
        }
    };

    return (
        <FormProvider {...methods}>
            <div className="wizard-content-container">
                <h3 cds-layout="m-t:md m-b:xl" cds-text="title">
                    vSphere Management Cluster Settings
                </h3>
                <div cds-layout="horizontal gap:md align:fill" key="section-holder">
                    <div cds-layout="vertical gap:xxl p-b:lg" key="cluster-name-section">
                        <ClusterName
                            field={VSPHERE_FIELDS.CLUSTERNAME}
                            clusterNameChange={onClusterNameChange}
                            placeholderClusterName={'my-vsphere-cluster'}
                        />
                        {VmTemplateSection(
                            VSPHERE_FIELDS.VMTEMPLATE,
                            osTemplates,
                            errors,
                            register,
                            osImagesLoading,
                            handleRefresh,
                            onFieldChange
                        )}
                    </div>
                    <div key="profile-type-section">
                        <NodeProfile
                            field={VSPHERE_FIELDS.NODE_PROFILE_TYPE}
                            nodeProfileTypes={nodeProfileTypes}
                            nodeProfileTypeChange={onNodeProfileChange}
                            selectedProfileId={vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.NODE_PROFILE_TYPE]}
                        />
                    </div>
                </div>
                <div cds-layout="col:12">{SshKeySection(VSPHERE_FIELDS.SSHKEY, errors, register, onFieldChange)}</div>
                <CdsButton cds-layout="m-t:lg" onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                    NEXT
                </CdsButton>
            </div>
        </FormProvider>
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
    osImagesLoading: boolean,
    handleRefresh: (event: MouseEvent<HTMLAnchorElement>) => void,
    onOsImageSelected: (osImage: string, field: VSPHERE_CLUSTER_SETTING_STEP_FIELDS) => void
) {
    const handleOsImageSelect = (event: ChangeEvent<HTMLSelectElement>) => {
        onOsImageSelected(event.target.value || '', field);
    };

    const fieldError = errors[field];
    return (
        <CdsFormGroup layout="vertical-inline" control-width="shrink">
            <div cds-layout="horizontal gap:lg align:vertical-center">
                <SpinnerSelect
                    className="select-md-width"
                    disabled={false}
                    label="OS Image"
                    handleSelect={handleOsImageSelect}
                    name={VSPHERE_FIELDS.VMTEMPLATE}
                    controlMessage="OS Images will be retrieved after selecting a Datacenter"
                    isLoading={osImagesLoading}
                    register={register}
                    error=""
                >
                    {VmTemplateDropdownOptions(vmTemplates)}
                </SpinnerSelect>
                <a href="/" className={'btn-refresh icon-blue'} onClick={handleRefresh} cds-text="secondary">
                    <CdsIcon shape="refresh" size="sm"></CdsIcon>{' '}
                    <span cds-layout="m-t:sm" className="vertical-mid">
                        REFRESH
                    </span>
                </a>
                {fieldError && (
                    <div>
                        &nbsp;<CdsControlMessage status="error">{fieldError.message}</CdsControlMessage>
                    </div>
                )}
            </div>
        </CdsFormGroup>
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
        <div cds-layout="m-t:lg">
            <CdsFormGroup layout="vertical">
                <CdsTextarea layout="vertical" className="textarea-ssh-key">
                    <label>SSH public key</label>
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
        [VSPHERE_FIELDS.NODE_PROFILE_TYPE]: nodeProfileValidation(),
        [VSPHERE_FIELDS.CLUSTERNAME]: clusterNameValidation(),
    };
}

function yupStringRequired(errorMessage: string) {
    return yup.string().nullable().required(errorMessage);
}
