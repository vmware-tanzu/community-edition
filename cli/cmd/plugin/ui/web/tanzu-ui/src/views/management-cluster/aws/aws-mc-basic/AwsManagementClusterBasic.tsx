// React imports
import React, { useContext } from 'react';

// App imports
import { AwsStore } from '../../../../state-management/stores/Store.aws';
import { AWS_MC_BASIC_TAB_NAMES } from './AwsManagementClusterBasic.constants';
import Wizard from '../../../../shared/components/wizard/Wizard';
import ManagementClusterSettings from '../../../../shared/components/management-cluster-settings-step/ManagementClusterSettings';
import ManagementCredentials from '../aws-mc-common/management-credential-step/ManagementCredentials';
import useAwsDeployment from '../../../../shared/services/awsDeployment';

function AwsManagementClusterBasic() {
    const { awsState, awsDispatch } = useContext(AwsStore);
    const { deployOnAws } = useAwsDeployment();

    return (
        <Wizard tabNames={AWS_MC_BASIC_TAB_NAMES} state={awsState} dispatch={awsDispatch}>
            <ManagementCredentials />
            <ManagementClusterSettings deploy={deployOnAws} defaultData={awsState} />
        </Wizard>
    );
}

export default AwsManagementClusterBasic;
