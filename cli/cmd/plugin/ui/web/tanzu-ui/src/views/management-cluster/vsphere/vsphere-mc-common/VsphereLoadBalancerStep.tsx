// React imports
import React, { ChangeEvent, useContext, useState } from 'react';
import { FieldError, FieldErrors, RegisterOptions, SubmitHandler, useForm, UseFormRegisterReturn } from 'react-hook-form';
// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsSelect } from '@cds/react/select';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { ENDPOINT_PROVIDER_IDS, ENDPOINT_PROVIDERS, IPFAMILIES, VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { VsphereStore } from '../Store.vsphere.mc';
import { CdsInput } from '@cds/react/input';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { yupServerTest } from './vsphere.credential.form.schema';

interface VsphereLoadBalancerFormInputs {
    [VSPHERE_FIELDS.CLUSTER_ENDPOINT]: string;
}

export function VsphereLoadBalancerStep(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm, handleValueChange } = props;
    const { vsphereState } = useContext(VsphereStore);
    const [endpoint, setEndpoint] = useState<string>('');
    const ipFamily = vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.IPFAMILY] || IPFAMILIES.IPv4;
    const vsphereLoadBalancerStepFormSchema = yup.object({
        [VSPHERE_FIELDS.CLUSTER_ENDPOINT]: yupServerTest(ipFamily).required('vSphere server name is required'),
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
                {EndpointProviderSection(endpoint, ipFamily, errors, register, onEndpointChange)}
                <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                    NEXT
                </CdsButton>
            </div>
        </div>
    );
}

function EndpointProviderSection(
    endpoint: string,
    ipFamily: string,
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
                        <option key={ENDPOINT_PROVIDER_IDS.KUBE_VIP} value={ENDPOINT_PROVIDER_IDS.KUBE_VIP}>
                            {ENDPOINT_PROVIDERS[ENDPOINT_PROVIDER_IDS.KUBE_VIP]}
                        </option>
                    </select>
                    <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>
                </CdsSelect>

                <CdsInput layout="compact">
                    <label cds-layout="p-b:xs">Control Plane Endpoint ({ipFamily})</label>
                    <input {...register(VSPHERE_FIELDS.CLUSTER_ENDPOINT)} onChange={handleEndpointChange} defaultValue={endpoint} />
                    {err && <CdsControlMessage status="error">{err.message}</CdsControlMessage>}
                    {!err && <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>}
                </CdsInput>
            </div>
        </CdsFormGroup>
    );
}
