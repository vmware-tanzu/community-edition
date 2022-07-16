// React imports
import React, { useContext } from 'react';

// App imports
import { AwsStore } from '../../../../state-management/stores/Store.aws';
import { AWS_MC_BASIC_TAB_NAMES } from '../../../../shared/constants/NavRoutes.constants';
import Wizard from '../../../../shared/components/wizard/Wizard';
import ManagementClusterSettings from '../../../../shared/components/management-cluster-settings-step/ManagementClusterSettings';
import ManagementCredentials from './management-credential-step/ManagementCredentials';
import useAwsDeployment from '../../../../shared/services/awsDeployment';
import { AwsService, AWSVirtualMachine } from '../../../../swagger-api';

function AwsManagementClusterBasic() {
    const { awsState, awsDispatch } = useContext(AwsStore);
    const { deployOnAws } = useAwsDeployment();
    const getAWSImageMethod = async (region: string) => {
        const awsImageList = AwsService.getAwsosImages(region);
        return awsImageList;
    };

    return (
        <Wizard tabNames={AWS_MC_BASIC_TAB_NAMES} state={awsState} dispatch={awsDispatch}>
            <ManagementCredentials />
            <ManagementClusterSettings<AWSVirtualMachine>
                deploy={deployOnAws}
                defaultData={awsState}
                getImageMethod={getAWSImageMethod}
                clusterName={'Amazon Machine Image(AMI)'}
            />
        </Wizard>
    );
}

export default AwsManagementClusterBasic;
