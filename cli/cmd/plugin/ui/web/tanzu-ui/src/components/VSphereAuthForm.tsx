// React imports
import React, { ChangeEvent, useContext } from 'react';

// Library imports
import { CdsInput } from '@cds/react/input';
import { CdsFormGroup } from '@cds/react/forms';
import { CdsButton } from '@cds/react/button';

// App imports
import { TEXT_CHANGE } from '../state-management/actions/actionTypes';
import { Store } from '../state-management/stores/store';


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
                            onChange={(e) => serverOnchange(e, 'VCENTER_SERVER')}
                            value={state.data['VCENTER_SERVER']||''}/>
                    </CdsInput>
                    <CdsInput>
                        <label>USERNAME</label>
                        <input placeholder="Username"
                            onChange={(e) => { serverOnchange(e, 'VCENTER_USERNAME');}}
                            value={state.data['VCENTER_USERNAME']||''}/>
                    </CdsInput>
                    <CdsInput>
                        <label>PASSWORD</label>
                        <input placeholder="Password"
                            onChange={(e) => { serverOnchange(e, 'VCENTER_PASSWORD');}}
                            value={state.data['VCENTER_PASSWORD']||''}/>
                    </CdsInput>
                </div>
                <CdsButton>Next</CdsButton>
            </CdsFormGroup>
        </form>
    );
}

export default VSphereAuthForm;