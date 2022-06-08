// React imports
import React, { useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
// App imports
import { FormInputs } from '../../aws/wizard-basic/management-credential-step/ManagementCredentials';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { vsphereCredentialFormSchema } from './vsphere.credential.form.schema';

export function VsphereCredentialsStep(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const [connected, setConnection] = useState(true); // TODO: initial state should be false
    const methods = useForm<FormInputs>({
        resolver: yupResolver(vsphereCredentialFormSchema),
    });

    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;

    const canContinue = (): boolean => {
        return connected && Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    // TODO: return true if reasonable data is entered in the form fields
    const connectionDataEntered = (): boolean => {
        return true;
    };

    const handleConnect = () => {
        // TODO: login to vSphere and call setConnection(true) on success
    };

    // TODO: add form fields (and disconnect on changes to them)
    return (
        <div>
            <div className="wizard-content-container">
                <h2 cds-layout="m-t:lg">vSphere Credentials</h2>
            </div>
            <div cds-layout="p-t:lg">
                <CdsButton onClick={handleConnect} disabled={connected || !connectionDataEntered()}>
                    <CdsIcon shape="connect" size="md"></CdsIcon>
                    {connected ? 'CONNECTED' : 'CONNECT'}
                </CdsButton>
            </div>
            <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                NEXT
            </CdsButton>
        </div>
    );
}
