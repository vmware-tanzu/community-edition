// App imports
import { FormAction } from '../../shared/types/types';
import { INPUT_CHANGE, INPUT_CLEAR, SET_DEFAULTS } from '../actions/Form.actions';
import { ReducerDescriptor } from '../../shared/utilities/Reducer.utils';

export const STORE_SECTION_FORM = 'dataForm';

interface FormState {
    [key: string]: any;
}

function setInputField(state: FormState, action: FormAction): FormState {
    return {
        ...state,
        [action.field]: action.payload,
    };
}

function clearInputField(state: FormState, action: FormAction): FormState {
    const result = { ...state };
    delete result[action.field];
    return result;
}

function formReducer(state: FormState, action: FormAction) {
    let newState;
    switch (action.type) {
        case INPUT_CHANGE:
            newState = setInputField(state, action);
            break;
        case INPUT_CLEAR:
            newState = clearInputField(state, action);
            break;
        case SET_DEFAULTS:
            newState = { ...state, ...action.payload };
            break;
        default:
            console.error(`formReducer ignoring unrecognized action: ${JSON.stringify(action)}`);
            newState = { ...state };
    }
    console.log(`After ${action.type}, new form state: ${JSON.stringify(newState)}`);
    return newState;
}

export const formReducerDescriptor = {
    name: 'form reducer',
    reducer: formReducer,
    actionTypes: [INPUT_CHANGE, INPUT_CLEAR, SET_DEFAULTS],
    storeSection: STORE_SECTION_FORM,
} as ReducerDescriptor;
