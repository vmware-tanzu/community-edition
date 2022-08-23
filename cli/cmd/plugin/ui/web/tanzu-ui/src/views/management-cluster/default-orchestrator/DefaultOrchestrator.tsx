import { addErrorInfo, removeErrorInfo } from '../../../shared/utilities/Error.util';
import { CancelablePromise } from '../../../swagger-api';
import { RESOURCE } from '../../../state-management/actions/Resources.actions';
import { ResourceAction, ResourceWithDefaultAction, StoreDispatch } from '../../../shared/types/types';

export interface InitResourcesProps<RESOURCE_TYPE> {
    dispatch: StoreDispatch;
    errorObject: { [fieldName: string]: any };
    setErrorObject: (errorObject: { [fieldName: string]: any }) => void;
    fetcher: () => CancelablePromise<RESOURCE_TYPE[]>;
    resourceName: string;
    fieldName?: string;
    fxnSelectDefault?: (resources: RESOURCE_TYPE[]) => RESOURCE_TYPE | undefined;
    fxnReturnValues?: (resources: RESOURCE_TYPE[], defaultValue?: RESOURCE_TYPE) => void;
}

export class DefaultOrchestrator {
    static async initResources<RESOURCE_TYPE>(props: InitResourcesProps<RESOURCE_TYPE>): Promise<RESOURCE_TYPE[]> {
        const { dispatch, errorObject, setErrorObject, resourceName, fxnSelectDefault, fxnReturnValues } = props;
        const fieldName = props.fieldName ?? props.resourceName;
        try {
            console.log(`initResources() calling fetcher for ${fieldName}`);
            const resources = (await props.fetcher()) || [];
            console.log(`initResources() fetcher received ${resources.length} item(s) for ${fieldName}`);
            if (fxnSelectDefault) {
                saveResourceDataWithDefault<RESOURCE_TYPE>({ resources, resourceName, fieldName, fxnSelectDefault, dispatch });
            } else {
                saveResourceData(dispatch, RESOURCE.ADD_RESOURCES, resourceName, resources);
            }
            setErrorObject(removeErrorInfo(errorObject, fieldName));
            if (fxnReturnValues) {
                fxnReturnValues(resources, fxnSelectDefault ? fxnSelectDefault(resources) : undefined);
            }
            return resources;
        } catch (error) {
            clearPreviousResourceData(dispatch, RESOURCE.DELETE_RESOURCES, resourceName);
            setErrorObject(addErrorInfo(errorObject, error, fieldName));
            if (fxnReturnValues) {
                fxnReturnValues([], undefined);
            }
            return [];
        }
    }

    static clearResourceData(dispatch: StoreDispatch, resourceName: string) {
        dispatch({
            type: RESOURCE.DELETE_RESOURCES,
            resourceName: resourceName,
            payload: [],
        } as ResourceAction);
    }
}

export function clearPreviousResourceData(dispatch: StoreDispatch, actionType: RESOURCE, resourceName: string) {
    dispatch({
        type: actionType,
        resourceName: resourceName,
        payload: [],
    } as ResourceAction);
}

export function saveResourceData(dispatch: StoreDispatch, actionType: RESOURCE, resourceName: string, currentValues: any[]) {
    console.log(`saveResourceData() dispatching event for ${resourceName}`);

    dispatch({
        type: actionType,
        resourceName: resourceName,
        payload: currentValues,
    } as ResourceAction);
}

interface ResourceDataWithDefaultParams<RESOURCE_TYPE> {
    dispatch: StoreDispatch;
    resourceName: string;
    fieldName?: string;
    resources: RESOURCE_TYPE[];
    fxnSelectDefault: (resources: RESOURCE_TYPE[]) => RESOURCE_TYPE | undefined;
}

function saveResourceDataWithDefault<RESOURCE_TYPE>(params: ResourceDataWithDefaultParams<RESOURCE_TYPE>) {
    console.log(`saveResourceDataWithDefault() dispatching event for ${params.resourceName}`);
    params.dispatch({
        type: RESOURCE.ADD_RESOURCES_WITH_DEFAULT,
        resourceName: params.resourceName,
        payload: params.resources,
        fieldName: params.fieldName,
        defaultValue: params.fxnSelectDefault(params.resources),
    } as ResourceWithDefaultAction);
}
