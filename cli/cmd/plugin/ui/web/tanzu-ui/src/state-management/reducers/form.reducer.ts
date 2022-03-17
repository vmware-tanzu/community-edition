// App imports
import { TEXT_CHANGE } from '../actions/actionTypes';
import { Action } from '../../types/types';

interface FormState {
    VCENTER_SERVER?: string,
    VCENTER_USERNAME?: string,
    VCENTER_PASSWORD?: string
}

export function formReducer (state: FormState, action: Action) {
    let newState = { ...state };
    switch (action.type) {
    case TEXT_CHANGE:
        newState =  {
            [action.payload.name]: action.payload.value
        };
    }
    return newState;
}