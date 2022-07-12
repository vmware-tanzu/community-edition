// React imports
import React, { createContext, ReactNode, useReducer } from 'react';

// App imports
import { STORE_SECTION_FORM } from '../reducers/Form.reducer';
import { StoreDispatch } from '../../shared/types/types';
import umcReducer from '../reducers/Wizard.reducer';

const initialState = {
    [STORE_SECTION_FORM]: {
        // Cluster settings basic
        CLUSTER_NAME: '',

        // Cluster Settings advanced
        CONTROL_PLANE_NODES_COUNT: '1',
        WORKER_NODES_COUNT: '0',
        CLUSTER_PROVIDER: 'CALICO',

        // Cluster Network Settings
        IP_ADDRESS: '127.0.0.1',
        HOST_PORT_MAPPING: '80',
        NODE_PORT_MAPPING: '80',
        NODE_TO_HOST_PORT_MAPPING: '127.0.0.1/80:80/tcp',
        CLUSTER_SERVICE_CIDR: '100.64.0.0/13',
        CLUSTER_POD_CIDR: '100.96.0.0/11',
        CLUSTER_PROTOCOL: 'tcp',
    },
};

const UmcStore = createContext<{
    umcState: { [key: string]: any };
    umcDispatch: StoreDispatch;
}>({
    umcState: initialState,
    umcDispatch: () => null,
});

const UmcProvider: React.FC<{ children: ReactNode }> = ({ children }: { children: ReactNode }) => {
    const [umcState, umcDispatch] = useReducer(umcReducer, initialState);

    return <UmcStore.Provider value={{ umcState, umcDispatch }}>{children}</UmcStore.Provider>;
};

export { UmcStore, UmcProvider };
