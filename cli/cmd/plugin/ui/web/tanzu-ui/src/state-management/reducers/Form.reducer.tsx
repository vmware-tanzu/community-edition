// App imports
import {
    CCVAR_CHANGE,
    INPUT_CHANGE,
} from '../actions/Form.actions';
import { DEPLOYMENT_STATUS_CHANGED } from '../actions/Deployment.actions';
import { Action } from '../../shared/types/types';
import { ensureDataPath } from './index';
import { Deployments } from '../../shared/models/Deployments';

interface FormState {
    [key: string]: any;
    deployments: Deployments;
}

function createNewState(state: FormState, action: Action): FormState {
    return {
        ...state,
        [action.field]: action.payload
    }
}

function createNewCcVarState(state: FormState, action: Action): FormState {
    const newState = { ...state }
    const clusterName = action.locationData
    if (!clusterName) {
        console.error(
            `Form reducer unable to store ccvar data from this action: ${JSON.stringify(action)}, because no cluster name was provided!`)
    }
    const dataPath = `ccAttributes.${clusterName}`
    const leafObject = ensureDataPath(dataPath, newState)
    if (!action.payload) {
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
    case CCVAR_CHANGE:
        newState = createNewCcVarState(state, action)
        break;
    case DEPLOYMENT_STATUS_CHANGED:
        newState = {
            ...state,
            deployments: {
                ...action.payload
            }
        }
        break
    default:
        newState = { ...state };
    }
    console.log(`New state: ${JSON.stringify(newState)}`);
    return newState;
}
