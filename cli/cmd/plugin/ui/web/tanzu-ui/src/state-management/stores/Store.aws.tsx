// React imports
import React, {
    createContext,
    ReactNode,
    useReducer,
} from 'react';

// App imports
import { StoreDispatch } from '../../shared/types/types';
import awsReducer from '../reducers/Aws.reducer';

const initialState = {
    data: {
        PROFILE: 'profile1',
        SECRET_ACCESS_KEY: '123',
        SESSION_TOKEN: '4343',
        ACCESS_KEY_ID: '555',
        REGION: '',
        EC2_KEY_PAIR: '',
        CLUSTER_NAME: 'testcluster'
    },
};

const AwsStore = createContext<{
    awsState: { [key: string]: any };
    awsDispatch: StoreDispatch;
}>({
    awsState: initialState,
    awsDispatch: () => null,
});

const AwsProvider: React.FC<{ children: ReactNode }> = ({
    children,
}: {
    children: ReactNode;
}) => {
    const [awsState, awsDispatch] = useReducer(
        awsReducer,
        initialState
    );

    return (
        <AwsStore.Provider value={{ awsState, awsDispatch }}>
            {children}
        </AwsStore.Provider>
    );
};

export { AwsStore, AwsProvider };
