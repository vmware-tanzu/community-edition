// React imports
import React, { useContext, useEffect } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';

// App imports
import { useTabStatus, useWizardForm } from '../shared/services/form.service';
import { SUBMIT_FORM } from '../state-management/actions/Form.actions';
import { Store } from '../state-management/stores/Store';
import { ChildProps } from '../shared/components/wizard/StepNav';

function ClusterSettings(props: ChildProps | any) {
    const { state, dispatch } = useContext(Store);
    const { register, handleFormSubmit, formState: { errors }, getValues, setValue } = useWizardForm();

    const submitForm = (data: {[key: string]: string}) => {
        dispatch({
            type: SUBMIT_FORM,
            payload: data
        });
    };
    const blurHandler = () => {
        if (Object.keys(errors).length === 0) {
            submitForm(getValues());
        }
    };
    const { currentStep, setTabStatus, tabStatus } = props;
    useTabStatus(errors, { currentStep, setTabStatus, tabStatus });

    const clusterNameValue: string = state.data['CLUSTER_NAME'];

    useEffect(() => {
        setValue('CLUSTER_NAME', clusterNameValue);
    }, [clusterNameValue, setValue]);

    return (
        <CdsFormGroup layout="vertical-inline" control-width="shrink">
            <div cds-layout="horizontal gap:lg align:vertical-center">
                <CdsInput>
                    <label>CLUSTER NAME</label>
                    <input placeholder="CLUSTER NAME" 
                        {...register('CLUSTER_NAME', { required: true })}
                        defaultValue={state.data['CLUSTER_NAME']||''} onBlur={() => blurHandler()}/>
                    { errors['CLUSTER_NAME'] && <CdsControlMessage status="error">{errors['CLUSTER_NAME'].message}</CdsControlMessage> }
                </CdsInput>
            </div>
            <CdsButton onClick={() => {
                handleFormSubmit({ ...props, submitForm });
            }}>Next</CdsButton>
        </CdsFormGroup>
    );
}

export default ClusterSettings;
