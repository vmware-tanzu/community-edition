// React imports
import React, { useContext } from 'react';
import { SubmitHandler, FormProvider, useForm } from 'react-hook-form';

// Library imports
import { CdsButton } from '@cds/react/button';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import { ClusterName, clusterNameValidation } from '../../../../shared/components/FormInputComponents/ClusterName/ClusterName';
import { FormAction } from '../../../../shared/types/types';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { UNMANAGED_CLUSTER_FIELDS } from '../../unmanaged-cluster-common/UnmanagedCluster.constants';
import { UNMANAGED_PLACEHOLDER_VALUES } from '../../unmanaged-cluster-common/unmanaged.defaults';
import { UmcStore } from '../../../../state-management/stores/Store.umc';
import UseUpdateTabStatus from '../../../../shared/components/wizard/UseUpdateTabStatus.hooks';

export interface FormInputs {
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]: string;
}

function createYupSchemaObject() {
    return {
        [UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]: clusterNameValidation(),
    };
}
const unmanagedClusterBasicSettingStepFormSchema = yup.object(createYupSchemaObject()).required();

function UnmanagedClusterSettings(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm, updateTabStatus } = props;
    const methods = useForm<FormInputs>({
        resolver: yupResolver(unmanagedClusterBasicSettingStepFormSchema),
        mode: 'all',
    });
    const {
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;

    const { umcDispatch } = useContext(UmcStore);

    // update tab status bar
    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }

    const canContinue = (): boolean => {
        return Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const onFieldChange = (field: UNMANAGED_CLUSTER_FIELDS, data: string) => {
        umcDispatch({
            type: INPUT_CHANGE,
            field,
            payload: data,
        } as FormAction);
    };

    return (
        <FormProvider {...methods}>
            <div className="cluster-settings-container" cds-layout="m:lg">
                <div cds-layout="p-b:lg" cds-text="title">
                    Cluster settings
                </div>
                <div cds-layout="grid gap:m" key="section-holder">
                    <div cds-layout="col:6" key="cluster-name-section">
                        <ClusterName
                            field={UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME}
                            clusterNameChange={(value) => {
                                onFieldChange(UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME, value);
                            }}
                            placeholderClusterName={UNMANAGED_PLACEHOLDER_VALUES.CLUSTER_NAME}
                        />
                    </div>
                </div>
                <div cds-layout="p-t:lg">
                    <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                        NEXT
                    </CdsButton>
                </div>
            </div>
        </FormProvider>
    );
}

export default UnmanagedClusterSettings;
