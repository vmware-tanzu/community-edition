import { Dispatch, Reducer, ReducerAction } from 'react';

export interface Action {
    type: string,
    field: string,
    payload?: any
}
export type StoreDispatch = Dispatch<ReducerAction<Reducer<any, any>>>;