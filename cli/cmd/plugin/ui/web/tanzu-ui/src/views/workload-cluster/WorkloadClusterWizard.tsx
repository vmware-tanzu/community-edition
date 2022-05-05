// React imports
import React, { useContext } from 'react';

// Library imports
// App imports
import './WorkloadClusterWizard.scss';
import ClusterAttributeStep from './ClusterAttributeStep';
import { ClusterClassDefinition } from '../../shared/models/ClusterClass';
import ClusterTopologyStep from './ClusterTopologyStep';
import { FakeClusterClassAws, FakeClusterClassAzure, FakeClusterClassDocker, FakeClusterClassVsphere } from './fake-cluster-class-service';
import { ManagementCluster } from '../../shared/models/ManagementCluster';
import SelectManagementCluster from './SelectManagementCluster';
import { WcStore } from '../../state-management/stores/Store.wc';
import Wizard from '../../shared/components/wizard/Wizard';

const fakeServiceRetrievesManagementClusterObjects = (): ManagementCluster[] => {
    return [
        { name: 'shimon-test-cluster-1', provider: 'aws', created: '10/22/2021', description: 'This cluster should be deleted soon' },
        { name: 'some-other-cluster', provider: 'vsphere', created: '4/13/2022', description: 'a very high-level cluster' },
        { name: 'docker-foobar-cluster', provider: 'docker', created: '2/14/2022', description: 'a local fun cluster' },
        { name: 'azure-clown-cluster', provider: 'azure', created: '3/15/2022', description: 'beware, Caesar, a backstabbing cluster' },
    ];
};

const fakeServiceRetrievesClusterClassDefinition = (mc: string): ClusterClassDefinition | undefined => {
    switch (mc) {
    case 'some-other-cluster':
        return FakeClusterClassVsphere
    case 'shimon-test-cluster-1':
        return FakeClusterClassAws
    case 'docker-foobar-cluster':
        return FakeClusterClassDocker
    case 'azure-clown-cluster':
        return FakeClusterClassAzure
    default:
        return undefined
    }
}

const wcTabNames = ['Select a Management Cluster', 'Cluster topology', 'Cluster attributes'] as string[];

function WorkloadClusterWizard (props: any) {
    return (
        <Wizard tabNames={wcTabNames} {...useContext(WcStore)} >
            <SelectManagementCluster
                retrieveManagementClusters={fakeServiceRetrievesManagementClusterObjects}
                selectedManagementCluster=""
            />
            <ClusterTopologyStep></ClusterTopologyStep>
            <ClusterAttributeStep
                retrieveClusterClassDefinition={fakeServiceRetrievesClusterClassDefinition}
            ></ClusterAttributeStep>
        </Wizard>
    );
}

export default WorkloadClusterWizard;
