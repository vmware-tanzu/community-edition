// React imports
import React from 'react';

// App imports
import { DOCKER_MC_BASIC_TAB_NAMES } from './DockerManagementClusterBasic.constants';
import Wizard from '../../../../shared/components/wizard/Wizard';
import McPrerequisiteStep from '../docker-mc-common/mc-prerequisite-step/McPrerequisiteStep';
import useDockerDeployment from '../../../../shared/services/dockerDeployment';
import DockerClusterSettings from '../docker-mc-common/mc-cluster-settings-step/DockerClusterSettings';

function DockerManagementClusterBasic() {
    const { deployOnDocker } = useDockerDeployment();
    return (
        <Wizard tabNames={DOCKER_MC_BASIC_TAB_NAMES}>
            <McPrerequisiteStep />
            <DockerClusterSettings deploy={deployOnDocker} />
        </Wizard>
    );
}

export default DockerManagementClusterBasic;
