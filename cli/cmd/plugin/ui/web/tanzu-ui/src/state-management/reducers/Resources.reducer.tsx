import { ReducerDescriptor } from '../../shared/utilities/Reducer.utils';
import { RESOURCE } from '../actions/Resources.actions';
import { ResourceAction } from '../../shared/types/types';

export const STORE_SECTION_RESOURCES = 'resources';

function addResource(state: any, action: ResourceAction) {
    const newState = { ...state };
    if (action.segment) {
        if (!newState[action.resourceName]) {
            newState[action.resourceName] = {};
        }
        newState[action.resourceName][action.segment] = action.payload;
    } else {
        newState[action.resourceName] = action.payload;
    }
    return newState;
}
function deleteResource(state: any, action: ResourceAction) {
    const newState = { ...state };
    if (action.segment) {
        if (newState[action.resourceName]) {
            delete newState[action.resourceName][action.segment];
        }
    } else {
        delete newState[action.resourceName];
    }
    return newState;
}
function resourcesReducer(state: any, action: ResourceAction) {
    if (!action.resourceName) {
        console.error(`resourcesReducer received action ${JSON.stringify(action)} which has no resourceName!`);
        return { ...state };
    }
    let newState;
    if (action.type === RESOURCE.ADD_RESOURCES) {
        newState = addResource(state, action);
    } else if (action.type === RESOURCE.DELETE_RESOURCES) {
        newState = deleteResource(state, action);
    } else {
        console.error(`resourcesReducer ignoring unrecognized action type: ${JSON.stringify(action)}`);
        newState = { ...state };
    }
    console.log(`After action ${JSON.stringify(action)} -->\nNew resources state: ${JSON.stringify(newState)}`);
    return newState;
}

export function getResource<RESOURCE_TYPE>(resourceName: string, store: any, segment?: string): RESOURCE_TYPE | undefined {
    if (!store || !store[STORE_SECTION_RESOURCES] || !resourceName) {
        return undefined;
    }
    if (!segment) {
        return store[STORE_SECTION_RESOURCES][resourceName] as RESOURCE_TYPE;
    }
    if (!store[STORE_SECTION_RESOURCES][resourceName]) {
        return undefined;
    }
    return store[STORE_SECTION_RESOURCES][resourceName][segment] as RESOURCE_TYPE;
}

export const resourceReducerDescriptor = {
    name: 'resource reducer',
    reducer: resourcesReducer,
    storeSection: STORE_SECTION_RESOURCES,
    actionTypes: [RESOURCE.ADD_RESOURCES, RESOURCE.DELETE_RESOURCES],
} as ReducerDescriptor;
