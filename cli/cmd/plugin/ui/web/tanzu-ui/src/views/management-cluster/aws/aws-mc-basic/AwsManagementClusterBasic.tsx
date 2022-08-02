// React imports
import React from 'react';

// App imports
import { AWS_MC_BASIC_TAB_NAMES } from './AwsManagementClusterBasic.constants';
import Wizard from '../../../../shared/components/wizard/Wizard';
import ManagementCredentials from '../aws-mc-common/management-credential-step/ManagementCredentials';
import useAwsDeployment from '../../../../shared/services/awsDeployment';
import AwsClusterSettingsStep from '../AwsClusterSettingsStep';

function AwsManagementClusterBasic() {
    const { deployOnAws } = useAwsDeployment();

    return (
        <Wizard tabNames={AWS_MC_BASIC_TAB_NAMES}>
            <ManagementCredentials />
            <AwsClusterSettingsStep deploy={deployOnAws} />
        </Wizard>
    );
}

export default AwsManagementClusterBasic;
