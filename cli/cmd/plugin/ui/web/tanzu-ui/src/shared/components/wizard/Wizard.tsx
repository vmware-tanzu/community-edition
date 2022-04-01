import _ from 'lodash';
import React, { useState } from 'react';

import StepWizard from 'react-step-wizard';
import styled from 'styled-components';
import { STATUS } from '../../constants/App.constants';
import StepNav from './StepNav';

const WizardContainer = styled.div`
    border: 1px solid #0F171C;
    border-top: none;
    overflow: hidden;
`;

function Wizard(props: any) {
    const [tabStatus, setTabStatus] = useState([STATUS.CURRENT, ..._.times(props.children.length - 1, () => STATUS.DISABLED)]);

    return(
        <WizardContainer>
            <StepWizard initialStep={1} nav={<StepNav tabStatus={tabStatus}/>}>
                {props.children.map((child: any, index: number) => React.cloneElement(child, { tabStatus, setTabStatus, key: index }))}
            </StepWizard>
        </WizardContainer>
    );
}

export default Wizard;
