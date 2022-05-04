// React imports
import React, { useContext } from 'react';

// Library imports
// App imports
import Wizard from '../../shared/components/wizard/Wizard';
import SelectManagementCluster from './SelectManagementCluster';
import { ManagementCluster } from '../../shared/models/ManagementCluster';
import { WcStore } from '../../state-management/stores/Store.wc';
import ClusterTopologyStep from './ClusterTopologyStep';
import './WorkloadClusterWizard.scss';
import ClusterAttributeStep from './ClusterAttributeStep';
import { ClusterClassDefinition, ClusterClassVariableType } from '../../shared/models/ClusterClass';

const fakeServiceRetrievesManagementClusterObjects = (): ManagementCluster[] => {
    return [
        { name: 'shimon-test-cluster-1', provider: 'aws', created: '10/22/2021', description: 'just fooling around' },
        { name: 'some-other-cluster', provider: 'vsphere', created: '1/13/2022', description: 'a very serious cluster' }
    ];
};

const fakeServiceRetrievesClusterClassDefinition = (mc: string): ClusterClassDefinition | undefined => {
    if (mc === '') {
        return undefined
    }
    if (mc === 'some-other-cluster') {
        return {
            name: 'tkg-vsphere-default',
            requiredVariables: [
                { name: 'VSPHERE_CONTROL_PLANE_ENDPOINT', valueType: ClusterClassVariableType.STRING,
                    description: 'kube-apiserver endpoint (IP) for the workload cluster' },
                { name: 'IS_WINDOWS_WORKLOAD_CLUSTER', valueType: ClusterClassVariableType.BOOLEAN,
                    description: 'Is this a Windows-based workload cluster?' },
            ],
            optionalVariables: [
                { name: 'CLUSTER_PLAN', valueType: ClusterClassVariableType.STRING,
                    description: 'plan used for workload cluster: dev or prod',
                    defaultValue: 'dev',
                    possibleValues: ['dev', 'prod']
                }
            ]
        }
    }
    if (mc === 'shimon-test-cluster-1') {
        return {
            name: 'tkg-aws-default',
            requiredVariables: [],
            optionalVariables: [
                { name: 'CLUSTER_PLAN', valueType: ClusterClassVariableType.STRING,
                    description: 'plan used for workload cluster: dev or prod',
                    defaultValue: 'dev',
                    possibleValues: ['dev', 'prod']
                },
                { name: 'AWS_VPC_ID', valueType: ClusterClassVariableType.STRING,
                    description: 'VPC id',
                    defaultValue: '123',
                }
            ],
        }
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
