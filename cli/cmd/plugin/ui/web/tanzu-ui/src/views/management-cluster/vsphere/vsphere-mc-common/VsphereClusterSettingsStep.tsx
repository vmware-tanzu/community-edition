// React imports
import React from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
// Library imports
import { CdsButton } from '@cds/react/button';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
// App imports
import { FormInputs } from '../../aws/wizard-basic/management-credential-step/ManagementCredentials';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { vsphereClusterSettingsFormSchema } from './vsphere.clusterSettings.form.schema';

export function VsphereClusterSettingsStep(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const methods = useForm<FormInputs>({
        resolver: yupResolver(vsphereClusterSettingsFormSchema),
    });

    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;

    const canContinue = (): boolean => {
        return Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    return (
        <div>
            <div className="wizard-content-container">
                <h2 cds-layout="m-t:lg">vSphere Cluster Settings</h2>
            </div>
            <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                NEXT
            </CdsButton>
        </div>
    );
}
