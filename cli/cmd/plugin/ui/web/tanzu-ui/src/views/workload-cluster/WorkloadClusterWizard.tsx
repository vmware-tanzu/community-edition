// React imports
import React, { useContext } from 'react';

// Library imports
// App imports
import './WorkloadClusterWizard.scss';
import { CancelablePromise, ManagementCluster, ManagementService } from '../../swagger-api';
import ClusterAttributeStep from './ClusterAttributeStep';
import ClusterTopologyStep from './ClusterTopologyStep';
import SelectManagementCluster from './SelectManagementCluster';
import { WcStore } from '../../state-management/stores/Store.wc';
import Wizard from '../../shared/components/wizard/Wizard';

const retrieveManagementClusterObjects = (): CancelablePromise<Array<ManagementCluster>> => {
    return ManagementService.getMgmtClusters();
};
// TODO: implement retrieveAvailableClusterClasses
const retrieveAvailableClusterClasses = (mcName: string): CancelablePromise<Array<string>> => {
    console.log(`(Pretending to) retrieve cluster classes for MC ${mcName}`);
    return new CancelablePromise<Array<string>>((resolve) => resolve(['tkg-vsphere-default', 'custom-cluster-class']));
};

const wcTabNames = ['Select a Management Cluster', 'Cluster topology', 'Cluster attributes'] as string[];

function WorkloadClusterWizard(props: any) {
    return (
        <Wizard tabNames={wcTabNames} {...useContext(WcStore)}>
            <SelectManagementCluster retrieveManagementClusters={retrieveManagementClusterObjects} selectedManagementCluster="" />
            <ClusterTopologyStep />
            <ClusterAttributeStep retrieveAvailableClusterClasses={retrieveAvailableClusterClasses} />
        </Wizard>
    );
}

export default WorkloadClusterWizard;
