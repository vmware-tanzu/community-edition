import React, { createContext, Dispatch, ReactNode, Reducer, ReducerAction, useReducer } from 'react';
import mainReducer from '../reducers';

const initialState = {
    app: {
        appEnv: ''
    },
    ui: {
        navExpanded: false,
        isDeployInProgress: false,
        appBanner: {
            display: false,
            message: '',
            text: '',
            status: 'success'
        }
    },
    data: {
        // TODO: should be a list of deployments
        deployments: {
            type: '',
            status: '',
            provider: '',
            configPath: ''
        },
        VCENTER_SERVER: '1.1.1.1',
        VCENTER_USERNAME: 'admin',
        VCENTER_PASSWORD: 'password',
        CLUSTER_NAME: 'mycluster',

    }
};
export type StoreDispatch = Dispatch<ReducerAction<Reducer<any, any>>>;
const Store = createContext<{
    state: {[key: string]: any},
    dispatch: StoreDispatch
}>({
    state: initialState,
    dispatch: () => null
});


const AppProvider: React.FC<{ children: ReactNode}> = ({ children } : { children: ReactNode}) => {
    const [state, dispatch] = useReducer(mainReducer, initialState);

    return (
        <Store.Provider value={{ state, dispatch }}>
            {children}
        </Store.Provider>
    );
};

export { Store, AppProvider };