// React imports
import React, { useContext } from 'react';

// Library imports
import styled from 'styled-components';

// App imports
import TestRender from '../../components/TestRender';
import Wizard from '../../shared/components/wizard/Wizard';
import SelectManagementCluster from './SelectManagementCluster';
import { ManagementCluster } from '../../shared/models/ManagementCluster';
import { WcProvider, WcStore } from '../../state-management/stores/Store.wc';
import ClusterTopologyStep from './ClusterTopologyStep';


const FormContainer = styled.div`
    padding: 50px 0;
`;

const fakeServiceRetrievesManagementClusterObjects = (): ManagementCluster[] => {
    return [
        { name: 'shimon-test-cluster-1', provider: 'aws', created: '10/22/2021', description: 'just fooling around' },
        { name: 'some-other-cluster', provider: 'vsphere', created: '1/13/2022', description: 'a very serious cluster' }
    ];
};

const wcTabNames = ['Select a Management Cluster', 'Cluster topology', 'Cluster attributes'] as string[];

function WorkloadClusterWizard () {
    return (
        <WcProvider>
            <div cds-layout="vertical gap:lg gap@md:xl col@sm:12">
                <FormContainer>
                    <Wizard tabNames={wcTabNames} {...useContext(WcStore)} >
                        <SelectManagementCluster
                                                 retrieveManagementClusters={fakeServiceRetrievesManagementClusterObjects}
                                                 selectedManagementCluster=""
                                                 />
                        <ClusterTopologyStep></ClusterTopologyStep>
                        <TestRender></TestRender>
                    </Wizard>
                </FormContainer>
            </div>
        </WcProvider>
    );
}

export default WorkloadClusterWizard;
