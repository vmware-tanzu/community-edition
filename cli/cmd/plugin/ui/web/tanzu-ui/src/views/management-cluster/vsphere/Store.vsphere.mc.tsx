// React imports
import React, { createContext, ReactNode, useReducer } from 'react';

// App imports
import { STORE_SECTION_FORM } from '../../../state-management/reducers/Form.reducer';
import { STORE_SECTION_RESOURCES } from '../../../state-management/reducers/Resources.reducer';
import { StoreDispatch } from '../../../shared/types/types';
import { VSPHERE_FIELDS } from './VsphereManagementCluster.constants';
import vsphereReducer from '../../providers/vsphere/Vsphere.reducer';

const initialState = {
    [STORE_SECTION_FORM]: {
        [VSPHERE_FIELDS.SERVERNAME]: '',
        [VSPHERE_FIELDS.USERNAME]: '',
        [VSPHERE_FIELDS.PASSWORD]: '',
        [VSPHERE_FIELDS.DATACENTER]: undefined,
    },
    [STORE_SECTION_RESOURCES]: {},
};

const VsphereStore = createContext<{
    vsphereState: { [key: string]: any };
    vsphereDispatch: StoreDispatch;
}>({
    vsphereState: initialState,
    vsphereDispatch: () => null,
});

const VsphereProvider: React.FC<{ children: ReactNode }> = ({ children }: { children: ReactNode }) => {
    const [vsphereState, vsphereDispatch] = useReducer(vsphereReducer, initialState);

    return <VsphereStore.Provider value={{ vsphereState, vsphereDispatch }}>{children}</VsphereStore.Provider>;
};

export { VsphereStore, VsphereProvider };
