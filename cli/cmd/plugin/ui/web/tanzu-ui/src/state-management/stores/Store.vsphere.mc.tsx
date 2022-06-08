// React imports
import React, { createContext, ReactNode, useReducer } from 'react';
// App imports
import { StoreDispatch } from '../../shared/types/types';
import wizardReducer from '../reducers/Wizard.reducer';

const initialState = {
    data: {
        SERVER: '',
        USERNAME: '',
        PASSWORD: '',
    },
};

const VsphereStore = createContext<{
    vsphereState: { [key: string]: any };
    vsphereDispatch: StoreDispatch;
}>({
    vsphereState: initialState,
    vsphereDispatch: () => null,
});

const VsphereProvider: React.FC<{ children: ReactNode }> = ({ children }: { children: ReactNode }) => {
    const [vsphereState, vsphereDispatch] = useReducer(wizardReducer, initialState);

    return <VsphereStore.Provider value={{ vsphereState, vsphereDispatch }}>{children}</VsphereStore.Provider>;
};

export { VsphereStore, VsphereProvider };
