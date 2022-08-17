import { Dispatch, Reducer, ReducerAction } from 'react';
import { RESOURCE } from '../../state-management/actions/Resources.actions';
import { AWS_FIELDS } from '../../views/management-cluster/aws/aws-mc-basic/AwsManagementClusterBasic.constants';
import { AZURE_FIELDS } from '../../views/management-cluster/azure/azure-mc-basic/AzureManagementClusterBasic.constants';
/**
 * Actions are used with the dispatch() functionality, which causes a message to be sent which will affect a store's state.
 * All dispatch() calls send an Action param, so there are many possible Action objects.
 * Note that because the StoreDispatch (defined below) uses the Action type, all custom reducers must accept an Action (which
 * may be an extension of the original Action interface below)
 */
export interface Action {
    type: string; // type of action, e.g. INPUT_CHANGE
    payload?: any; // the payload of the action, generally the new value
}

export interface FormAction extends Action {
    field: string; // name of form field related to the action
}

export interface ResourceAction extends Action {
    resourceName: string;
}

export interface DynamicFormAction extends FormAction {
    locationData?: any; // data needed for storing the payload, generally only used when store location is dynamic (cf CCVAR_CHANGE)
}

export interface DynamicCategoryToggleAction extends Action {
    category: string;
}

export interface KeyOfStringToString {
    [key: string]: string;
}

export type StoreDispatch = Dispatch<ReducerAction<Reducer<any, Action>>>;

// export interface Fields {
//     fields: AZURE_FIELDS | AWS_FIELDS;
// }
