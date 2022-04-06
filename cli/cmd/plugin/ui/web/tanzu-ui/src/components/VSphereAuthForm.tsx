// React imports
import React, { useContext } from 'react';

// Library imports
import { CdsInput } from '@cds/react/input';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsButton } from '@cds/react/button';
import { yupResolver } from '@hookform/resolvers/yup';

// App imports
import { SUBMIT_FORM } from '../state-management/actions/Form.actions';
import { Store } from '../state-management/stores/Store';
import { authFormSchema } from './vsphere.auth.form.schema';
import { blurHandler, onBlurWithDependencies, useTabStatus, useWizardForm } from '../shared/services/form.service';
import { dependencyMap } from './DependencyMap';
import { ChildProps } from '../shared/components/wizard/StepNav';

function VSphereAuthForm (props: ChildProps | any) {
    const { state, dispatch } = useContext(Store);
    const { register, handleFormSubmit, formState: { errors }, getValues } = useWizardForm({
        resolver: yupResolver(authFormSchema)
    });
    const submitForm = (data: {[key: string]: string}) => {
        dispatch({
            type: SUBMIT_FORM,
            payload: data
        });
    };

    const { currentStep, setTabStatus, tabStatus } = props;
    useTabStatus(errors, { currentStep, setTabStatus, tabStatus });
    return (
        <CdsFormGroup layout="vertical-inline" control-width="shrink">
            <div cds-layout="horizontal gap:lg align:vertical-center">
                <CdsInput>
                    <label>VCENTER SERVER</label>
                    <input placeholder="IP OR FQDN" 
                        {...register('VCENTER_SERVER')}
                        defaultValue={state.data['VCENTER_SERVER']||''} onBlur={() => blurHandler(errors, submitForm, getValues())}/>
                    { errors['VCENTER_SERVER'] && <CdsControlMessage status="error">{errors['VCENTER_SERVER'].message}</CdsControlMessage> }
                </CdsInput>
                
                <CdsInput>
                    <label>USERNAME</label>
                    <input placeholder="Username"
                        {...register('VCENTER_USERNAME')}
                        defaultValue={state.data['VCENTER_USERNAME']||''} onBlur={() => blurHandler(errors, submitForm, getValues())}/>
                    { errors.VCENTER_USERNAME && <CdsControlMessage status="error">{errors.VCENTER_USERNAME.message}</CdsControlMessage> }
                </CdsInput>
                <CdsInput>
                    <label>PASSWORD</label>
                    <input placeholder="Password"
                        {...register('VCENTER_PASSWORD')}
                        type="password"
                        defaultValue={state.data['VCENTER_PASSWORD']||''} onBlur={(e) => {
                            blurHandler(errors, submitForm, getValues());
                            onBlurWithDependencies({
                                name: 'VCENTER_PASSWORD',
                                dependencyMap,
                                dispatch
                            });
                        }}/>
                    { errors.VCENTER_PASSWORD && <CdsControlMessage status="error">{errors.VCENTER_PASSWORD.message}</CdsControlMessage> }
                </CdsInput>
            </div>
            <CdsButton onClick={() => {
                handleFormSubmit({ ...props, submitForm });
            }}>Next</CdsButton>
        </CdsFormGroup>
    );
}

export default VSphereAuthForm;