// App imports
import { SUBMIT_FORM } from '../actions/actionTypes';
import { Action } from '../../types/types';

interface FormState {
    VCENTER_SERVER?: string,
    VCENTER_USERNAME?: string,
    VCENTER_PASSWORD?: string
}

export function formReducer (state: FormState, action: Action) {
    let newState = { ...state };
    switch (action.type) {
    case SUBMIT_FORM:
        newState =  {
            ...action.payload
        };
    }
    console.log(newState);
    return newState;
}