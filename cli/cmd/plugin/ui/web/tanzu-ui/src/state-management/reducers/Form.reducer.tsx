// App imports
import { RESET_DEPENDENT_FIELDS, SUBMIT_FORM } from '../actions/Form.actions';
import { Action } from '../../shared/types/types';

interface FormState {
    VCENTER_SERVER?: string,
    VCENTER_USERNAME?: string,
    VCENTER_PASSWORD?: string
}

export function formReducer (state: FormState, action: Action) {
    let newState = { ...state };
    if (action.type === SUBMIT_FORM) {
        newState =  {
            ...newState,
            ...action.payload
        };
    } else if (action.type === RESET_DEPENDENT_FIELDS) {
        const resetFields = action.payload.fields.reduce((acc: {[key: string]: string}, cur: string) => {
            acc[cur] = '';
            return acc;
        }, {});
        newState = {
            ...newState,
            ...resetFields
        };
    }
    console.log(newState);
    return newState;
}