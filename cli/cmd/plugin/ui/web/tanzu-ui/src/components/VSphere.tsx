import React from 'react';
import styled from 'styled-components';

import DisplayFormData from './DisplayFormData';
import TestRender from './TestRender';
import VSphereAuthForm from './VSphereAuthForm';


const FormContainer = styled.div`
    padding: 50px 0;
`;

function VSphere () {
    return (
        <FormContainer>
            <VSphereAuthForm></VSphereAuthForm>
            <DisplayFormData></DisplayFormData>
            <TestRender></TestRender>
        </FormContainer>
    );
}

export default VSphere;