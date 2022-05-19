// React imports
import React from 'react';

// App imports
import { AwsProvider } from '../../../state-management/stores/Store.aws';
import AwsManagementClusterBasic from './wizard-basic/AwsManagementClusterBasic';
import './AwsManagementCluster.scss';
import AwsLogo from '../../../assets/aws.svg';

function AwsManagementCluster() {
    return (
        <AwsProvider>
            <div className="aws-management-container" cds-layout="grid col:12">
                <div cds-layout="col:8">
                    <div cds-layout="col:12 p-b:lg">
                        <span cds-text="title">
                            <img src={AwsLogo} className="logo logo-42" cds-layout="m-r:md" alt="aws logo" />
                            Create Management Cluster on Amazon Web Services
                        </span>
                    </div>

                    {/* Disable Basic/Advanced banner until advanced settings available */}
                    {/* <ConfigBanner /> */}
                    <AwsManagementClusterBasic />
                </div>
                <div cds-layout="col:4" className="image-container">
                    <div className="mgmt-cluster-admins-img"></div>
                </div>
            </div>
        </AwsProvider>
    );
}

export default AwsManagementCluster;
