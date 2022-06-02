// React imports
import React, { useContext } from 'react';

// App imports
import { AwsStore } from '../../../state-management/stores/Store.aws';
import { UMC_BASIC_TAB_NAMES } from '../../../shared/constants/NavRoutes.constants';
import Wizard from '../../../shared/components/wizard/Wizard';
import ManagementClusterSettings from '../../../shared/components/management-cluster-settings-step/ManagementClusterSettings';
import ManagementCredentials from '../../management-cluster/aws/wizard-basic/management-credential-step/ManagementCredentials';
import useAwsDeployment from '../../../shared/services/awsDeployment';

function UnmanagedClusterBasic() {
    const { awsState, awsDispatch } = useContext(AwsStore);
    const { deployOnAws } = useAwsDeployment();

    return (
        <Wizard tabNames={UMC_BASIC_TAB_NAMES} state={awsState} dispatch={awsDispatch}>
            <ManagementClusterSettings deploy={deployOnAws} defaultData={awsState} />
            <ManagementCredentials />
        </Wizard>
    );
}

export default UnmanagedClusterBasic;
