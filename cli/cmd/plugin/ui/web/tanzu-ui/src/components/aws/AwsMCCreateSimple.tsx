// React imports
import React, { useContext } from 'react';

// App imports
import { AwsStore } from '../../state-management/stores/Store.aws';
import { TAB_NAMES } from '../../shared/constants/NavRoutes.constants';
import Wizard from '../../shared/components/wizard/Wizard';
import ManagementClusterSettings from '../../shared/components/management-cluster-settings-step/ManagementClusterSettings';
import ManagementCredentials from './management-credential-step/ManagementCredentials';
import useAwsDeployment from '../../shared/services/awsDeployment';

function AwsMCCreateSimple() {
    const { awsState, awsDispatch } = useContext(AwsStore);
    const { deployOnAws } = useAwsDeployment();

    return (
        <Wizard tabNames={TAB_NAMES.awsManagementClusterCreateSimple} state={awsState} dispatch={awsDispatch}>
            <ManagementCredentials />
            <ManagementClusterSettings deploy={deployOnAws} defaultData={awsState} />
        </Wizard>
    );
}

export default AwsMCCreateSimple;
