// React imports
import React from 'react';

// App imports
import { DockerProvider } from '../../../state-management/stores/Docker.store';
import DockerManagementClusterSimple from './wizard-simple/DockerManagementClusterSimple';
import DockerLogo from '../../../assets/docker.svg';
import './DockerManagementCluster.scss';

function DockerManagementCluster() {
    return (
        <DockerProvider>
            <div className="aws-management-container" cds-layout="grid col:12">
                <div cds-layout="col:8">
                    <div cds-layout="col:12 p-b:lg">
                        <span cds-text="title">
                            <img src={DockerLogo} className="logo logo-42" cds-layout="m-r:md" alt="aws logo" />
                            Create Management Cluster on Docker
                        </span>
                    </div>

                    {/* Disable Simple/Advanced banner until advanced settings available */}
                    {/* <ConfigBanner /> */}
                    <DockerManagementClusterSimple />
                </div>
                <div cds-layout="col:4" className="image-container">
                    <div className="mgmt-cluster-admins-img"></div>
                </div>
            </div>
        </DockerProvider>
    );
}

export default DockerManagementCluster;
