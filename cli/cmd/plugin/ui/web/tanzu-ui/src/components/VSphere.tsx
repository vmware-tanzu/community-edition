// React imports
import React from 'react';

// Library imports
import styled from 'styled-components';

// App imports
import DisplayFormData from './DisplayFormData';
import TestRender from './TestRender';
import VSphereAuthForm from './VSphereAuthForm';
import Wizard from '../shared/components/wizard/Wizard';
import ClusterSettings from './ClusterSettings';


const FormContainer = styled.div`
    padding: 50px 0;
`;

function VSphere () {
    return (
        <FormContainer>
            <Wizard>
                <VSphereAuthForm></VSphereAuthForm>
                <ClusterSettings></ClusterSettings>
                <DisplayFormData></DisplayFormData>
                <TestRender></TestRender>
                <TestRender></TestRender>
            </Wizard>
        </FormContainer>
    );
}

export default VSphere;