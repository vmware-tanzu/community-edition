// React imports
import React from 'react';

// App imports
import { DockerProvider } from '../../../state-management/stores/Docker.store';
import DockerManagementClusterBasic from './docker-mc-basic/DockerManagementClusterBasic';
import DockerLogo from '../../../assets/docker.svg';
import './DockerManagementCluster.scss';

function DockerManagementCluster() {
    return (
        <DockerProvider>
            <div cds-layout="grid col:12">
                <div cds-layout="col:8">
                    <div cds-layout="col:12 p-b:lg">
                        <span cds-text="title">
                            <img src={DockerLogo} className="logo logo-42" cds-layout="m-r:md" alt="docker logo" />
                            Create Management Cluster on Docker
                        </span>
                    </div>
                    {/* Disable Basic/Advanced banner until advanced settings available */}
                    {/* <ConfigBanner /> */}
                    <DockerManagementClusterBasic />
                </div>
                <div cds-layout="col:4" className="image-container">
                    <div className="mgmt-cluster-admins-img"></div>
                </div>
            </div>
        </DockerProvider>
    );
}

export default DockerManagementCluster;
