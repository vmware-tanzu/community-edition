import React, { ChangeEvent, useContext } from 'react';

import { CdsInput } from '@cds/react/input';
import { CdsFormGroup } from '@cds/react/forms';
import { CdsButton } from '@cds/react/button';

import { TEXT_CHANGE } from '../constants/actionTypes';
import { Store } from '../stores/store';


function VSphereAuthForm () {
    const { state, dispatch } = useContext(Store);
    const serverOnchange = (e: ChangeEvent<HTMLInputElement>, fieldName: string) => {
        dispatch({
            type: TEXT_CHANGE,
            payload: {
                name: fieldName,
                value: e.target.value
            }
        });
    };
    return (
        <form>
            <CdsFormGroup layout="vertical-inline" control-width="shrink">
                <div cds-layout="horizontal gap:lg align:vertical-center">
                    <CdsInput>
                        <label>VCENTER SERVER</label>
                        <input placeholder="IP OR FQDN" 
                            onChange={(e) => serverOnchange(e, 'VCENETER_SERVER')}
                            value={state.data['VCENETER_SERVER']||''}/>
                    </CdsInput>
                    <CdsInput>
                        <label>USERNAME</label>
                        <input placeholder="Username"
                            onChange={(e) => { serverOnchange(e, 'VCENETER_USERNAME');}}
                            value={state.data['VCENETER_USERNAME']||''}/>
                    </CdsInput>
                    <CdsInput>
                        <label>PASSWORD</label>
                        <input placeholder="Password"
                            onChange={(e) => { serverOnchange(e, 'VCENETER_PASSWORD');}}
                            value={state.data['VCENETER_PASSWORD']||''}/>
                    </CdsInput>
                </div>
                <CdsButton>Next</CdsButton>
            </CdsFormGroup>
        </form>
    );
}

export default VSphereAuthForm;