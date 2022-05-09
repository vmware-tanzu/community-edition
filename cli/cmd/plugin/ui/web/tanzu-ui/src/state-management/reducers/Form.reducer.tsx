// App imports
import {
    INPUT_CHANGE,
/*
    RESET_DEPENDENT_FIELDS,
    SUBMIT_FORM,
*/
} from '../actions/Form.actions';
import { Action } from '../../shared/types/types';
import { ensureDataPath } from './index';

interface FormState {
    [key: string]: any;
}

function createNewState(state: FormState, action: Action): FormState {
    const newState = { ...state }
    const leafObject = ensureDataPath(action.dataPath, newState)
    if (action.removeFieldIfEmpty && !action.payload) {
        delete leafObject[action.field]
    } else {
        leafObject[action.field] = action.payload
    }
    return newState
}

export function formReducer(state: FormState, action: Action) {
    let newState;
    switch (action.type) {
    case INPUT_CHANGE:
        newState = createNewState(state, action)
        break;
    default:
        newState = { ...state };
    }
    console.log(`New state: ${JSON.stringify(newState)}`);
    return newState;
}
