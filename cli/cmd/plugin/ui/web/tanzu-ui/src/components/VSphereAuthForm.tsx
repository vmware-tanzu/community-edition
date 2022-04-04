// React imports
import React, { useContext } from 'react';

// Library imports
import { CdsInput } from '@cds/react/input';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsButton } from '@cds/react/button';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';

// App imports
import { SUBMIT_FORM } from '../state-management/actions/Form.actions';
import { Store } from '../state-management/stores/Store';
import { authFormSchema } from './vsphere.auth.form.schema';

function VSphereAuthForm () {
    const { state, dispatch } = useContext(Store);
    const { register, handleSubmit, formState: { errors } } = useForm({
        resolver: yupResolver(authFormSchema)
    });
    const submitForm = (data: {[key: string]: string}) => {
        dispatch({
            type: SUBMIT_FORM,
            payload: data
        });
    };
    return (
        <CdsFormGroup layout="vertical-inline" control-width="shrink">
            <div cds-layout="horizontal gap:lg align:vertical-center">
                <CdsInput>
                    <label>VCENTER SERVER</label>
                    <input placeholder="IP OR FQDN" 
                        {...register('VCENTER_SERVER')}
                        defaultValue={state.data['VCENTER_SERVER']||''}/>
                    { errors['VCENTER_SERVER'] && <CdsControlMessage status="error">{errors['VCENTER_SERVER'].message}</CdsControlMessage> }
                </CdsInput>
                
                <CdsInput>
                    <label>USERNAME</label>
                    <input placeholder="Username"
                        {...register('VCENTER_USERNAME')}
                        defaultValue={state.data['VCENTER_USERNAME']||''}/>
                    { errors.VCENTER_USERNAME && <CdsControlMessage status="error">{errors.VCENTER_USERNAME.message}</CdsControlMessage> }
                </CdsInput>
                <CdsInput>
                    <label>PASSWORD</label>
                    <input placeholder="Password"
                        {...register('VCENTER_PASSWORD')}
                        type="password"
                        defaultValue={state.data['VCENTER_PASSWORD']||''}/>
                    { errors.VCENTER_PASSWORD && <CdsControlMessage status="error">{errors.VCENTER_PASSWORD.message}</CdsControlMessage> }
                </CdsInput>
            </div>
            <CdsButton onClick={handleSubmit((data) => submitForm(data))}>Next</CdsButton>
        </CdsFormGroup>
    );
}

export default VSphereAuthForm;