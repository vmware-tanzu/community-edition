// React imports
import React, { createContext, Dispatch, ReactNode, Reducer, ReducerAction, useReducer } from 'react';

// App imports
import mainReducer from '../reducers';
import { STORE_SECTION_APP } from '../reducers/App.reducer';
import { STORE_SECTION_DEPLOYMENT } from '../reducers/Deployment.reducer';
import { STORE_SECTION_UI } from '../reducers/Ui.reducer';

const initialState = {
    [STORE_SECTION_APP]: {
        appEnv: '',
        appRoute: '',
    },
    [STORE_SECTION_UI]: {
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
    [STORE_SECTION_DEPLOYMENT]: {
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
