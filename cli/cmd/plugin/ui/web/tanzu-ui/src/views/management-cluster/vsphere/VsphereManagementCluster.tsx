// React imports
import React from 'react';
// App imports
import VsphereLogo from '../../../assets/vsphere.svg';
import VsphereManagementClusterBasic from './vsphere-mc-basic/VsphereManagementClusterBasic';
import { VsphereProvider } from './Store.vsphere.mc';

function VsphereManagementCluster() {
    return (
        <VsphereProvider>
            <div className="vsphere-management-container" cds-layout="grid col:12">
                <div cds-layout="col:8">
                    <div cds-layout="col:12 p-b:lg">
                        <span cds-text="title">
                            <img src={VsphereLogo} className="logo logo-42" cds-layout="m-r:md" alt="vsphere logo" />
                            Create Management Cluster on vSphere
                        </span>
                    </div>

                    {/* Disable Basic/Advanced banner until advanced settings available */}
                    {/* <ConfigBanner /> */}
                    <VsphereManagementClusterBasic />
                </div>
                <div cds-layout="col:4" className="image-container">
                    <div className="mgmt-cluster-admins-img"></div>
                </div>
            </div>
        </VsphereProvider>
    );
}

export default VsphereManagementCluster;
