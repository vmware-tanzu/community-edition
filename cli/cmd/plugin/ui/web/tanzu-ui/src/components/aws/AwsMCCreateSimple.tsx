// React imports
import React, { useContext } from 'react';

// App imports
import { AwsStore } from '../../state-management/stores/Store.aws';
import { TAB_NAMES } from '../../shared/constants/NavRoutes.constants';
import ManagementClusterSettings from './management-cluster-settings-step/ManagementClusterSettings';
import ManagementCredentials from './management-credential-step/ManagementCredentials';
import Wizard from '../../shared/components/wizard/Wizard';

function AwsMCCreateSimple() {
    const { awsState, awsDispatch } = useContext(AwsStore);
    return (
        <Wizard
            tabNames={TAB_NAMES.awsManagementClusterCreateSimple}
            state={awsState}
            dispatch={awsDispatch}
        >
            <ManagementCredentials />
            <ManagementClusterSettings />
        </Wizard>
    );
}

export default AwsMCCreateSimple;
