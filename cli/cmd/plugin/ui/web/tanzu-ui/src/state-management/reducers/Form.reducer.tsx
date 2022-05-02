// App imports
import {
    INPUT_CHANGE,
    RESET_DEPENDENT_FIELDS,
    SUBMIT_FORM,
} from '../actions/Form.actions';
import { Action } from '../../shared/types/types';

interface FormState {
    [key: string]: any;
}

export function formReducer(state: FormState, action: Action) {
    let newState = { ...state };
    switch (action.type) {
        case INPUT_CHANGE:
            newState = {
                ...state,
                [action.field]: action.payload,
            };
            break;
        default:
            newState = { ...state };
    }
    console.log(`New state: ${JSON.stringify(newState)}`);
    return newState;
    // let newState = { ...state };
    // if (action.type === SUBMIT_FORM) {
    //     newState = {
    //         ...newState,
    //         ...action.payload,
    //     };
    // } else if (action.type === RESET_DEPENDENT_FIELDS) {
    //     const resetFields = action.payload.fields.reduce(
    //         (acc: { [key: string]: string }, cur: string) => {
    //             acc[cur] = '';
    //             return acc;
    //         },
    //         {}
    //     );
    //     newState = {
    //         ...newState,
    //         ...resetFields,
    //     };
    // }
    // console.log(newState);
    // return newState;
}
