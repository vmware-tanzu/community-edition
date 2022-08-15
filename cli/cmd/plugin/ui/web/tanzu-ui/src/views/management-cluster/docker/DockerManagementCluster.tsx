// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { CdsIconButton } from '@cds/react/button';

// App imports
import DockerManagementClusterBasic from './docker-mc-basic/DockerManagementClusterBasic';
import { DockerProvider } from '../../../state-management/stores/Docker.store';
import DockerLogo from '../../../assets/docker.svg';
import './DockerManagementCluster.scss';
import { NavRoutes } from '../../../shared/constants/NavRoutes.constants';

function DockerManagementCluster() {
    const navigate = useNavigate();

    return (
        <DockerProvider>
            <div cds-layout="col:12">
                <div cds-layout="p-b:lg">
                    <span cds-text="title">
                        <CdsIconButton
                            cds-layout="p-t:md"
                            action="flat"
                            status="primary"
                            onClick={() => navigate(NavRoutes.MANAGEMENT_CLUSTER_SELECT_PROVIDER)}
                        >
                            <CdsIcon shape="arrow" direction="left" size="lg"></CdsIcon>
                        </CdsIconButton>
                        <img src={DockerLogo} className="logo logo-42" cds-layout="m-r:md" alt="docker logo" />
                        Create a Management Cluster on Docker
                    </span>
                </div>
                {/* Disable Basic/Advanced banner until advanced settings available */}
                {/* <ConfigBanner /> */}
                <DockerManagementClusterBasic />
            </div>
        </DockerProvider>
    );
}

export default DockerManagementCluster;
