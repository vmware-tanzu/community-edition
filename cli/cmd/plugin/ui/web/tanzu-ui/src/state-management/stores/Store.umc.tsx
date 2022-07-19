// React imports
import React, { createContext, ReactNode, useReducer } from 'react';

// App imports
import { STORE_SECTION_FORM } from '../reducers/Form.reducer';
import { StoreDispatch } from '../../shared/types/types';
import umcReducer from '../reducers/Wizard.reducer';
import { UNMANAGED_CLUSTER_FIELDS } from '../../views/unmanaged-cluster/unmanaged-cluster-wizard-page/UnmanagedCluster.constants';

const initialState = {
    [STORE_SECTION_FORM]: {
        // Cluster settings basic
        [UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]: 'unmanged-cluster',

        // Cluster Settings advanced
        [UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT]: '1',
        [UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT]: '0',
        [UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROVIDER]: 'kind',

        // Cluster Network Settings
        [UNMANAGED_CLUSTER_FIELDS.CNI_PROVIDER]: 'CALICO',
        [UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR]: '100.64.0.0/13',
        [UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR]: '100.96.0.0/11',
        [UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS]: '127.0.0.1',
        [UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING]: '80',
        [UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING]: '80',
        [UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROTOCOL]: 'tcp',
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
