// React imports
import React, { ChangeEvent } from 'react';

// Library imports
import { CdsControlMessage } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { CdsButton } from '@cds/react/button';
import { SubmitHandler, useForm } from 'react-hook-form';
<<<<<<< HEAD
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
=======
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea

// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
<<<<<<< HEAD
import { isK8sCompliantString } from '../../../../shared/validations/Validation.service';
=======
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea

export interface FormInputs {
    CLUSTER_NAME: string;
}

<<<<<<< HEAD
const unmanagedClusterBasicSettingStepFormSchema = yup
    .object({
        CLUSTER_NAME: yup
            .string()
            .nullable()
            .required('Please enter a name for your unmanaged cluster')
            .test(
                '',
                'Cluster name must contain only lower case letters and hyphen',
                (value) => value !== null && isK8sCompliantString(value)
            ),
    })
    .required();

=======
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea
function UnmanagedClusterSettings(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;

    const {
        register,
        handleSubmit,
        formState: { errors },
<<<<<<< HEAD
    } = useForm<FormInputs>({ resolver: yupResolver(unmanagedClusterBasicSettingStepFormSchema) });
=======
    } = useForm<FormInputs>();
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea

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
            <div cds-layout="p-b:lg" cds-text="title">
                Cluster settings
            </div>
            <div cds-layout="grid gap:md">
<<<<<<< HEAD
                <div cds-layout="col@sm:4 p-b:md">{ClusterName()}</div>
            </div>
            <div cds-layout="horizontal gap:md">
                <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
            </div>
        </div>
    );

    function ClusterName() {
        return (
            <div>
                <CdsInput>
                    <label cds-layout="p-b:xs" cds-text="section">
                        Cluster name
                    </label>
                    <input
                        {...register('CLUSTER_NAME')}
                        placeholder="cluster-name"
                        onChange={handleClusterNameChange}
                        defaultValue="test-cluster"
                    ></input>
                    {errors['CLUSTER_NAME'] && <CdsControlMessage status="error">{errors['CLUSTER_NAME'].message}</CdsControlMessage>}
                </CdsInput>
                <div>
                    <p className="description" cds-layout="m-t:sm">
                        Can only contain lowercase alphanumeric characters and dashes.
                    </p>
                    <p className="description" cds-layout="m-t:sm">
=======
                <div cds-layout="col@sm:4 p-b:md">
                    <CdsInput>
                        <label cds-layout="p-b:xs" cds-text="section">
                            Cluster name
                        </label>
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
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea
                        The name will be used to reference your cluster in the Tanzu CLI and kubectl.
                    </p>
                </div>
            </div>
<<<<<<< HEAD
        );
    }
=======
            <div cds-layout="horizontal gap:md">
                <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
            </div>
        </div>
    );
>>>>>>> 1b7fc0511e6f37a8425f3d8f43a47d85cd0536ea
}

export default UnmanagedClusterSettings;
