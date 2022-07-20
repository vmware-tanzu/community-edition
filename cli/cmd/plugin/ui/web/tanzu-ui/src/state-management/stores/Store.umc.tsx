// React imports
import React, { createContext, ReactNode, useReducer } from 'react';

// App imports
import { STORE_SECTION_FORM } from '../reducers/Form.reducer';
import { StoreDispatch } from '../../shared/types/types';
import umcReducer from '../reducers/Wizard.reducer';
import { UNMANAGED_CLUSTER_FIELDS } from '../../views/unmanaged-cluster/unmanaged-cluster-wizard-page/UnmanagedCluster.constants';
import { UNMANAGED_DEFAULT_VALUES } from '../../shared/constants/defaults/unmanaged.defaults';

const initialState = {
    [STORE_SECTION_FORM]: {
        // Cluster settings basic
        [UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME]: '',
        ...UNMANAGED_DEFAULT_VALUES,
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
