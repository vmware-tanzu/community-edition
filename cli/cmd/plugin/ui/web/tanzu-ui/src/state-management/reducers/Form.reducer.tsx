// App imports
import { CCVAR_CHANGE, INPUT_CHANGE } from '../actions/Form.actions';
import { DEPLOYMENT_STATUS_CHANGED } from '../actions/Deployment.actions';
import { Action } from '../../shared/types/types';
import { ensureDataPath, getDataPath, pruneDataPath } from './index';
import { Deployments } from '../../shared/models/Deployments';

interface FormState {
    [key: string]: any;
    deployments: Deployments;
}

// NOTE: the field path separate cannot be a period, because yup chokes when a field name has a period in it
export const FIELD_PATH_SEPARATOR = '___';
const DATASTORE_PATH_SEPARATOR = '.';

function createNewState(state: FormState, action: Action): FormState {
    return {
        ...state,
        [action.field]: action.payload,
    };
}

export function getFieldData(fieldName: string, clusterName: string, state: any): any {
    const { simpleFieldName, fieldDataPath } = parseFullFieldName(fieldName);
    const dataPath = 'data' + DATASTORE_PATH_SEPARATOR + fullDataPath(clusterName, fieldDataPath);
    const leafObject = getDataPath(dataPath, DATASTORE_PATH_SEPARATOR, state);
    return leafObject ? leafObject[simpleFieldName] : undefined;
}

// NOTE: the fullFieldName is a string separated by FIELD_PATH_SEPARATOR which is in the form:
// [category name]___path___with___segments___[field name]
// for the category name, take the first segment (we generally don't care about the category name)
// for the data path, throw away the first and last segments and take what's left
// for the simpleFieldName, use the last segment
function parseFullFieldName(fullFieldName: string): { category: string; simpleFieldName: string; fieldDataPath: string } {
    const pathPartsFromFieldName = fullFieldName.split(FIELD_PATH_SEPARATOR);
    const category = pathPartsFromFieldName[0];
    const simpleFieldName = pathPartsFromFieldName[pathPartsFromFieldName.length - 1];
    const fieldDataPath =
        pathPartsFromFieldName.length > 2
            ? pathPartsFromFieldName.slice(1, pathPartsFromFieldName.length - 1).join(DATASTORE_PATH_SEPARATOR)
            : '';
    return { category, fieldDataPath, simpleFieldName };
}

function fullDataPath(clusterName: string, fieldPath: string | undefined) {
    return `ccAttributes.${clusterName}${fieldPath ? DATASTORE_PATH_SEPARATOR + fieldPath : ''}`;
}

function createNewCcVarState(state: FormState, action: Action): FormState {
    const newState = { ...state };
    const clusterName = action.locationData.clusterName;
    if (!clusterName) {
        console.error(
            `Form reducer unable to store ccvar data from this action: ${JSON.stringify(action)}, because no cluster name was provided!`
        );
    }

    let fieldName = '';
    let fieldPath;
    if (action.locationData.fieldPath) {
        // we're using the fieldPath sent to us to get the datapath and the simple field name
        const dataPathElements = action.locationData.fieldPath.split(DATASTORE_PATH_SEPARATOR);
        // the last part of the path is the simple field name
        fieldName = dataPathElements[dataPathElements.length - 1];
        // all the first parts of the path together are the field path
        fieldPath = dataPathElements.slice(0, dataPathElements.length - 1).join(DATASTORE_PATH_SEPARATOR);
    } else {
        // we're using the field name of the action to find the data path and the simple field name
        const { simpleFieldName, fieldDataPath } = parseFullFieldName(action.field);
        fieldName = simpleFieldName;
        fieldPath = fieldDataPath;
    }

    const dataPath = fullDataPath(clusterName, fieldPath);
    const leafObject = ensureDataPath(dataPath, DATASTORE_PATH_SEPARATOR, newState);
    if (!action.payload) {
        delete leafObject[fieldName];
        pruneDataPath(dataPath, DATASTORE_PATH_SEPARATOR, newState);
    } else {
        leafObject[fieldName] = action.payload;
    }
    return newState;
}

export function formReducer(state: FormState, action: Action) {
    let newState;
    switch (action.type) {
        case INPUT_CHANGE:
            newState = createNewState(state, action);
            break;
        case CCVAR_CHANGE:
            newState = createNewCcVarState(state, action);
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
    console.log(`New state: ${JSON.stringify(newState)}`);
    return newState;
}
