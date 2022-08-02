// React imports
import React, { ChangeEvent, useContext } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { isK8sCompliantString } from '../../../../shared/validations/Validation.service';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { UNMANAGED_CLUSTER_FIELDS } from '../../unmanaged-cluster-common/UnmanagedCluster.constants';
import { UNMANAGED_PLACEHOLDER_VALUES } from '../../unmanaged-cluster-common/unmanaged.defaults';
import { UmcStore } from '../../../../state-management/stores/Store.umc';

export interface FormInputs {
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]: string;
}

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

function UnmanagedClusterSettings(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;

    const { umcState } = useContext(UmcStore);

    const methods = useForm<FormInputs>({
        resolver: yupResolver(unmanagedClusterBasicSettingStepFormSchema),
    });

    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (Object.keys(errors).length === 0 && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const handleClusterNameChange = (event: ChangeEvent<HTMLInputElement>) => {
        setValue(UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME, event.target.value, { shouldValidate: true });
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME, event.target.value, currentStep, errors);
        }
    };

    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <div cds-layout="p-b:lg" cds-text="title">
                Cluster settings
            </div>
            <div cds-layout="grid">
                <div cds-layout="col@sm:8">
                    <div cds-layout="vertical gap:lg">
                        <div cds-layout="grid gap:md">
                            <div cds-layout="col@sm:6">{ClusterName()}</div>
                        </div>
                        <div cds-layout="horizontal gap:md">
                            <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );

    function ClusterName() {
        const errorClusterName = errors[UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME];
        return (
            <div>
                <CdsInput layout="vertical">
                    <label cds-layout="p-b:xs" cds-text="section">
                        Cluster name
                    </label>
                    <input
                        {...register(UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME)}
                        placeholder={UNMANAGED_PLACEHOLDER_VALUES[UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]}
                        onChange={handleClusterNameChange}
                        defaultValue={umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]}
                    ></input>
                    {errorClusterName && <CdsControlMessage status="error">{errorClusterName.message}</CdsControlMessage>}
                </CdsInput>
                <div>
                    <p className="description" cds-layout="m-t:sm">
                        Can only contain lowercase alphanumeric characters and dashes.
                    </p>
                    <p className="description" cds-layout="m-t:sm">
                        The name will be used to reference your cluster in the Tanzu CLI and kubectl.
                    </p>
                </div>
            </div>
        );
    }
}

export default UnmanagedClusterSettings;
