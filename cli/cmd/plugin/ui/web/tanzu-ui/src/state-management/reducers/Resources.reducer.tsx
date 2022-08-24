import { ReducerDescriptor } from '../../shared/utilities/Reducer.utils';
import { RESOURCE } from '../actions/Resources.actions';
import { ResourceAction } from '../../shared/types/types';

export const STORE_SECTION_RESOURCES = 'resources';

function resourcesReducer(state: any, action: ResourceAction) {
    const newState = { ...state };
    if (!action.resourceName) {
        console.error(`resourcesReducer received action ${JSON.stringify(action)} which has no resourceName!`);
        return newState;
    }
    if (action.type === RESOURCE.ADD_RESOURCES) {
        newState[action.resourceName] = action.payload;
    } else if (action.type === RESOURCE.DELETE_RESOURCES) {
        delete newState[action.resourceName];
    } else {
        console.error(`resourcesReducer received unrecognized action type: ${JSON.stringify(action)}`);
    }
    console.log(`New resources state: ${JSON.stringify(newState)} after action ${JSON.stringify(action)}`);
    return newState;
}

export function getResource<RESOURCE_TYPE>(resourceName: string, store: any): RESOURCE_TYPE | undefined {
    if (!store || !store[STORE_SECTION_RESOURCES] || !resourceName) {
        return undefined;
    }
    return store[STORE_SECTION_RESOURCES][resourceName] as RESOURCE_TYPE;
}

export const resourceReducerDescriptor = {
    name: 'resource reducer',
    reducer: resourcesReducer,
    storeSection: STORE_SECTION_RESOURCES,
    actionTypes: [RESOURCE.ADD_RESOURCES, RESOURCE.DELETE_RESOURCES],
} as ReducerDescriptor;
