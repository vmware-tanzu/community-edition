// React imports
import React, { useContext } from 'react';

// Library imports
// App imports
import './WorkloadClusterWizard.scss';
import { CancelablePromise, ManagementCluster, ManagementService } from '../../swagger-api';
import ClusterAttributeStep from './ClusterAttributeStep';
import ClusterTopologyStep from './ClusterTopologyStep';
import SelectManagementCluster from './SelectManagementCluster';
import { WcStore } from './Store.wc';
import Wizard from '../../shared/components/wizard/Wizard';

const retrieveManagementClusterObjects = (): CancelablePromise<Array<ManagementCluster>> => {
    return ManagementService.getMgmtClusters();
};
// TODO: implement retrieveAvailableClusterClasses
const retrieveAvailableClusterClasses = (mcName: string): CancelablePromise<Array<string>> => {
    let result = [] as string[];
    console.log(`(Pretending to) retrieve cluster classes for MC ${mcName}`);
    if (mcName.includes('vsphere')) {
        result = ['tkg-vsphere-default', 'custom-cluster-class'];
    } else if (mcName.includes('aws')) {
        result = ['tkg-aws-default'];
    }
    return new CancelablePromise<Array<string>>((resolve) => resolve(result));
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
