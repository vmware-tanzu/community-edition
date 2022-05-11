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
        // Auth - Credential Profile
        PROFILE: '',

        // Auth - Temporary Credentials
        SECRET_ACCESS_KEY: '',
        SESSION_TOKEN: '',
        ACCESS_KEY_ID: '',

        // Auth - General
        EC2_KEY_PAIR: '',

        // Region
        REGION: '',

        // Cluster Data
        CLUSTER_NAME: 'my-aws-cluster',
        CLUSTER_PLAN: 'dev',
        CLUSTER_DEV_NODE_TYPE: 't3a.large',
        CLUSTER_PROD_NODE_TYPE: 'm6a.xlarge',
        CLUSTER_WORKER_NODE_TYPE: 't3a.large',

        // VPC New
        // VPC_NAME: 'temp-vpc-name',
        VPC_CIDR: '10.0.0.0/16',

        // other?
        CREATE_CLOUDFORMATION_STACK: true,
        ENABLE_AUDIT_LOGGING: true,
        ENABLE_BASTION_HOST: true,
        ENABLE_CEIP_PARTICIPATION: true,
        ENABLE_MACHINE_HEALTH_CHECK: true,

        // Kubernetes Networking
        CLUSTER_SERVICE_CIDR: '100.64.0.0/13',
        CLUSTER_POD_CIDR: '100.96.0.0/11',

        // HTTP Proxy & Load Balancer
        HTTP_PROXY_ENABLED: false,
        LOAD_BALANCER_SCHEME_INTERNAL: false,
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
