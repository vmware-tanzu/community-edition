// React imports
import React, { createContext, ReactNode, useReducer } from 'react';

// App imports
import awsReducer from '../../../../views/providers/aws/Aws.reducer';
import { AWS_DEFAULT_VALUES } from '../../../../shared/constants/defaults/aws.defaults';
import { AWS_FIELDS } from '../../../../views/management-cluster/aws/aws-mc-basic/AwsManagementClusterBasic.constants';
import { StoreDispatch } from '../../../../shared/types/types';
import { STORE_SECTION_AWS_RESOURCES } from '../../../../views/providers/aws/AwsResources.reducer';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';

const initialState = {
    [STORE_SECTION_FORM]: {
        // Auth - Credential Profile
        [AWS_FIELDS.PROFILE]: '',

        // Auth - General
        [AWS_FIELDS.EC2_KEY_PAIR]: '',

        // Auth - Temporary Credentials
        [AWS_FIELDS.SECRET_ACCESS_KEY]: '',
        [AWS_FIELDS.SESSION_TOKEN]: '',
        [AWS_FIELDS.ACCESS_KEY_ID]: '',
        // Region
        [AWS_FIELDS.REGION]: '',
        ...AWS_DEFAULT_VALUES,
    },
    [STORE_SECTION_AWS_RESOURCES]: {},
};

const AwsStore = createContext<{
    awsState: { [key: string]: any };
    awsDispatch: StoreDispatch;
}>({
    awsState: initialState,
    awsDispatch: () => null,
});

const AwsProvider: React.FC<{ children: ReactNode }> = ({ children }: { children: ReactNode }) => {
    const [awsState, awsDispatch] = useReducer(awsReducer, initialState);

    return <AwsStore.Provider value={{ awsState, awsDispatch }}>{children}</AwsStore.Provider>;
};

export { AwsStore, AwsProvider };
