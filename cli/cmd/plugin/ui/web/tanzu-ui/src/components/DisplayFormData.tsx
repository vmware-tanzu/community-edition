import React, { useContext } from 'react';
import styled from 'styled-components';
import { Store } from '../stores/store';


const Container = styled.div`
    margin-top: 30px;
`;
function DisplayFormData() {
    const { state } = useContext(Store);
    return (
        <Container>{
            Object.entries(state.data).map(([key, val]) => {
                return (<div key={key}> {key} : {val}</div>);
            })
        }</Container>
    );
}

export default DisplayFormData;