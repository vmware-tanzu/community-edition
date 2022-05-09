import { Dispatch, Reducer, ReducerAction } from 'react';

export interface Action {
    type: string,
    field: string,
    payload?: any,
    dataPath?: string,
    removeFieldIfEmpty?: boolean,
}
export type StoreDispatch = Dispatch<ReducerAction<Reducer<any, any>>>;
