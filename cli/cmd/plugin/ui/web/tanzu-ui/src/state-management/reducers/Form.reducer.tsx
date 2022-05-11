// App imports
import {
    CCVAR_CHANGE,
    INPUT_CHANGE,
} from '../actions/Form.actions';
import { Action } from '../../shared/types/types';
import { ensureDataPath } from './index';

interface FormState {
    [key: string]: any;
}

// NOTE: the field path separate cannot be a period, because yup chokes when a field name has a period in it
export const FIELD_PATH_SEPARATOR = '___'
const DATASTORE_PATH_SEPARATOR = '.'

function createNewState(state: FormState, action: Action): FormState {
    return {
        ...state,
        [action.field]: action.payload
    }
}

function createNewCcVarState(state: FormState, action: Action): FormState {
    const newState = { ...state }
    const clusterName = action.locationData.clusterName
    if (!clusterName) {
        console.error(
            `Form reducer unable to store ccvar data from this action: ${JSON.stringify(action)}, because no cluster name was provided!`)
    }
    // for the data path, take the name of the field, parse on path separator and throw away the last segment
    // for the simpleFieldName, use the last segment
    const pathParts = action.field.split(FIELD_PATH_SEPARATOR)
    const path = pathParts.length > 1 ? pathParts.slice(0, pathParts.length - 1).join(DATASTORE_PATH_SEPARATOR) : ''
    const simpleFieldName = pathParts[pathParts.length-1]

    const dataPath = `ccAttributes.${clusterName}${path ? DATASTORE_PATH_SEPARATOR + path : ''}`
    const leafObject = ensureDataPath(dataPath, DATASTORE_PATH_SEPARATOR, newState)
    if (!action.payload) {
        delete leafObject[simpleFieldName]
    } else {
        leafObject[simpleFieldName] = action.payload
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
    default:
        newState = { ...state };
    }
    console.log(`New state: ${JSON.stringify(newState)}`);
    return newState;
}
