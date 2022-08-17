// React imports
import React, { createContext, ReactNode, useReducer } from 'react';

// App imports
import { StoreDispatch } from '../../shared/types/types';
import azureReducer from '../../views/providers/azure/Azure.reducer';
import { STORE_SECTION_FORM } from '../reducers/Form.reducer';
import { AzureCloud } from '../../shared/constants/App.constants';
import { AZURE_DEFAULT_VALUES } from '../../shared/constants/defaults/azure.defaults';
import { STORE_SECTION_AZURE_RESOURCES } from '../../views/providers/azure/AzureResources.reducer';
import { AZURE_FIELDS } from '../../views/management-cluster/azure/azure-mc-basic/AzureManagementClusterBasic.constants';

const initialState = {
    [STORE_SECTION_FORM]: {
        [AZURE_FIELDS.TENANT_ID]: '',
        [AZURE_FIELDS.CLIENT_ID]: '',
        [AZURE_FIELDS.CLIENT_SECRET]: '',
        [AZURE_FIELDS.SUBSCRIPTION_ID]: '',
        [AZURE_FIELDS.AZURE_ENVIRONMENT]: AzureCloud.PUBLIC,
        [AZURE_FIELDS.REGION]: '',
        [AZURE_FIELDS.SSH_PUBLIC_KEY]: '',
        [AZURE_FIELDS.OS_IMAGE]: {},
        [AZURE_FIELDS.CLUSTER_NAME]: '',
        ...AZURE_DEFAULT_VALUES,
    },
    [STORE_SECTION_AZURE_RESOURCES]: {},
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
