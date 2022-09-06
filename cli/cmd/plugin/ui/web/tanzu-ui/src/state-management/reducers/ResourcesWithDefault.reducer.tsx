import { FormAction, ResourceAction, ResourceWithDefaultAction } from '../../shared/types/types';
import { RESOURCE } from '../actions/Resources.actions';
import { ReducerDescriptor } from '../../shared/utilities/Reducer.utils';
import { resourceReducerDescriptor, STORE_SECTION_RESOURCES } from './Resources.reducer';
import { formReducerDescriptor, STORE_SECTION_FORM } from './Form.reducer';
import { INPUT_CHANGE, INPUT_CLEAR } from '../actions/Form.actions';

// This reducer updates TWO sections of the store: the resources section and the form section
function resourcesWithDefaultReducer(state: any, action: ResourceWithDefaultAction) {
    return updateFormSection(updateResourcesSection(state, action), action);
}

function updateResourcesSection(state: any, action: ResourceWithDefaultAction) {
    const newState = { ...state };
    const resourcesActionType = action.type === RESOURCE.DELETE_RESOURCES_WITH_DEFAULT ? RESOURCE.DELETE_RESOURCES : RESOURCE.ADD_RESOURCES;
    const resourcesAction = { ...action, type: resourcesActionType } as ResourceAction;
    newState[STORE_SECTION_RESOURCES] = resourceReducerDescriptor.reducer(state[STORE_SECTION_RESOURCES], resourcesAction);
    return newState;
}

function updateFormSection(state: any, action: ResourceWithDefaultAction) {
    const newState = { ...state };
    const field = action.fieldName ?? action.resourceName;
    const actionType = action.type === RESOURCE.DELETE_RESOURCES_WITH_DEFAULT ? INPUT_CLEAR : INPUT_CHANGE;
    const formAction = { type: actionType, field, payload: action.defaultValue } as FormAction;
    newState[STORE_SECTION_FORM] = formReducerDescriptor.reducer(newState[STORE_SECTION_FORM], formAction);
    return newState;
}

// NOTE: because this reducer updates TWO sections of the store, it does not specify a singular section in this descriptor
export const resourceWithDefaultReducerDescriptor = {
    name: 'resource-with-default reducer',
    reducer: resourcesWithDefaultReducer,
    actionTypes: [RESOURCE.ADD_RESOURCES_WITH_DEFAULT, RESOURCE.DELETE_RESOURCES_WITH_DEFAULT],
} as ReducerDescriptor;
