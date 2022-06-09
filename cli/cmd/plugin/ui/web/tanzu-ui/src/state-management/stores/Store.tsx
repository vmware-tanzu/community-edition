// React imports
import React, { createContext, Dispatch, ReactNode, Reducer, ReducerAction, useReducer } from 'react';

// App imports
import mainReducer from '../reducers';

const initialState = {
    app: {
        appEnv: '',
        appRoute: '',
    },
    ui: {
        navExpanded: true,
        isDeployInProgress: false,
        currentRoute: '',
        appBanner: {
            display: false,
            message: '',
            text: '',
            status: 'success',
        },
    },
    data: {
        // TODO: convert to list of deployments; should be updated when deployment started
        deployments: {
            type: 'management-cluster',
            status: '',
            provider: 'aws',
            configPath: '~/.config/tanzu/tkg/clusterconfigs/fcrjpbtumf.yaml',
        },
    },
};
export type StoreDispatch = Dispatch<ReducerAction<Reducer<any, any>>>;
const Store = createContext<{
    state: { [key: string]: any };
    dispatch: StoreDispatch;
}>({
    state: initialState,
    dispatch: () => null,
});

const AppProvider: React.FC<{ children: ReactNode }> = ({ children }: { children: ReactNode }) => {
    const [state, dispatch] = useReducer(mainReducer, initialState);

    return <Store.Provider value={{ state, dispatch }}>{children}</Store.Provider>;
};

export { Store, AppProvider };
