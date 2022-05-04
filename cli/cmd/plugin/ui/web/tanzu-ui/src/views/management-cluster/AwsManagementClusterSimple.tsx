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
            <div className="aws-management-container" cds-layout="grid col:12">
                <div cds-layout="col:8">
                    <h2>
                        <div cds-layout="m-r:lg" className="aws-sm-logo"></div>
                        <span cds-text="heading">
                            Create Management Cluster on Amazon Web Services
                        </span>
                    </h2>
                    <ConfigBanner />
                    <AwsMCCreateSimple />
                </div>
                <div cds-layout="col:4" className="image-container">
                    <div className="mgmt-cluster-admins-img"></div>
                </div>
            </div>
        </AwsProvider>
    );
}

export default AwsManagementClusterSimple;
