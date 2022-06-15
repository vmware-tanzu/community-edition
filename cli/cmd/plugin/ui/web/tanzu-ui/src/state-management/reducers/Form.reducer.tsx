// App imports
import { DEPLOYMENT_STATUS_CHANGED } from '../actions/Deployment.actions';
import { Deployments } from '../../shared/models/Deployments';
import { FormAction } from '../../shared/types/types';
import { INPUT_CHANGE } from '../actions/Form.actions';
import { ReducerDescriptor } from '../../shared/utilities/Reducer.utils';

export const STORE_SECTION_FORM = 'dataForm';

interface FormState {
    [key: string]: any;
    deployments: Deployments;
}

function createNewState(state: FormState, action: FormAction): FormState {
    return {
        ...state,
        [action.field]: action.payload,
    };
}

function formReducer(state: FormState, action: FormAction) {
    let newState;
    switch (action.type) {
        case INPUT_CHANGE:
            newState = createNewState(state, action);
            break;
        case DEPLOYMENT_STATUS_CHANGED:
            newState = {
                ...state,
                deployments: {
                    ...action.payload,
                },
            };
            break;
        default:
            newState = { ...state };
    }
    console.log(`New form state: ${JSON.stringify(newState)}`);
    return newState;
}

export const formReducerDescriptor = {
    name: 'form reducer',
    reducer: formReducer,
    actionTypes: [INPUT_CHANGE, DEPLOYMENT_STATUS_CHANGED],
    storeSection: STORE_SECTION_FORM,
} as ReducerDescriptor;
