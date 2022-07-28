// React imports
import React, { createContext, ReactNode, useReducer } from 'react';

// App imports
import { StoreDispatch } from '../../shared/types/types';
import azureReducer from '../reducers/Wizard.reducer';
import { STORE_SECTION_FORM } from '../reducers/Form.reducer';
import { AzureCloud } from '../../shared/constants/App.constants';
import { AZURE_DEFAULT_VALUES } from '../../shared/constants/defaults/azure.defaults';

const initialState = {
    [STORE_SECTION_FORM]: {
        TENANT_ID: '',
        CLIENT_ID: '',
        CLIENT_SECRET: '',
        SUBSCRIPTION_ID: '',
        AZURE_ENVIRONMENT: AzureCloud.PUBLIC,
        REGION: '',
        SSH_PUBLIC_KEY: '',
        IMAGE_INFO: {},
        ...AZURE_DEFAULT_VALUES,
    },
};

const AzureStore = createContext<{
    azureState: { [key: string]: any };
    azureDispatch: StoreDispatch;
}>({
    azureState: initialState,
    azureDispatch: () => null,
});

const AzureProvider: React.FC<{ children: ReactNode }> = ({ children }: { children: ReactNode }) => {
    const [azureState, azureDispatch] = useReducer(azureReducer, initialState);

    return <AzureStore.Provider value={{ azureState, azureDispatch }}>{children}</AzureStore.Provider>;
};

export { AzureStore, AzureProvider };
