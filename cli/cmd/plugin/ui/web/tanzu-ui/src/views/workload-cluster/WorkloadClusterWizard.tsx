// React imports
import React, { useContext } from 'react';

// Library imports
// App imports
import './WorkloadClusterWizard.scss';
import { CancelablePromise, ManagementCluster, ManagementService } from '../../swagger-api';
import ClusterAttributeStep from './ClusterAttributeStep';
import { ClusterClassDefinition } from '../../shared/models/ClusterClass';
import ClusterTopologyStep from './ClusterTopologyStep';
import { FakeClusterClassAws, FakeClusterClassAzure, FakeClusterClassDocker, FakeClusterClassVsphere } from './fake-cluster-class-service';
import SelectManagementCluster from './SelectManagementCluster';
import { WcStore } from '../../state-management/stores/Store.wc';
import Wizard from '../../shared/components/wizard/Wizard';

const retrieveManagementClusterObjects = (): CancelablePromise<Array<ManagementCluster>> => {
    return ManagementService.getMgmtClusters();
};

const wcTabNames = ['Select a Management Cluster', 'Cluster topology', 'Cluster attributes'] as string[];

function WorkloadClusterWizard(props: any) {
    return (
        <Wizard tabNames={wcTabNames} {...useContext(WcStore)}>
            <SelectManagementCluster retrieveManagementClusters={retrieveManagementClusterObjects} selectedManagementCluster="" />
            <ClusterTopologyStep></ClusterTopologyStep>
            <ClusterAttributeStep></ClusterAttributeStep>
        </Wizard>
    );
}

export default WorkloadClusterWizard;
