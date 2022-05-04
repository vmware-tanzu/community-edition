// React imports
import React, {
    createContext,
    ReactNode,
    useReducer,
} from 'react';

// App imports
import { StoreDispatch } from '../../shared/types/types';
import wcReducer from '../reducers/Wc.reducer';

const initialState = {
    data: {
        CLUSTER_CLASS_VARIABLE_VALUES: {},
        SELECTED_MANAGEMENT_CLUSTER: '',
        SELECTED_WORKER_NODE_INSTANCE_TYPE: '',
        WORKLOAD_CLUSTER_NAME: '',
    },
    ui: {
        wcCcOptionalExpanded: false,
        wcCcRequiredExpanded: true,
    }
};

const WcStore = createContext<{
    state: { [key: string]: any };
    dispatch: StoreDispatch;
}>({
    state: initialState,
    dispatch: () => null,
});

const WcProvider: React.FC<{ children: ReactNode }> = ({ children }: { children: ReactNode; }) => {
    const [state, dispatch] = useReducer(
        wcReducer,
        initialState
    );

    return (
        <WcStore.Provider value={{ state, dispatch }}>
            {children}
        </WcStore.Provider>
    );
};

export { WcStore, WcProvider };
