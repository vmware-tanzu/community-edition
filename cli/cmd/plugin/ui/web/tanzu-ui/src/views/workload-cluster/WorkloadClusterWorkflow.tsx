// React imports
import React, { useContext } from 'react';

// Library imports
import styled from 'styled-components';

// App imports
import WorkloadClusterWizard from './WorkloadClusterWizard';
import { WcProvider } from '../../state-management/stores/Store.wc';


const FormContainer = styled.div`
    padding: 50px 0;
`;

function WorkloadClusterWorkflow () {
    return (
        <WcProvider>
            <FormContainer>
                <WorkloadClusterWizard />
            </FormContainer>
        </WcProvider>
    );
}

export default WorkloadClusterWorkflow;
