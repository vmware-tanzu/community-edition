// React imports
import React from 'react';

// App imports
import { AzureProvider } from '../../../state-management/stores/Azure.store';
import AzureLogo from '../../../assets/azure.svg';
import AzureManagementClusterBasic from './azure-mc-basic/AzureManagementClusterBasic';
import './AzureManagementCluster.scss';

function AzureManagementCluster() {
    return (
        <AzureProvider>
            <div cds-layout="grid col:12">
                <div cds-layout="col:8">
                    <div cds-layout="col:12 p-b:lg">
                        <span cds-text="title">
                            <img src={AzureLogo} className="logo logo-42" cds-layout="m-r:md" alt="azure logo" />
                            Create Management Cluster on Azure
                        </span>
                    </div>

                    {/* Disable Basic/Advanced banner until advanced settings available */}
                    {/* <ConfigBanner /> */}
                    <AzureManagementClusterBasic />
                </div>
                <div cds-layout="col:4" className="image-container">
                    <div className="mgmt-cluster-admins-img"></div>
                </div>
            </div>
        </AzureProvider>
    );
}

export default AzureManagementCluster;
