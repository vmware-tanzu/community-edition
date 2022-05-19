// React imports
import React from 'react';

// Library imports
import styled from 'styled-components';

// App imports
import DisplayFormData from './DisplayFormData';
import TestRender from './TestRender';
import VSphereAuthForm from './VSphereAuthForm';
import Wizard from '../shared/components/wizard/Wizard';
import ManagementClusterSettings from './ManagementClusterSettings';

const FormContainer = styled.div`
    padding: 50px 0;
`;

function VSphere() {
    return (
        <div cds-layout="vertical gap:lg gap@md:xl col@sm:12">
            <FormContainer>
                <Wizard>
                    <VSphereAuthForm></VSphereAuthForm>
                    <ManagementClusterSettings></ManagementClusterSettings>
                    <DisplayFormData></DisplayFormData>
                    <TestRender></TestRender>
                    <TestRender></TestRender>
                </Wizard>
            </FormContainer>
        </div>
    );
}

export default VSphere;
