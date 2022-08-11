import { AzureResourceAction } from '../../../shared/types/types';
import { ReducerDescriptor } from '../../../shared/utilities/Reducer.utils';
import { AZURE_ADD_RESOURCES, AZURE_DELETE_RESOURCES } from '../../../state-management/actions/Resources.actions';

export const STORE_SECTION_AZURE_RESOURCES = 'resources';

// The resources reducer builds up state by ensuring an object associated with action.region,
// and then assigning object[action.resourceName] = action.payload.
// So if the action object were:
// { region: us-east-1, resourceName: osImages, payload: [obj1, obj2, obj3] }, then we would expect
// to see state['us-east-1']['osImages'] === [obj1, obj2, obj3]
function azureResourcesReducer(state: any, action: AzureResourceAction) {
    const newState = { ...state };
    if (!action.resourceName) {
        console.error(`azureResourcesReducer received action ${JSON.stringify(action)} which has no resourceName!`);
        return newState;
    }
    if (action.type === AZURE_ADD_RESOURCES) {
        newState[action.resourceName] = action.payload;
    } else if (action.type === AZURE_DELETE_RESOURCES) {
        delete newState[action.resourceName];
    }
    console.log(`New resources state: ${JSON.stringify(newState)} after action ${JSON.stringify(action)}`);
    return newState;
}

export function getResource(resourceName: string, store: any) {
    if (!store || !store[STORE_SECTION_AZURE_RESOURCES] || !resourceName) {
        return undefined;
    }
    return store[STORE_SECTION_AZURE_RESOURCES][resourceName];
}

export const azureResourceReducerDescriptor = {
    name: 'azure resource reducer',
    reducer: azureResourcesReducer,
    storeSection: STORE_SECTION_AZURE_RESOURCES,
    actionTypes: [AZURE_ADD_RESOURCES, AZURE_DELETE_RESOURCES],
} as ReducerDescriptor;
