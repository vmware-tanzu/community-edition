// React imports
import React from 'react';

// App imports
import { AwsProvider } from '../../state-management/stores/Store.aws';
import AwsMCCreateSimple from '../../components/aws/AwsMCCreateSimple';
import ConfigBanner from '../../shared/components/ConfigBanner/ConfigBanner';
import './AwsManagementClusterSimple.scss';

function AwsManagementClusterSimple() {
    return (
        <AwsProvider>
            <div className="aws-management-container">
                <h2>
                    <div className="aws-sm-logo logo-space"></div>
                    <span>
                        Create Management Cluster on Amazon Web Services
                    </span>
                </h2>
                <ConfigBanner />
                <AwsMCCreateSimple />
                <div className="mgmt-cluster-admins-container">
                    <div className="mgmt-cluster-admins"></div>
                </div>
            </div>
        </AwsProvider>
    );
}

export default AwsManagementClusterSimple;
