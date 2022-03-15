import { TEXT_CHANGE } from '../constants/actionTypes';
import { Action } from '../types/types';

interface FormState {
    VCENETER_SERVER?: string,
    VCENETER_USERNAME?: string,
    VCENETER_PASSWORD?: string
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