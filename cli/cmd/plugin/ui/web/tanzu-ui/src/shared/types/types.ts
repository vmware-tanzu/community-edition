import { Dispatch, Reducer, ReducerAction } from 'react';

export interface Action {
    type: string; // type of action, e.g. INPUT_CHANGE
    field: string; // name of form field related to the action
    payload?: any; // the payload of the action, generally the new value
    locationData?: any; // data needed for storing the payload, generally only used when store location is dynamic (cf CCVAR_CHANGE)
}
export type StoreDispatch = Dispatch<ReducerAction<Reducer<any, any>>>;
