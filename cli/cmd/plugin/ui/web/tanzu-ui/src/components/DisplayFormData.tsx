// React imports
import { CdsButton } from '@cds/react/button';
import React, { useContext } from 'react';

// Library imports
import styled from 'styled-components';

// App imports
import { Store } from '../state-management/stores/Store';
import { STATUS } from '../shared/constants/App.constants';
import { ChildProps } from '../shared/components/wizard/StepNav';

const Container = styled.div`
    margin-top: 30px;
`;

function DisplayFormData(props: ChildProps | any) {
    const { state } = useContext(Store);
    return (
        <Container>{
            Object.entries(state.data).map(([key, val]) => {
                return (<div key={key}> {key} : {val}</div>);
            })
        }
        <CdsButton onClick={() => {
            props.goToStep(props.currentStep + 1);
            const tabStatus = [...props.tabStatus];
            tabStatus[props.currentStep - 1] = STATUS.VALID;
            props.setTabStatus(tabStatus);
        }}>Next</CdsButton>
        </Container>
    );
}

export default DisplayFormData;