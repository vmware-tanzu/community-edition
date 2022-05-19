// React imports
import React, { useContext } from 'react';

// App imports
import { DockerStore } from '../../../../state-management/stores/Docker.store';
import { DOCKER_MC_BASIC_TAB_NAMES } from '../../../../shared/constants/NavRoutes.constants';
import Wizard from '../../../../shared/components/wizard/Wizard';
import McPrerequisiteStep from './mc-prerequisite-step/McPrerequisiteStep';
import ManagementClusterSettings from '../../../../shared/components/management-cluster-settings-step/ManagementClusterSettings';
import useDockerDeployment from '../../../../shared/services/dockerDeployment';

function DockerManagementClusterBasic() {
    const { dockerState, dockerDispatch } = useContext(DockerStore);
    const { deployOnDocker } = useDockerDeployment();
    return (
        <Wizard tabNames={DOCKER_MC_BASIC_TAB_NAMES} state={dockerState} dispatch={dockerDispatch}>
            <McPrerequisiteStep />
            <ManagementClusterSettings
                defaultData={dockerState}
                deploy={deployOnDocker}
                message="A single node Management Cluster will be created
                                on your local workstation."
            />
        </Wizard>
    );
}

export default DockerManagementClusterBasic;
