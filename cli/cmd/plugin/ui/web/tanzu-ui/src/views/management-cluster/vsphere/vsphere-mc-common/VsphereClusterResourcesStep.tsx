// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';
// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
// App imports
import { NavRoutes } from '../../../../shared/constants/NavRoutes.constants';

export function VsphereClusterResourcesStep() {
    const navigate = useNavigate();

    const navigateToDeploymentProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    const handleMCCreation = () => {
        // TODO: gather the payload, set data into store, call endpoint to create the MC
        navigateToDeploymentProgress();
    };

    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <h3>vSphere Cluster Resources</h3>
            <div cds-layout="grid gap:md cols:12">
                <div>Placeholder for vSphere cluster resources</div>
            </div>
            <CdsButton cds-layout="col:start-1" status="success" onClick={handleMCCreation}>
                <CdsIcon shape="cluster" size="sm"></CdsIcon>
                Create Management cluster
            </CdsButton>
        </div>
    );
}
