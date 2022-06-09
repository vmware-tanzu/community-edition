// React imports
import React, { ChangeEvent } from 'react';

// Library imports
import { CdsControlMessage } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { CdsButton } from '@cds/react/button';
import { SubmitHandler, useForm } from 'react-hook-form';

// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';

export interface FormInputs {
    CLUSTER_NAME: string;
}

function UnmanagedClusterSettings(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<FormInputs>();

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const handleClusterNameChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'CLUSTER_NAME', event.target.value, currentStep, errors);
        }
    };

    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <h2>Cluster settings</h2>
            <div cds-layout="grid gap:md">
                <div cds-layout="col@sm:4 p-b:md">
                    <CdsInput>
                        <label cds-layout="p-b:md">Cluster name</label>
                        <input
                            {...register('CLUSTER_NAME')}
                            placeholder="cluster-name"
                            onChange={handleClusterNameChange}
                            defaultValue="test-cluster"
                        ></input>
                        {errors['CLUSTER_NAME'] && <CdsControlMessage status="error">{errors['CLUSTER_NAME'].message}</CdsControlMessage>}
                    </CdsInput>
                    <p className="description" cds-layout="m-t:sm">
                        Can only contain lowercase alphanumeric characters and dashes.
                        <br></br>
                        <br></br>
                        The name will be used to reference your cluster in the Tanzu CLI and kubectl.
                    </p>
                </div>
            </div>
            <div cds-layout="horizontal gap:md">
                <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
            </div>
        </div>
    );
}

export default UnmanagedClusterSettings;
