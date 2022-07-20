// React imports
import React, { ChangeEvent, useContext, useState } from 'react';
import { FieldError, FieldErrors, RegisterOptions, SubmitHandler, useForm, UseFormRegisterReturn } from 'react-hook-form';
// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { CdsSelect } from '@cds/react/select';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import {
    ENDPOINT_PROVIDERS,
    ENDPOINT_PROVIDERS_DISPLAY,
    IP_FAMILIES_DISPLAY,
    IP_FAMILIES,
    VSPHERE_FIELDS,
} from '../VsphereManagementCluster.constants';
import { VsphereStore } from '../Store.vsphere.mc';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { yupServerTest } from './vsphere.credential.form.schema';

interface VsphereLoadBalancerFormInputs {
    [VSPHERE_FIELDS.CLUSTER_ENDPOINT]: string;
}

export function VsphereLoadBalancerStep(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm, handleValueChange } = props;
    const { vsphereState } = useContext(VsphereStore);
    const [endpoint, setEndpoint] = useState<string>('');
    const ipFamilyId = vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.IPFAMILY] || IP_FAMILIES.IPv4;
    const vsphereLoadBalancerStepFormSchema = yup.object({
        [VSPHERE_FIELDS.CLUSTER_ENDPOINT]: yupServerTest(ipFamilyId).required('vSphere server name is required'),
    });
    const methods = useForm<VsphereLoadBalancerFormInputs>({
        resolver: yupResolver(vsphereLoadBalancerStepFormSchema),
    });

    const {
        handleSubmit,
        formState: { errors },
        register,
        setValue,
    } = methods;

    const canContinue = () => {
        return Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<VsphereLoadBalancerFormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const onEndpointChange = (endpoint: string) => {
        setEndpoint(endpoint);
        setValue && setValue(VSPHERE_FIELDS.CLUSTER_ENDPOINT, endpoint, { shouldValidate: true });
        handleValueChange && handleValueChange(INPUT_CHANGE, VSPHERE_FIELDS.CLUSTER_ENDPOINT, endpoint, currentStep, errors);
    };

    return (
        <div>
            <div className="wizard-content-container">
                <h2 cds-layout="m-t:lg">vSphere Load Balancer Settings</h2>
                {EndpointProviderSection(endpoint, ipFamilyId, errors, register, onEndpointChange)}
                <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                    NEXT
                </CdsButton>
            </div>
        </div>
    );
}

function EndpointProviderSection(
    endpoint: string,
    ipFamily: IP_FAMILIES,
    errors: any,
    register: any,
    onEndpointChange: (endpoint: string) => void
) {
    const handleEndpointChange = (event: ChangeEvent<HTMLInputElement>) => {
        onEndpointChange(event.target.value);
    };
    const err = errors[VSPHERE_FIELDS.CLUSTER_ENDPOINT];
    return (
        <CdsFormGroup layout="vertical-inline" control-width="shrink">
            <div cds-layout="horizontal gap:lg align:vertical-center p-b:sm">
                <CdsSelect layout="vertical" controlWidth="shrink">
                    <label cds-layout="p-b:xs">Control Plane Endpoint Provider</label>
                    <select>
                        <option key={ENDPOINT_PROVIDERS.KUBE_VIP} value={ENDPOINT_PROVIDERS.KUBE_VIP}>
                            {ENDPOINT_PROVIDERS_DISPLAY[ENDPOINT_PROVIDERS.KUBE_VIP]}
                        </option>
                    </select>
                    <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>
                </CdsSelect>

                <CdsInput layout="compact">
                    <label cds-layout="p-b:xs">Control Plane Endpoint ({IP_FAMILIES_DISPLAY[ipFamily]})</label>
                    <input {...register(VSPHERE_FIELDS.CLUSTER_ENDPOINT)} onChange={handleEndpointChange} defaultValue={endpoint} />
                    {err && <CdsControlMessage status="error">{err.message}</CdsControlMessage>}
                    {!err && <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>}
                </CdsInput>
            </div>
        </CdsFormGroup>
    );
}
