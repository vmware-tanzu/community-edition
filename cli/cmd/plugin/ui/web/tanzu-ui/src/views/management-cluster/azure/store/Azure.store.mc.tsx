// React imports
import React, { createContext, ReactNode, useReducer } from 'react';

// App imports
import azureReducer from '../../../../views/providers/azure/Azure.reducer';
import { AzureCloud } from '../../../../shared/constants/App.constants';
import { AZURE_DEFAULT_VALUES } from '../../../../shared/constants/defaults/azure.defaults';
import { AZURE_FIELDS } from '../azure-mc-basic/AzureManagementClusterBasic.constants';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { STORE_SECTION_RESOURCES } from '../../../../state-management/reducers/Resources.reducer';
import { StoreDispatch } from '../../../../shared/types/types';

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
    [STORE_SECTION_RESOURCES]: {},
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
